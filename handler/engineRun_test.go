package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"errors"

	"github.com/qetuantuan/jengo_recap/api"
	"github.com/qetuantuan/jengo_recap/service"
)

type mockEngineRunService struct {
}

func (m *mockEngineRunService) CreateRun(params *api.EngineCreateRunParams) (*api.Run, error) {
	return &api.Run{Id: "id-123"}, nil
}

func (m *mockEngineRunService) DescribeRuns(p *api.EngineDescribeRunsParams) (api.Runs, error) {
	return nil, nil
}

func (m *mockEngineRunService) DescribeRun(runId string) (*api.Run, error) {
	return &api.Run{Id: "id-123"}, nil
}

func (m *mockEngineRunService) verifyEngineCreateRunParams(e *api.EngineCreateRunParams) (err error) {
	return nil
}

type mockInvalidRunService struct {
}

func (m *mockInvalidRunService) CreateRun(params *api.EngineCreateRunParams) (*api.Run, error) {
	return nil, errors.New("this is an error")
}

func (m *mockInvalidRunService) DescribeRuns(p *api.EngineDescribeRunsParams) (api.Runs, error) {
	return nil, errors.New("this is an error")
}

func (m *mockInvalidRunService) DescribeRun(runId string) (*api.Run, error) {
	return nil, errors.New("this is an error")
}

func (m *mockInvalidRunService) verifyEngineCreateRunParams(e *api.EngineCreateRunParams) (err error) {
	return nil
}

func TestRunHandler_CreateRun(t *testing.T) {
	p := api.EngineCreateRunParams{
		ProjectId: "project_id",
		UserId:    "user_id",
		EventId:   "event_id",
		Branch:    "branch",
		Compare:   "compare",
	}

	pStr, _ := json.Marshal(p)
	fmt.Printf(string(pStr))
	testMap := map[int](map[string]string){
		400: {"not a json": "err_msg"},
		200: {string(pStr): "ok_msg"},
	}

	h := NewCreateRunHandler(&mockEngineRunService{})

	for status, bodyRespMap := range testMap {
		for body, _ := range bodyRespMap {
			fmt.Printf("[TestRunHandler_CreateRun] expected_status_code: %d, body: %s \n", status, body)
			w := httptest.NewRecorder()
			r := &http.Request{Body: ioutil.NopCloser(strings.NewReader(body))}
			h.CreateRun(w, r)
			resp := w.Result()
			if resp.StatusCode != status {
				t.Error(fmt.Sprintf("[TestRunHandler_CreateRun] expected %d, actual %d.\n", status, resp.StatusCode))
			}
		}
	}
}

func TestRunHandler_DescribeRuns(t *testing.T) {
	testDescribeRuns(
		t,
		&mockEngineRunService{},
		map[string][]string{
			"project_id": []string{"project_id"},
			"run_id":     []string{"run_id"}},
		http.StatusOK,
	)
}

func testDescribeRuns(t *testing.T, service service.EngineRunServiceInterface, form map[string][]string, expStatus int) {
	h := NewDescribeRunsHandler(service)

	fmt.Printf("[testDescribeRuns] expected_status_code: %d, form: %s \n", expStatus, form)
	w := httptest.NewRecorder()
	r := &http.Request{Form: form}
	h.DescribeRuns(w, r)
	resp := w.Result()
	if resp.StatusCode != expStatus {
		t.Error(fmt.Sprintf("[testDescribeRuns] expected %d, actual %d.\n", expStatus, resp.StatusCode))
	}
}

func TestRunHandler_DescribeRun(t *testing.T) {
	testDescribeRun(
		t,
		&mockEngineRunService{},
		"run_ut_id",
		http.StatusOK,
	)

	testDescribeRun(
		t,
		&mockInvalidRunService{},
		"run_ut_id",
		http.StatusInternalServerError,
	)
}

func testDescribeRun(t *testing.T, service service.EngineRunServiceInterface, runId string, expStatus int) {
	h := NewDescribeRunHandler(service)

	fmt.Printf("[testDescribeRun] expected_status_code: %d, runId: %s \n", expStatus, runId)
	w := httptest.NewRecorder()
	r := &http.Request{}
	h.DescribeRun(w, r)
	resp := w.Result()
	if resp.StatusCode != expStatus {
		t.Error(fmt.Sprintf("[testDescribeRuns] expected %d, actual %d.\n", expStatus, resp.StatusCode))
	}
}
