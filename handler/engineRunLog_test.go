package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"errors"

	"github.com/qetuantuan/jengo_recap/api"
	"github.com/qetuantuan/jengo_recap/constant"
	"github.com/qetuantuan/jengo_recap/service"
)

type mockEngineRunLogService struct {
}

func (m *mockEngineRunLogService) GetRunLog(params *api.GetRunLogParams) (*api.RunLog, error) {
	return &api.RunLog{}, nil
}

func (m *mockEngineRunLogService) PutLog(logs []byte, buildId, runId string) (id string, err error) {
	return "", nil
}

func (m *mockEngineRunLogService) GetLog(id string) (log api.RunLog, err error) {
	return
}

type mockInvalidRunLogService struct {
}

func (m *mockInvalidRunLogService) GetRunLog(params *api.GetRunLogParams) (*api.RunLog, error) {
	return &api.RunLog{}, errors.New("this is an error")
}

func (m *mockInvalidRunLogService) PutLog(logs []byte, buildId, runId string) (id string, err error) {
	return "", nil
}

func (m *mockInvalidRunLogService) GetLog(id string) (log api.RunLog, err error) {
	return
}

func TestRunLogHandler_GetRunLog(t *testing.T) {
	testGetRunLog(
		t,
		&mockEngineRunLogService{},
		map[string][]string{
			"offset": []string{"100"},
			"limit":  []string{"12"}},
		map[string][]string{},
		http.StatusOK, // normal
	)
	testGetRunLog(
		t,
		&mockEngineRunLogService{},
		map[string][]string{
			"offset": []string{"100"},
		},
		map[string][]string{},
		http.StatusBadRequest, // no limit
	)
	testGetRunLog(
		t,
		&mockEngineRunLogService{},
		map[string][]string{
			"limit": []string{"12"},
		},
		map[string][]string{},
		http.StatusBadRequest, // no offset
	)
	testGetRunLog(
		t,
		&mockEngineRunLogService{},
		map[string][]string{
			"offset": []string{"a"},
			"limit":  []string{"12"},
		},
		map[string][]string{},
		http.StatusBadRequest, // invalid offset
	)
	testGetRunLog(
		t,
		&mockInvalidRunLogService{},
		map[string][]string{
			"offset": []string{"100"},
			"limit":  []string{"12"},
		},
		map[string][]string{},
		http.StatusInternalServerError, // invalid RunLogService
	)
	testGetRunLog(
		t,
		&mockEngineRunLogService{},
		map[string][]string{
			"offset": []string{"100"},
			"limit":  []string{"12"},
		},
		map[string][]string{
			constant.ACCEPT_ENCODING: []string{"invalid"},
		},
		http.StatusNotAcceptable, // invalid request header
	)
	// TODO: mock GzipEncode
	// testGetRunLog(t, &mockEngineRunLogService{}, map[string][]string{"offset": []string{"100"}, "limit": []string{"12"}}, map[string][]string{constant.ACCEPT_ENCODING: []string{"gzip"}}, http.StatusInternalServerError)
}

func testGetRunLog(t *testing.T, service service.RunLogServiceInterface, form map[string][]string, headers map[string][]string, expStatus int) {
	h := NewGetRunLogHandler(service)

	fmt.Printf("[testGetRunLog] expected_status_code: %d, form: %s \n", expStatus, form)
	w := httptest.NewRecorder()
	r := &http.Request{Form: form, Header: headers}
	h.GetRunLog(w, r)
	resp := w.Result()
	if resp.StatusCode != expStatus {
		t.Error(fmt.Sprintf("[testGetRunLog] expected %d, actual %d.\n", expStatus, resp.StatusCode))
	}
}
