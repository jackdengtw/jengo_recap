package service

import (
	"github.com/golang/glog"
	"github.com/qetuantuan/jengo_recap/api"
	"github.com/qetuantuan/jengo_recap/constant"
	"github.com/qetuantuan/jengo_recap/dao"
	"github.com/qetuantuan/jengo_recap/model"
	"github.com/qetuantuan/jengo_recap/scm"
)

type RunServiceInterface interface {
	GetBuildsByFilter(filter map[string]string, maxCount, offset int) (runs api.Builds, err error)
	GetBuildsByIds(buildIds []string) (builds api.Builds, err error)
	UpdateRun(buildId string, run api.Run) (err error)
	PartialUpdateRun(buildId string, run api.Run, runInterface map[string]interface{}) (err error)
	InsertRun(run api.Run) (build api.Build, err error)
}

type RunService struct {
	Md        *dao.ProjectMongoDao
	GithubScm *scm.GithubScm
}

func (r *RunService) GetBuildsByFilter(filter map[string]string, maxCount, offset int) (builds api.Builds, err error) {
	var bs model.Builds
	bs, err = r.Md.GetBuildsByFilter(filter, maxCount, offset)
	if err != nil {
		glog.Error("get project failed! error: ", err)
		err = MongoError
		return
	}
	builds = bs.ToApiObj()
	return
}
func (r *RunService) PartialUpdateRun(buildId string, run api.Run, runInterface map[string]interface{}) (err error) {
	err = r.Md.UpdateRunProperties(buildId, run.Id, runInterface)
	if err != nil {
		glog.Warningf("partial update run failed! error:%v", err)
		err = MongoError
		return
	}

	if run.Status == constant.RUN_FINISHED {
		err = r.Md.UpdateDynamicProjectInfo(run.ProjectId, run.State, "", "")
		if err != nil {
			glog.Warningf("update project:%s state to %s failed! error:%v", run.ProjectId, run.State, err)
			err = MongoError
			return
		}
		glog.Infof("update project:%s state to %s success!", run.ProjectId, run.State)
	}
	return
}
func (r *RunService) UpdateRun(buildId string, run api.Run) (err error) {
	mr := model.NewRunFrom(&run)
	err = r.Md.UpdateRun(buildId, *mr)
	if err != nil {
		glog.Warningf("update run failed! error:%v", err)
		err = MongoError
		return
	}

	if run.Status == constant.RUN_FINISHED {
		err = r.Md.UpdateDynamicProjectInfo(run.ProjectId, run.State, "", "")
		if err != nil {
			glog.Warningf("update project:%s state to %s failed! error:%v", run.ProjectId, run.State, err)
			err = MongoError
			return
		}
		glog.Infof("update project:%s state to %s success!", run.ProjectId, run.State)
	}
	return
}
func (r *RunService) GetBuildsByIds(buildIds []string) (builds api.Builds, err error) {
	var bs model.Builds
	bs, err = r.Md.GetBuilds(buildIds)
	if err != nil {
		glog.Error("get project failed! error: ", err)
		err = MongoError
		return
	}
	builds = bs.ToApiObj()
	return
}

//if there is run's parent build insert to build,otherwise buildNo+1 and create new build
func (r *RunService) InsertRun(run api.Run) (build api.Build, err error) {
	commitId := *run.HeadCommit.ID
	branch := run.Branch
	var b model.Build
	b, err = r.Md.FindBuildByBranchCommit(run.ProjectId, commitId, branch) //find build by commitId+branch+projectId
	if err == dao.ErrorBuildNotFind {
		var buildNo = 1
		var buildId = ""
		buildNo, err = r.Md.GetBuildIndex(run.ProjectId)
		if err != nil {
			glog.Warningf("get build index from project when insert build failed!err:%v", err)
			err = MongoError
			return
		}
		b = model.Build{
			ProjectId: run.ProjectId,
			UserId:    run.UserId,
			Numero:    buildNo,
			Runs:      []api.Run{api.Run(run)},
			CommitId:  commitId,
			Branch:    branch,
		}
		buildId, err = r.Md.InsertBuild(b)
		if err != nil {
			glog.Warningf("insert new build to mongodb failed! error:%v", err)
			err = MongoError
			return
		}
		state := ""
		if run.Status == constant.RUN_FINISHED {
			state = run.State
		}
		err = r.Md.UpdateDynamicProjectInfo(b.ProjectId, state, buildId, branch)
		if err != nil {
			glog.Warningf("update project:%s latest build id :%s failed! error:%v", b.ProjectId, buildId, err)
			err = MongoError
			return
		}
		return
	}

	if err != nil {
		glog.Warningf("find build by commit from mongodb failed! error:%v", err)
		err = MongoError
		return
	}

	mr := model.NewRunFrom(&run)
	err = r.Md.InsertRun(b.Id, *mr)
	if err != nil {
		glog.Warningf("insert run to mongodb failed! error:%v", err)
		err = MongoError
		return
	}
	if run.Status == constant.RUN_FINISHED {
		err = r.Md.UpdateDynamicProjectInfo(run.ProjectId, run.State, "", run.Branch)
		if err != nil {
			glog.Warningf("insert project:%s set state to %s failed! error:%v", run.ProjectId, run.State, err)
			err = MongoError
			return
		}
		glog.Infof("update project:%s state to %s success!", run.ProjectId, run.State)
	}
	return
}
