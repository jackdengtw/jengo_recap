package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/qetuantuan/jengo_recap/api"
	"github.com/qetuantuan/jengo_recap/service"
)

type mockProjectService struct {
}

func (u *mockProjectService) UpdateProjects(userId string) (projects []api.Project, err error) {
	if userId != "u_github_123" {
		err = errors.New("No such User,update faild")
		return
	}
	return nil, nil
}

func (u *mockProjectService) GetProject(projectId string) (project api.Project, err error) {
	if projectId == "123" {
		project = api.Project{Meta: api.ProjectMeta{Id: "123"}}
		return
	}
	err = service.NotFoundError
	return
}

func (u *mockProjectService) GetProjectsByFilter(filter map[string]interface{}, limitCount, offset int) (projects []api.Project, err error) {

	return
}
func (u *mockProjectService) SwitchProject(userId, projectId string, enable bool) error {

	return nil
}

func TestProjectHandler_UpdateProjects(t *testing.T) {
	testUpdateProjects(
		t,
		&mockProjectService{},
		http.StatusOK,
	)

}

func testUpdateProjects(t *testing.T, service service.ProjectServiceInterface, expStatus int) {
	h := NewUpdateProjectHandler(service)

	expectedMap := map[string]int{
		"/v0.1/user/u_github_123/projects/action?method=update": 200,
		"/v0.1/user/u_github_456/projects/action?method=update": 500,
		"/v0.1/user/u_github_123/projects/action?method=delete": 400,
	}

	for url, expStatus := range expectedMap {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(h.Method(), url, nil)

		m := mux.NewRouter()
		q := "{" + h.Query() + "}"
		m.HandleFunc(h.Pattern(), h.UpdateProject).Methods(h.Method()).Queries(h.Query(), q)
		m.ServeHTTP(w, r)
		resp := w.Result()
		if resp.StatusCode != expStatus {
			t.Error(" expected %d, actual %d.\n and the resp is \n", expStatus, resp.StatusCode, resp)

		}
	}
}
func TestProjectHandler_GetProject(t *testing.T) {
	testGetProject(
		t,
		&mockProjectService{},
		http.StatusOK,
	)
}

func testGetProject(t *testing.T, service service.ProjectServiceInterface, expStatus int) {
	h := NewGetProjectHandler(service)

	expectedMap := map[string]int{
		"/v0.1/user/u_github_123/project/123": 200,
		"/v0.1/user/u_github_123/project/456": 404,
	}

	for url, expStatus := range expectedMap {
		if h.Query() != "" {
			t.Error("nil")
		}
		req, _ := http.NewRequest(h.Method(), url, nil)
		w := httptest.NewRecorder()
		m := mux.NewRouter()
		m.HandleFunc(h.Pattern(), h.GetProject).Methods(h.Method())
		m.ServeHTTP(w, req)

		resp := w.Result()
		if resp.StatusCode != expStatus {
			t.Error("expected %d, actual %d.\n", expStatus, resp.StatusCode)
		}

	}
}

func TestProjectHandler_GetProjects(t *testing.T) {
	testGetProjects(
		t,
		&mockProjectService{},
		http.StatusOK,
	)
}

func testGetProjects(t *testing.T, service service.ProjectServiceInterface, expStatus int) {
	h := NewGetProjectsHandler(service)

	expectedMap := map[string]int{
		"/v0.1/user/u_github_123/projects?max_count=10&offset=1&group_by=scm":             200,
		"/v0.1/user/u_github_123/projects?max_count=10&offset=1&group_by=other":           400,
		"/v0.1/user/u_github_123/projects?max_count=10&offset=1&group_by=scm&enable=true": 200,
	}

	for url, expStatus := range expectedMap {
		req, _ := http.NewRequest(h.Method(), url, nil)
		w := httptest.NewRecorder()
		m := mux.NewRouter()
		m.HandleFunc(h.Pattern(), h.GetProjects).Methods(h.Method())
		m.ServeHTTP(w, req)

		resp := w.Result()
		if resp.StatusCode != expStatus {
			t.Error("expected %d, actual %d ,resp is.\n", expStatus, resp.StatusCode)
		}

	}
}

func TestProjectHandler_SwitchProject(t *testing.T) {
	testSwitchProject(
		t,
		&mockProjectService{},
		http.StatusOK,
	)

}

func testSwitchProject(t *testing.T, service service.ProjectServiceInterface, expStatus int) {
	h := NewSwitchProjectHandler(service)

	expectedMap := map[string]int{
		"/v0.1/user/u_github_123/project/123/action?method=enable":  200,
		"/v0.1/user/u_github_456/project/123/action?method=disable": 200,
		"/v0.1/user/u_github_123/project/123/action?method=delete":  400,
	}

	for url, expStatus := range expectedMap {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(h.Method(), url, nil)

		m := mux.NewRouter()
		q := "{" + h.Query() + "}"
		m.HandleFunc(h.Pattern(), h.SwitchProject).Methods(h.Method()).Queries(h.Query(), q)
		m.ServeHTTP(w, r)
		resp := w.Result()
		if resp.StatusCode != expStatus {
			t.Errorf(" expected %d, actual %d.\n and the resp is \n", expStatus, resp.StatusCode, resp)

		}
	}
}
