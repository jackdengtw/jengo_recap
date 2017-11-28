package handler

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/qetuantuan/jengo_recap/api"
	"github.com/qetuantuan/jengo_recap/service"
)

type mockRunLogService struct {
}

func (u *mockRunLogService) PutLog(logs []byte, buildId, runId string) (id string, err error) {
	return "", nil
}

func (u *mockRunLogService) GetLog(id string) (log api.RunLog, err error) {
	if id == "123" {
		log = api.RunLog{Id: "121"}
		return
	} else {
		err = service.NotFoundError
	}
	return
}
func TestPutLogHandler_PutLog(t *testing.T) {
	testPutLog(
		t,
		&mockRunLogService{},
		url.Values{
			"content": {"content"},
		},
		http.StatusOK,
	)

}

func testPutLog(t *testing.T, service service.RunLogServiceInterface, form url.Values, expStatus int) {
	h := NewPutLogHandler(service)

	expectedMap := map[string]int{
		"/v0.1/user/u_github_123/project/123/build/1/run/123/log":  200,
		"/v0.1/user/u_github_123/project/123/build/2/run/1234/log": 200,
		"/v0.1/user/u_github_123/project/123/run/1235/log":         404,
	}

	for url, expStatus := range expectedMap {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(h.Method(), url, nil)
		r.Body = ioutil.NopCloser(strings.NewReader(form.Encode()))
		m := mux.NewRouter()
		m.HandleFunc(h.Pattern(), h.PutLog).Methods(h.Method())
		m.ServeHTTP(w, r)
		resp := w.Result()
		if resp.StatusCode != expStatus {
			t.Error(" expected %d, actual %d.\n and the resp is \n", expStatus, resp.StatusCode, resp)

		}
	}
}
func TestGetLogHandler_GetLog(t *testing.T) {
	testGetLog(
		t,
		&mockRunLogService{},
		http.StatusOK,
	)
}

func testGetLog(t *testing.T, service service.RunLogServiceInterface, expStatus int) {
	h := NewGetLogHandler(service)

	expectedMap := map[string]int{
		"/v0.1/log/123": 200,
		"/v0.1/log/456": 404,
	}

	for url, expStatus := range expectedMap {
		req, _ := http.NewRequest(h.Method(), url, nil)
		w := httptest.NewRecorder()
		m := mux.NewRouter()
		m.HandleFunc(h.Pattern(), h.GetLog).Methods(h.Method())
		m.ServeHTTP(w, req)

		resp := w.Result()
		if resp.StatusCode != expStatus {
			t.Error("expected %d, actual %d.\n", expStatus, resp.StatusCode)
		}

	}
}
