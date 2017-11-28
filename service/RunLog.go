package service

import (
	"io/ioutil"
	"math/rand"
	"time"

	"github.com/golang/glog"
	"github.com/qetuantuan/jengo_recap/api"
	"github.com/qetuantuan/jengo_recap/dao"
	"github.com/qetuantuan/jengo_recap/model"
	"gopkg.in/mgo.v2"
)

type RunLogServiceInterface interface {
	RunLogGetter
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

type RunLogGetter interface {
	GetRunLog(params *api.GetRunLogParams) (*api.RunLog, error)
}

const max_log_length = 500

var FileDir = "/tmp/engine/"

func (u *RunLogService) GetRunLog(params *api.GetRunLogParams) (*api.RunLog, error) {
	f := FileDir + params.RunId + ".log"
	s := params.RunId

	var logObj = &model.RunLog{}
	if params.Limit == 0 {
		params.Limit = 100
	}

	realLimit := params.Limit

	if len(s) >= params.Offset+params.Limit {
		logObj.Content = s[params.Offset : params.Offset+params.Limit]
	} else {
		if params.Offset+params.Limit > max_log_length {
			realLimit = max_log_length - params.Offset
			if realLimit < 0 {
				logObj.Content = ""
				return logObj.ToApiObj(), nil
			}
		}
		logObj.Content = s + GetRandomString(params.Offset+realLimit-len(s))
		if err := ioutil.WriteFile(f, []byte(logObj.Content), 0644); err != nil {
			return logObj.ToApiObj(), err
		}
	}

	logObj.RunId = params.RunId
	logObj.FileName = params.RunId
	logObj.Length = realLimit

	return logObj.ToApiObj(), nil
}

func GetRandomString(len int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= len; i++ {
		result = append(result, bytes[r.Intn(62)])
		if i%20 == 0 {
			result = append(result, byte('\n'))
			i++
		}
	}
	return string(result)
}
