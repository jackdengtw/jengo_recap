package handler

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/qetuantuan/jengo_recap/api"
	"github.com/qetuantuan/jengo_recap/model"
	"github.com/qetuantuan/jengo_recap/service"
)

type mockRunService struct {
}

func (r *mockRunService) GetBuildsByFilter(filter map[string]string, maxCount, offset int) (runs api.Builds, err error) {
	return
}
func (r *mockRunService) GetBuildsByIds(buildIds []string) (builds api.Builds, err error) {
	if len(buildIds) == 1 && buildIds[0] == "1" {
		b1 := api.Build{Id: "1"}
		builds = append(builds, b1)
		return
	}
	return
}
func (r *mockRunService) UpdateRun(buidId string, run api.Run) (err error) {
	return
}
func (r *mockRunService) PartialUpdateRun(buildId string, run api.Run, runInterface map[string]interface{}) (err error) {
	return
}
func (r *mockRunService) InsertRun(run api.Run) (build api.Build, err error) {
	return
}

func TestGetBuildsByProjectHandler_GetBuildsByProject(t *testing.T) {
	testGetBuildsByProject(
		t,
		&mockRunService{},
		http.StatusOK,
	)

}

func testGetBuildsByProject(t *testing.T, service service.RunServiceInterface, expStatus int) {
	h := NewGetBuildsByFilterHandler(service)

	expectedMap := map[string]int{
		"/v0.1/user/u_github_123/project/123/builds?branch=master": 200,
		"/v0.1/user/u_github_123/project/456/builds":               200,
	}

	for url, expStatus := range expectedMap {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(h.Method(), url, nil)
		m := mux.NewRouter()
		m.HandleFunc(h.Pattern(), h.GetBuildsByFilter).Methods(h.Method())
		m.ServeHTTP(w, r)
		resp := w.Result()

		if resp.StatusCode != expStatus {
			t.Error(" expected %d, actual %d.\n and the resp is \n", expStatus, resp.StatusCode, resp)

		}
	}
}
func TestGetBuildsByIdsHandler_GetBuildsByIds(t *testing.T) {
	testGetBuildsByIds(
		t,
		&mockRunService{},
		url.Values{
			"ids": {"ids"}},
		http.StatusOK,
	)

}

func testGetBuildsByIds(t *testing.T, service service.RunServiceInterface, form url.Values, expStatus int) {
	h := NewGetBuildsByIdsHandler(service)

	expectedMap := map[string]int{
		"/v0.1/builds?ids=1": 200,
		"/v0.1/builds":       404,
	}
	//   /v0.1/builds
	for url, expStatus := range expectedMap {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(h.Method(), url, nil)
		r.Body = ioutil.NopCloser(strings.NewReader(form.Encode()))
		m := mux.NewRouter()
		q := "{" + h.Query() + "}"
		m.HandleFunc(h.Pattern(), h.GetBuildsByIds).Methods(h.Method()).Queries(h.Query(), q)
		m.ServeHTTP(w, r)
		resp := w.Result()
		if resp.StatusCode != expStatus {
			t.Error(" expected %d, actual %d.\n and the resp is \n", expStatus, resp.StatusCode, resp)

		}
	}
}
func TestUpdateRunHandler_UpdateRun(t *testing.T) {
	testUpdateRun(
		t,
		&mockRunService{},
		http.StatusOK,
	)

}

func testUpdateRun(t *testing.T, service service.RunServiceInterface, expStatus int) {
	h := NewUpdateRunHandler(service)
	run := model.Run{
		Id:      "1",
		EventId: "2",
		Branch:  "main"}

	runReq, _ := json.Marshal(run)
	expectedMap := map[string]int{
		"/v0.1/user/u_github_123/project/123/build/1/run/1": 200,
		"/v0.1/user/u_github_1234/project/456/build/1":      404,
	}

	for url, expStatus := range expectedMap {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(h.Method(), url, bytes.NewBuffer([]byte(runReq)))
		//r.Body = ioutil.NopCloser(strings.NewReader(string(bs)))
		m := mux.NewRouter()
		m.HandleFunc(h.Pattern(), h.UpdateRun).Methods(h.Method())
		m.ServeHTTP(w, r)
		resp := w.Result()
		if resp.StatusCode != expStatus {
			t.Error(" expected %d, actual %d.\n and the resp is \n", expStatus, resp.StatusCode, resp)

		}
	}
}
func TestPartialUpdateRunHandler_PartialUpdateRun(t *testing.T) {
	testPartialUpdateRun(
		t,
		&mockRunService{},
		http.StatusOK,
	)

}

func testPartialUpdateRun(t *testing.T, service service.RunServiceInterface, expStatus int) {
	h := NewPartialUpdateRunHandler(service)
	run := model.Run{
		Id:      "1",
		EventId: "2",
		Branch:  "main"}

	runReq, _ := json.Marshal(run)
	expectedMap := map[string]int{
		"/v0.1/user/u_github_123/project/123/build/1/run/1": 200,
		"/v0.1/user/u_github_1234/project/456/build/1":      404,
	}

	for url, expStatus := range expectedMap {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(h.Method(), url, bytes.NewBuffer([]byte(runReq)))
		//r.Body = ioutil.NopCloser(strings.NewReader(string(bs)))
		m := mux.NewRouter()
		m.HandleFunc(h.Pattern(), h.PartialUpdateRun).Methods(h.Method())
		m.ServeHTTP(w, r)
		resp := w.Result()
		if resp.StatusCode != expStatus {
			t.Error(" expected %d, actual %d.\n and the resp is \n", expStatus, resp.StatusCode, resp)

		}
	}
}

func TestInsertRunHandler_InsertRun(t *testing.T) {
	testInsertRun(
		t,
		&mockRunService{},
		url.Values{
			"commits": {"commits"},
			"branch":  {"branch"}},
		http.StatusOK,
	)

}

func testInsertRun(t *testing.T, service service.RunServiceInterface, form url.Values, expStatus int) {
	h := NewInsertRunHandler(service)
	var test = "123445656"
	run := model.Run{
		Id:         "1",
		EventId:    "2",
		Branch:     "main",
		HeadCommit: &api.PushEventCommit{ID: &test},
	}

	runReq, _ := json.Marshal(run)
	expectedMap := map[string]int{
		"/v0.1/user/u_github_123/project/123/run/1": 200,
		"/v0.1/user/u_github_123/project/456/run/1": 200,
		"/v0.1/user/u_github_123/project/123/run/2": 200,
	}

	for url, expStatus := range expectedMap {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(h.Method(), url, bytes.NewBuffer([]byte(runReq)))

		m := mux.NewRouter()
		m.HandleFunc(h.Pattern(), h.InsertRun).Methods(h.Method())
		m.ServeHTTP(w, r)
		resp := w.Result()
		if resp.StatusCode != expStatus {
			t.Error(" expected %d, actual %d.\n and the resp is \n", expStatus, resp.StatusCode, resp)

		}
	}
}
