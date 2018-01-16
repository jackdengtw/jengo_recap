package service

import (
	"fmt"
	"time"

	"github.com/golang/glog"
	"github.com/qetuantuan/jengo_recap/constant"
	"github.com/qetuantuan/jengo_recap/dao"
	"github.com/qetuantuan/jengo_recap/model"
	"github.com/qetuantuan/jengo_recap/queue"
	"github.com/qetuantuan/jengo_recap/task"
)

type BuildReader interface {
	ListBuilds(p *model.EngineListBuildsParams) (model.Builds, error)
	GetBuild(buildId string) (model.Build, error)
}

type BuildCreator interface {
	CreateBuild(params *model.EngineCreateBuildParams) (buildId string, err error)
}

type EngineBuildService interface {
	BuildCreator
	BuildReader
}

type LocalEngineBuildService struct {
	BuildDao dao.EngineBuildDao
	Queue    queue.TaskQueue
}

func (s *LocalEngineBuildService) CreateBuild(params *model.EngineCreateBuildParams) (string, error) {
	glog.Infof("CreateBuildParams: %v, UserId: %s", params, params.UserId)
	var buildId string

	err := s.verifyEngineCreateBuildParams(params)
	if err != nil {
		glog.Errorf("EngineCreateBuildParams: %v, error: err(%s)\n", params, err.Error())
		return "", err
	}

	now := time.Now().UTC()
	var Build = &model.Build{
		RepoId:     params.RepoId,
		EventId:    &params.EventId,
		Status:     constant.BUILD_PRESTART,
		CreatedAt:  &now,
		UpdatedAt:  &now,
		Branch:     params.Branch,
		Commits:    params.Commits,
		HeadCommit: params.HeadCommit,
		UserId:     params.UserId,
	}
	glog.Infof("model.Build: %v", Build)
	buildId, err = s.BuildDao.InsertBuild(*Build)
	if err != nil {
		glog.Errorf("Failed to update Build status as finished. Error is %v", err)
		return "", err
	}

	tsk := task.General{
		Version: 1,
		Build:   Build,
		Status:  constant.BUILD_PRESTART,
	}
	glog.Infof("task: %v", tsk)
	if err = s.Queue.WriteMsgs(
		queue.TopicStartGroup1Name,
		[]interface{}{
			tsk,
		}); err != nil {
		glog.Errorf("Write parse msg failed: %v", err)
	}
	return buildId, err
}

func (s *LocalEngineBuildService) ListBuilds(p *model.EngineListBuildsParams) (model.Builds, error) {
	query := make(map[string]interface{})

	if p.RepoId != "" {
		query["Build.Repoid"] = p.RepoId
	}

	if p.UserId != "" {
		query["Build.userid"] = p.UserId
	}

	if p.BuildId != "" {
		query["Build._id"] = p.BuildId
	}

	if p.EventId != "" {
		query["Build.eventid"] = p.EventId
	}

	Builds, err := s.BuildDao.ListBuilds(query, p.Limit, p.Offset)
	if err != nil {
		glog.Errorf("List builds err: %v", err)
		return nil, err
	}

	return Builds, nil
}

func (s *LocalEngineBuildService) GetBuild(buildId string) (model.Build, error) {
	build, err := s.BuildDao.GetBuild(buildId)
	if err != nil {
		glog.Errorf("Get Build error: %v", err)
		return model.Build{}, err
	}
	return build, nil
}

func (s *LocalEngineBuildService) verifyEngineCreateBuildParams(e *model.EngineCreateBuildParams) (err error) {
	if e == nil || e.Repo == nil || e.Repo.ID == nil || e.Commits == nil { //|| e.HeadCommit == nil
		err = fmt.Errorf("nil error: either e is nil or e.Repo is nil or e.Repo.ID or e.Commits or e.HeadCommit is nil")
		return
	}
	if e.HeadCommit == nil {
		e.HeadCommit = &e.Commits[0]
	}
	return
}
