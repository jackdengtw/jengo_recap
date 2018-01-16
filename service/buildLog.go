package service

import (
	"io/ioutil"
	"math/rand"
	"time"

	"github.com/golang/glog"
	"github.com/qetuantuan/jengo_recap/api"
	"github.com/qetuantuan/jengo_recap/dao"
	"github.com/qetuantuan/jengo_recap/model"
)

type BuildLogReader interface {
	GetLog(params *api.GetBuildLogParams) (*model.BuildLog, error)
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
	Md *dao.LogMongoDao
}

const max_log_length = 500

var FileDir = "/tmp/engine/"

func (l *LocalBuildLogService) PutLog(content []byte) (logUri string, err error) {
	logUri, err = l.Md.AddLog(content)
	if err != nil {
		glog.Warningf("write logs to mongodb failed! error:%v", err)
		err = MongoError
		return

	}
	return
	/*
		if err = l.Md.UpdateBuildLog(buildId, BuildId, id); err != nil {
			glog.Warningf("update logid in Build info failed! error:%v ", err)
			err = MongoError
			return
		}
		return
	*/
}

func (l *LocalBuildLogService) AppendLog(logUri string, content []byte) (err error) {
	return
}

/*
func (l *LocalBuildLogService) GetLog(id string) (BuildLog api.BuildLog, err error) {
	var log model.BuildLog
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
	BuildLog = *log.ToApiObj()
	return
}
*/

func (u *LocalBuildLogService) GetLog(params *api.GetBuildLogParams) (*model.BuildLog, error) {
	f := FileDir + params.BuildId + ".log"
	s := params.BuildId

	var logObj = &model.BuildLog{}
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
				return logObj, nil
			}
		}
		logObj.Content = s + GetRandomString(params.Offset+realLimit-len(s))
		if err := ioutil.WriteFile(f, []byte(logObj.Content), 0644); err != nil {
			return logObj, err
		}
	}

	logObj.FileName = params.BuildId

	return logObj, nil
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
