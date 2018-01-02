package dao

import (
	"errors"

	"github.com/golang/glog"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/qetuantuan/jengo_recap/model"
)

type UserMongoDao struct {
	MongoDao
}

func (md *UserMongoDao) Init(d *MongoDao) (err error) {
	if d == nil {
		err = md.MongoDao.Init()
	} else {
		md.MongoDao = *d
		if !d.Inited {
			err = d.Init()
		}
	}
	return err
}

func (md *UserMongoDao) GetUser(userId string) (user model.User, err error) {
	session := md.GSession.Copy()
	defer session.Close()
	uc := session.DB(userDbName).C(userCol02)
	err = uc.FindId(userId).One(&user)
	if err != nil {
		glog.Error("err when GetUser because of :", err)
	}
	return
}

func (md *UserMongoDao) UpsertUser(user *model.User) (err error) {
	session := md.GSession.Copy()
	defer session.Close()
	uc := session.DB(userDbName).C(userCol02)
	_, err = uc.UpsertId(user.UserId, &user)
	if err != nil {
		glog.Error("err when CreateUser because of :", err)
	}
	return
}

func (md *UserMongoDao) GetUserByLogin(loginName string, auth string) (user *model.User, err error) {
	session := md.GSession.Copy()
	defer session.Close()
	uc := session.DB(userDbName).C(userCol02)
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

func (md *UserMongoDao) UpdateToken(userId, sourceType, scmId string, token []byte) (err error) {
	session := md.GSession.Copy()
	defer session.Close()
	uc := session.DB(userDbName).C(userCol02)
	err = uc.Update(bson.M{"_id": userId, sourceType + ".id": scmId},
		bson.M{"$set": bson.M{sourceType + ".$.token": token}})
	return
}
