package dao

import (
	"errors"

	"github.com/golang/glog"
	"github.com/qetuantuan/jengo_recap/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	SCM_COLUMN  = "scms"
	AUTH_COLUMN = "auths"
)

var dbName string = "users"
var userCol02 string = "user02"

type MongoDao struct {
	Url      string
	GSession *mgo.Session
}

type MgoLog int

func (m MgoLog) Output(calldepth int, s string) error {
	// if glog.V(glog.Level(calldepth)) {
	// TODO: Not sure why user service don't have glog info file log
	//       use error instead for now
	glog.Error(s)
	// }
	return nil
}

func (md *MongoDao) Init() (err error) {
	md.GSession, err = mgo.Dial(md.Url)
	var l MgoLog
	mgo.SetLogger(l)
	// TODO: make debug a command line option
	// This is to debug Mongo queries
	mgo.SetDebug(false)
	return err
}

func (md *MongoDao) GetUser(userId string) (user model.User, err error) {
	session := md.GSession.Copy()
	defer session.Close()
	uc := session.DB(dbName).C(userCol02)
	err = uc.FindId(userId).One(&user)
	if err != nil {
		glog.Error("err when GetUser because of :", err)
	}
	return
}

func (md *MongoDao) CreateUser(user *model.User) (err error) {
	session := md.GSession.Copy()
	defer session.Close()
	uc := session.DB(dbName).C(userCol02)
	_, err = uc.UpsertId(user.UserId, &user)
	if err != nil {
		glog.Error("err when CreateUser because of :", err)
	}
	return
}

func (md *MongoDao) GetUserByLogin(loginName string, auth string) (user *model.User, err error) {
	session := md.GSession.Copy()
	defer session.Close()
	uc := session.DB(dbName).C(userCol02)
	var users []model.User
	err = uc.Find(bson.M{
		"auths.login_name":     loginName,
		"auths.auth_source_id": auth,
		"auths.primary":        true,
	}).All(&users)
	if err != nil {
		glog.Error("err when GetUser because of :", err)
		return
	}
	if len(users) > 1 {
		err = errors.New("More than one users found with the same loginName and auth")
		glog.Errorf("%v: %v", err, users)
		return
	} else if len(users) < 1 {
		err = mgo.ErrNotFound
		return
	}

	return &users[0], nil
}

func (md *MongoDao) UpdateToken(userId, sourceType, scmId string, token []byte) (err error) {
	session := md.GSession.Copy()
	defer session.Close()
	uc := session.DB(dbName).C(userCol02)
	err = uc.Update(bson.M{"_id": userId, sourceType + ".id": scmId},
		bson.M{"$set": bson.M{sourceType + ".$.token": token}})
	return
}
