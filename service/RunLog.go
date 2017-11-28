package service

import (
	"github.com/golang/glog"
	"github.com/qetuantuan/jengo_recap/api"
	"github.com/qetuantuan/jengo_recap/dao"
	"github.com/qetuantuan/jengo_recap/model"
	"gopkg.in/mgo.v2"
)

type RunLogServiceInterface interface {
	PutLog(logs []byte, buildId, runId string) (id string, err error)
	GetLog(id string) (log api.RunLog, err error)
}

type RunLogService struct {
	Md *dao.ProjectMongoDao
}

func (l *RunLogService) PutLog(logs []byte, buildId, runId string) (id string, err error) {
	id, err = l.Md.AddLog(logs)
	if err != nil {
		glog.Warningf("write logs to mongodb failed! error:%v", err)
		err = MongoError
		return

	}
	if err = l.Md.UpdateRunLog(buildId, runId, id); err != nil {
		glog.Warningf("update logid in run info failed! error:%v ", err)
		err = MongoError
		return
	}
	return
}

func (l *RunLogService) GetLog(id string) (runLog api.RunLog, err error) {
	var log model.RunLog
	log, err = l.Md.GetLog(id)
	if err != nil {
		if err == mgo.ErrNotFound {
			glog.Warningf("No such log with logid : %v ", id)
			err = NotFoundError
			return
		} else {
			glog.Error("get log failed! error: ", err)
			err = MongoError
			return
		}
	}
	runLog = *log.ToApiObj()
	return
}
