package dao

import (
	"reflect"
	"strings"

	"github.com/golang/glog"
	"gopkg.in/mgo.v2"
)

type MongoDao struct {
	Url      string
	GSession *mgo.Session
	Inited   bool

	fieldJsonTypes map[string]reflect.Type
}

func (md *MongoDao) Init() (err error) {
	md.GSession, err = mgo.Dial(md.Url)
	/*
		var l MgoLog
		mgo.SetLogger(l)
		mgo.SetDebug(false)
	*/
	if err == nil {
		md.Inited = true
	}

	return err
}

func (md *MongoDao) initFieldJsonTypes(t reflect.Type) {
	if md.fieldJsonTypes == nil {
		md.fieldJsonTypes = make(map[string]reflect.Type)
	}

	// Not asserting a struct type since
	// otherwise panic when initing
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		ft := f.Type
		for ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}

		s := f.Tag.Get("json")
		s = strings.TrimSpace(strings.Split(s, ",")[0]) // name
		if s == "" {
			md.fieldJsonTypes[strings.ToLower(f.Name)] = ft
		} else {
			md.fieldJsonTypes[s] = ft
		}
	}
}

func (md *MongoDao) checkFieldType(updateData map[string]interface{}) bool {
	for k, v := range updateData {
		if t, ok := md.fieldJsonTypes[k]; ok {
			actual := reflect.ValueOf(v).Type()
			for actual.Kind() == reflect.Ptr {
				actual = actual.Elem()
			}
			if t.Kind() != actual.Kind() {
				glog.Errorf("Checking field type failed: expected is %v, actual is %v", t, actual)
				return false
			}
		} else {
			return false
		}
	}
	return true
}
