package dao

import (
	"fmt"
	"reflect"

	"github.com/golang/glog"
	"gopkg.in/mgo.v2/bson"

	"github.com/qetuantuan/jengo_recap/model"
)

type EngineBuildDaoInterface interface {
	UpdateBuildProperties(string, map[string]interface{}) error
	InsertBuild(*model.Build) error
	ListBuilds(map[string]interface{}, int, int) (model.Builds, error)
	GetBuild(string) (*model.Build, error)
}

type EngineBuildDao struct {
	MongoDao
}

func (md *EngineBuildDao) Init(d *MongoDao) (err error) {
	if d == nil {
		err = md.MongoDao.Init()
	} else {
		md.MongoDao = *d
		if !d.Inited {
			err = d.Init()
		}
	}

	md.initFieldJsonTypes(reflect.TypeOf(model.Build{}))

	fmt.Printf("%v\n", md.initFieldJsonTypes)

	return err
}

func (md *EngineBuildDao) UpdateBuildProperties(BuildId string, updateData map[string]interface{}) (err error) {
	if !md.checkFieldType(updateData) {
		err = ErrorTypeNotMatch
		return
	}
	session := md.GSession.Copy()
	defer session.Close()
	pc := session.DB(engineDbName).C(buildCol)
	err = pc.Update(bson.M{"_id": BuildId}, bson.M{"$set": updateData})
	if err != nil {
		glog.Errorf("Update failed for %s %v", BuildId, err)
	}
	return
}

func (md *EngineBuildDao) InsertBuild(Build *model.Build) (err error) {
	session := md.GSession.Copy()
	defer session.Close()
	pc := session.DB(engineDbName).C(buildCol)

	err = pc.Insert(*Build)
	if err != nil {
		glog.Errorf("UpsertId failed for %v %v", Build.Id, err)

	}
	return
}

func (md *EngineBuildDao) ListBuilds(query map[string]interface{}, limitCount int, offset int) (Builds model.Builds, err error) {
	session := md.GSession.Copy()
	defer session.Close()
	pc := session.DB(engineDbName).C(buildCol)
	var BuildCs []model.Build
	err = pc.Find(query).Skip(offset).Limit(limitCount).All(&BuildCs)
	for _, r := range BuildCs {
		Builds = append(Builds, r)
	}
	return
}

func (md *EngineBuildDao) GetBuild(id string) (Build *model.Build, err error) {
	session := md.GSession.Copy()
	defer session.Close()
	pc := session.DB(engineDbName).C(buildCol)
	var r model.Build
	err = pc.FindId(id).One(&r)
	return &r, err
}
