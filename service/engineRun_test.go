package service

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"bytes"
	"encoding/json"

	"github.com/qetuantuan/jengo_recap/api"
	"github.com/qetuantuan/jengo_recap/dao"
	"github.com/qetuantuan/jengo_recap/queue"
)

type mockProjectStoreClient struct {
}

func (m *mockProjectStoreClient) InsertRun(run *api.Run) (*http.Response, []byte, error) {
	body, _ := json.Marshal(map[string]string{"build_id": "1234"})
	resp := &http.Response{
		Status:        "200 OK",
		StatusCode:    200,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Body:          ioutil.NopCloser(bytes.NewBufferString(string(body))),
		ContentLength: int64(len(body)),
		Header:        make(http.Header, 0),
	}

	return resp, body, nil
}

func (m *mockProjectStoreClient) PutLog(userId string, projectId string, buildId string, runId string, runLog string) (*http.Response, []byte, error) {
	return &http.Response{}, nil, nil
}

func (m *mockProjectStoreClient) UpdateRun(run *api.Run, buildId string) (*http.Response, []byte, error) {
	return &http.Response{}, nil, nil
}

func (p *mockProjectStoreClient) UpdatePatchRun(run *api.PatchRun, buildId string) (*http.Response, []byte, error) {
	return &http.Response{}, nil, nil
}

func (p *mockProjectStoreClient) GetProject(userId string, projectId string) (project api.Project, err error) {
	return
}

func TestEngineRunService_CreateRun(t *testing.T) {
	d, server, session := SetupMongo()
	defer TearDown(d, server, session)

	s := &EngineRunService{&mockProjectStoreClient{}, &dao.RunDao{GSession: session}, queue.NewNativeTaskQueue()}

	// test CreateRun
	commit := api.PushEventCommit{}
	repoId := 123
	params := &api.EngineCreateRunParams{
		UserId:    "user_id",
		ProjectId: "project_id",
		Repo:      &api.PushEventRepository{ID: &repoId},
		Commits:   []api.PushEventCommit{commit},
	}
	run, err := s.CreateRun(params)
	if err != nil {
		t.Error(fmt.Sprintf("[TestEngineRunService_CreateRun][CreateRun] error in CreateRun"))
	}
	if !checkRunIsSame(run, params) {
		t.Error(fmt.Sprintf("[TestEngineRunService_CreateRun][CreateRun] the run is not the same"))
	}

	dParams := &api.EngineDescribeRunsParams{
		UserId:    run.UserId,
		ProjectId: run.EventId,
	}

	// test DescribeRuns
	runsAct, err := s.DescribeRuns(dParams)
	if err != nil {
		t.Error(fmt.Sprintf("[TestEngineRunService_CreateRun][DescribeRuns] error in DescribeRuns"))
	}
	if len(runsAct) != 1 {
		t.Error(fmt.Printf("[TestEngineRunService_CreateRun][DescribeRuns] the length of runs is not 1: %v", runsAct))
	}
	if !checkRunIsSame(&(runsAct[0]), params) {
		t.Error(fmt.Printf("[TestEngineRunService_CreateRun][DescribeRuns] the user_id is %s, the project_id is %s", runsAct[0].UserId, runsAct[0].ProjectId))
	}

	// test DescribeRun
	runAct, err := s.DescribeRun(run.Id)
	if err != nil {
		t.Error(fmt.Sprintf("[TestEngineRunService_CreateRun][DescribeRun] error in DescribeRun"))
	}
	if !checkRunIsSame(runAct, params) {
		t.Error(fmt.Printf("[TestEngineRunService_CreateRun][DescribeRun] the user_id is %s, the project_id is %s", runAct.UserId, runAct.ProjectId))
	}
}

func checkRunIsSame(run *api.Run, params *api.EngineCreateRunParams) bool {
	if run.ProjectId != params.ProjectId {
		return false
	} else if run.UserId != params.UserId {
		return false
	} else {
		return true
	}
}

func TestEngineRunService_DescribeRuns(t *testing.T) {
	d, server, session := SetupMongo()
	defer TearDown(d, server, session)

	s := &EngineRunService{&mockProjectStoreClient{}, &dao.RunDao{GSession: session}, queue.NewNativeTaskQueue()}
	p := &api.EngineDescribeRunsParams{
		ProjectId: "p_123456789_123456789",
	}
	runs, err := s.DescribeRuns(p)
	if err != nil {
		t.Error(fmt.Sprintf("[TestEngineRunService_DescribeRuns] error in DescribeRuns"))
	}
	if len(runs) != 0 {
		t.Error(fmt.Sprintf("[TestEngineRunService_DescribeRuns] the length of runs is not 0"))
	}
}

func TestEngineRunService_DescribeRun(t *testing.T) {
	d, server, session := SetupMongo()
	defer TearDown(d, server, session)

	s := &EngineRunService{&mockProjectStoreClient{}, &dao.RunDao{GSession: session}, queue.NewNativeTaskQueue()}
	run, err := s.DescribeRun("run_123456789_123456789")
	if err == nil {
		t.Error(fmt.Sprintf("[TestEngineRunService_DescribeRun] No error in DescribeRun"))
	}
	if run != nil {
		t.Error(fmt.Sprintf("[TestEngineRunService_DescribeRun] the run is not nil"))
	}
}
