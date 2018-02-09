package service

import (
	"github.com/golang/glog"
	"gopkg.in/mgo.v2"

	"github.com/qetuantuan/jengo_recap/vo"
	"github.com/qetuantuan/jengo_recap/dao"
	"github.com/qetuantuan/jengo_recap/model"
)

type BuildLogReader interface {
	GetLog(params *vo.GetBuildLogParams) ([]byte, error)
}

type BuildLogWriter interface {
	AppendLog(logUri string, content []byte) (err error)
	PutLog(content []byte) (logUri string, err error)
}

type BuildLogService interface {
	BuildLogReader
	BuildLogWriter
}

var _ BuildLogService = &LocalBuildLogService{}

type LocalBuildLogService struct {
	Md dao.LogDao
}

func (l *LocalBuildLogService) PutLog(content []byte) (logUri string, err error) {
	logUri, err = l.Md.AddLog(content) // row id as Uri temporarily
	if err != nil {
		glog.Warningf("write logs to mongodb failed! error:%v", err)
		err = MongoError
		return

	}
	return
}

func (l *LocalBuildLogService) AppendLog(logUri string, content []byte) (err error) {
	return
}

func (l *LocalBuildLogService) GetLog(params *vo.GetBuildLogParams) (content []byte, err error) {
	var log model.BuildLog
	log, err = l.Md.GetLog(params.Uri)
	if err != nil {
		if err == mgo.ErrNotFound {
			glog.Warningf("No such log with logid : %v ", params.Uri)
			err = NotFoundError
			return
		} else {
			glog.Error("get log failed! error: ", err)
			err = MongoError
			return
		}
	}
	if params.Limit == 0 {
		params.Limit = 100
	}

	realLimit := params.Limit

	if len(log.Content) >= params.Offset+params.Limit {
		content = log.Content[params.Offset : params.Offset+params.Limit]
	} else {
		realLimit = len(log.Content) - params.Offset
		if realLimit <= 0 {
			content = []byte{}
			return
		}
		content = log.Content[params.Offset:]
	}

	return
}
