package service

import (
	"fmt"
	"time"

	"github.com/golang/glog"
	"github.com/qetuantuan/jengo_recap/api"
	"github.com/qetuantuan/jengo_recap/client"
	"github.com/qetuantuan/jengo_recap/constant"
	"github.com/qetuantuan/jengo_recap/dao"
	"github.com/qetuantuan/jengo_recap/model"
	"github.com/qetuantuan/jengo_recap/queue"
	"github.com/qetuantuan/jengo_recap/task"
	"gopkg.in/mgo.v2/bson"
)

type RunCreator interface {
	CreateRun(params *api.EngineCreateRunParams) (*api.Run, error)
}

type RunDescriber interface {
	DescribeRuns(p *api.EngineDescribeRunsParams) (api.Runs, error)
	DescribeRun(runId string) (*api.Run, error)
}

type EngineRunServiceInterface interface {
	RunCreator
	RunDescriber
}

type EngineRunService struct {
	ProjectClient client.ProjectStoreClientInterface
	RunDao        dao.RunDaoInterface
	Queue         queue.TaskQueue
}

func (s *EngineRunService) CreateRun(params *api.EngineCreateRunParams) (*api.Run, error) {
	glog.Infof("CreateRunParams: %v, UserId: %s", params, params.UserId)

	err := s.verifyEngineCreateRunParams(params)
	if err != nil {
		glog.Errorf("EngineCreateRunParams: %v, error: err(%s)\n", params, err.Error())
		return nil, err
	}
	runId := bson.NewObjectId()
	rawId := runId.Hex()

	now := time.Now().UTC()
	var run = &model.InnerRun{
		Run: api.Run{
			Id:         rawId,
			ProjectId:  params.ProjectId,
			EventId:    params.EventId,
			Status:     constant.RUN_PRESTART,
			CreatedAt:  &now,
			UpdatedAt:  &now,
			Branch:     params.Branch,
			Commits:    params.Commits,
			HeadCommit: params.HeadCommit,
			UserId:     params.UserId,
		},
		Job:    params.ProjectId,
		HashId: []byte(rawId),
	}
	glog.Infof("model.Run: %v", run)
	err = s.RunDao.InsertRun(run)
	if err != nil {
		glog.Errorf("Failed to update run status as finished. Error is %v", err)
		return nil, err
	}

	tsk := task.General{
		Version: 1,
		Run:     run,
		Status:  constant.RUN_PRESTART,
	}
	glog.Infof("task: %v", tsk)
	if err = s.Queue.WriteMsgs(
		queue.TopicStartGroup1Name,
		[]interface{}{
			tsk,
		}); err != nil {
		glog.Errorf("Write parse msg failed: %v", err)
	}
	return &run.Run, err
}

func (s *EngineRunService) DescribeRuns(p *api.EngineDescribeRunsParams) (api.Runs, error) {
	query := make(map[string]interface{})

	if p.ProjectId != "" {
		query["run.projectid"] = p.ProjectId
	}

	if p.UserId != "" {
		query["run.userid"] = p.UserId
	}

	if p.RunId != "" {
		query["run._id"] = p.RunId
	}

	if p.EventId != "" {
		query["run.eventid"] = p.EventId
	}

	runs, err := s.RunDao.DescribeRuns(query, p.Limit, p.Offset)
	if err != nil {
		glog.Errorf("dao.DescribeRuns err: %v", err)
		return []api.Run{}, err
	}

	return runs.ToApiObj(), nil
}

func (s *EngineRunService) DescribeRun(runId string) (*api.Run, error) {
	run, err := s.RunDao.DescribeRun(runId)
	if err != nil {
		glog.Errorf("RunDao.DescribeRun error: %v", err)
		return nil, err
	}
	return &run.Run, nil
}

func (s *EngineRunService) deleteRun(run *model.InnerRun) {
	s.RunDao.DeleteRun(run)
}

func (s *EngineRunService) verifyEngineCreateRunParams(e *api.EngineCreateRunParams) (err error) {
	if e == nil || e.Repo == nil || e.Repo.ID == nil || e.Commits == nil { //|| e.HeadCommit == nil
		err = fmt.Errorf("nil error: either e is nil or e.Repo is nil or e.Repo.ID or e.Commits or e.HeadCommit is nil")
		return
	}
	if e.HeadCommit == nil {
		e.HeadCommit = &e.Commits[0]
	}
	return
}
