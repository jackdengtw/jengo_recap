package dao

import (
	"gopkg.in/mgo.v2/bson"

	"github.com/qetuantuan/jengo_recap/model"
)

type LogMongoDao struct {
	MongoDao
}

func (md *LogMongoDao) Init(d *MongoDao) (err error) {
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

func (md *LogMongoDao) AddLog(logs []byte) (id string, err error) {
	session := md.GSession.Copy()
	defer session.Close()
	lc := session.DB(repoDbName).C(logCol)
	oId := bson.NewObjectId()
	id = oId.Hex()
	err = lc.Insert(&model.BuildLog{Id: id, Content: string(logs)})
	return
}

func (md *LogMongoDao) GetLog(id string) (log model.BuildLog, err error) {
	session := md.GSession.Copy()
	defer session.Close()
	lc := session.DB(repoDbName).C(logCol)
	err = lc.FindId(id).One(&log)
	return
}
