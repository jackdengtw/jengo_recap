package dao

import (
	"github.com/golang/glog"

	"github.com/qetuantuan/jengo_recap/model"
)

type HookMongoDao struct {
	MongoDao
}

func (md *HookMongoDao) Init(d *MongoDao) (err error) {
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

func (md *HookMongoDao) AddHook(hook model.GithubHook) (err error) {
	session := md.GSession.Copy()
	defer session.Close()
	hc := session.DB(repoDbName).C(hookCol)
	_, err = hc.UpsertId(hook.Id, &hook)
	if err != nil {
		glog.Error("err when Create/add Hook because of :", err)
	}
	return
}

func (md *HookMongoDao) GetHook(HookId string) (hook model.GithubHook, err error) {
	session := md.GSession.Copy()
	defer session.Close()
	hc := session.DB(repoDbName).C(hookCol)
	err = hc.FindId(HookId).One(&hook)
	return
}

func (md *HookMongoDao) DeleteHook(HookId string) (err error) {
	session := md.GSession.Copy()
	defer session.Close()
	hc := session.DB(repoDbName).C(hookCol)
	err = hc.RemoveId(HookId)
	return
}
