package dao

import (
	"github.com/golang/glog"
	"github.com/qetuantuan/jengo_recap/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var dbName string = "engine"
var runCol string = "run"

type RunDaoInterface interface {
	UpdateRuns(model.InnerRuns) ([]string, []error)
	UpdateRun(*model.InnerRun) error
	UpdateRunProperties(string, map[string]interface{}) error
	InsertRuns(model.InnerRuns) ([]string, []error)
	InsertRun(*model.InnerRun) error
	DeleteRuns(model.InnerRuns) ([]string, []error)
	DeleteRun(*model.InnerRun) error
	DescribeRuns(map[string]interface{}, int, int) (model.InnerRuns, error)
	DescribeRun(string) (*model.InnerRun, error)
}

type RunDao struct {
	Url      string
	GSession *mgo.Session
}

func (md *RunDao) Init() (err error) {
	md.GSession, err = mgo.Dial(md.Url)
	var l MgoLog
	mgo.SetLogger(l)
	// TODO: make debug a command line option
	// This is to debug Mongo queries
	mgo.SetDebug(false)
	return err
}

func (md *RunDao) UpdateRuns(runs model.InnerRuns) (failedRunIds []string, err []error) {
	session := md.GSession.Copy()
	defer session.Close()
	pc := session.DB(dbName).C(runCol)
	for _, run := range runs {
		tmp_err := pc.UpdateId(run.HashId, bson.M{"$set": *run})
		if tmp_err != nil {
			failedRunIds = append(failedRunIds, run.Id)
			err = append(err, tmp_err)
			glog.Errorf("UpsertId failed for %v %v %v", run.Id, run.Id, err)
		}
	}
	return
}

func (md *RunDao) UpdateRun(run *model.InnerRun) (err error) {
	session := md.GSession.Copy()
	defer session.Close()
	pc := session.DB(dbName).C(runCol)
	err = pc.UpdateId(run.HashId, bson.M{"$set": *run})
	if err != nil {
		glog.Errorf("UpsertId failed for %v %v %v", run.Id, run.HashId, err)
	}
	return
}

func (md *RunDao) UpdateRunProperties(runId string, updateData map[string]interface{}) (err error) {
	session := md.GSession.Copy()
	defer session.Close()
	pc := session.DB(dbName).C(runCol)
	err = pc.Update(bson.M{"run._id": runId}, bson.M{"$set": updateData})
	if err != nil {
		glog.Errorf("Update failed for %s %v", runId, err)
	}
	return
}

func (md *RunDao) InsertRuns(runs model.InnerRuns) (failedRunIds []string, err []error) {
	session := md.GSession.Copy()
	defer session.Close()
	pc := session.DB(dbName).C(runCol)
	for _, run := range runs {
		tmp_err := pc.Insert(*run)
		if tmp_err != nil {
			failedRunIds = append(failedRunIds, run.Id)
			err = append(err, tmp_err)
			glog.Errorf("UpsertId failed for %v %v %v", run.Id, run.HashId, err)
		}
	}
	return
}

func (md *RunDao) InsertRun(run *model.InnerRun) (err error) {
	session := md.GSession.Copy()
	defer session.Close()
	pc := session.DB(dbName).C(runCol)

	err = pc.Insert(*run)
	if err != nil {
		glog.Errorf("UpsertId failed for %v %v %v", run.Id, run.HashId, err)

	}
	return
}

func (md *RunDao) DeleteRuns(runs model.InnerRuns) (failedRunIds []string, err []error) {
	session := md.GSession.Copy()
	defer session.Close()
	pc := session.DB(dbName).C(runCol)
	for _, run := range runs {
		tmp_err := pc.RemoveId(run.HashId)
		if tmp_err != nil {
			failedRunIds = append(failedRunIds, run.Id)
			err = append(err, tmp_err)
			glog.Errorf("DeleteRuns failed for %v %v %v", run.Id, run.HashId, err)
		}
	}
	return
}

func (md *RunDao) DeleteRun(run *model.InnerRun) (err error) {
	session := md.GSession.Copy()
	defer session.Close()
	pc := session.DB(dbName).C(runCol)
	err = pc.RemoveId([]byte(run.HashId))
	if err != nil {
		glog.Errorf("DeleteRuns failed for %v %v %v", run.Id, run.HashId, err)
	}
	return
}

func (md *RunDao) DescribeRuns(query map[string]interface{}, limitCount int, offset int) (runs model.InnerRuns, err error) {
	session := md.GSession.Copy()
	defer session.Close()
	pc := session.DB(dbName).C(runCol)
	var runCs []model.InnerRun
	err = pc.Find(query).Skip(offset).Limit(limitCount).All(&runCs)
	for _, r := range runCs {
		runs = append(runs, &r)
	}
	return
}

func (md *RunDao) DescribeRun(id string) (run *model.InnerRun, err error) {
	session := md.GSession.Copy()
	defer session.Close()
	pc := session.DB(dbName).C(runCol)
	var r model.InnerRun
	err = pc.FindId([]byte(id)).One(&r)
	return &r, err
}
