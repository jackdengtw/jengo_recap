package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/qetuantuan/jengo_recap/model"
	"github.com/qetuantuan/jengo_recap/service"

	"github.com/golang/glog"
	"github.com/gorilla/mux"

	"io/ioutil"
	"net/url"
	"strings"
)

type mockUserService struct {
}

func (u *mockUserService) CreateUser(loginName string, auth string, token string) (userId string, err error) {
	if loginName == "tutuwho520" {
		return "", service.CreateConflictError
	}
	user1 := &model.User{UserId: "user_123"}
	userId = user1.UserId
	return userId, nil
}

func (u *mockUserService) GetUser(userId string) (apiUser *model.ApiUser, err error) {
	glog.Info("in GetUser", userId, "\n")
	if userId == "u_github_123" {
		apiUser = &model.ApiUser{UserId: "user_123"}
	} else {
		err = service.UserNotFoundError
	}
	return
}

func (u *mockUserService) GetUserByLogin(loginName string, auth string) (apiUser *model.ApiUser, err error) {
	if loginName == "tutuwho520" {
		apiUser = &model.ApiUser{UserId: "user_123"}
		return
	} else {
		err = service.UserNotFoundError
	}
	return
}

func (u *mockUserService) UpdateScmToken(userId, scmId, tokenStr string) error {
	return nil
}

func TestUserHandler_CreateUser(t *testing.T) {
	testCreateUser(
		t,
		&mockUserService{},
		url.Values{
			"login_name": {"login_name"},
			"auth":       {"auth"},
			"token":      {"token"}},
		http.StatusOK,
	)

}

func testCreateUser(t *testing.T, service service.UserServiceInterface, form url.Values, expStatus int) {
	h := NewCreateUserHandler(service)

	expectedMap := map[string]int{
		"/v0.2/user?login_name=tutuwho520&auth=github.com&&token=12345": 409,
		"/v0.2/user?login_name=test521&auth=github.com&&token=12345":    201,
		"/v0.2/user?login_name=test521&auth=":                           400,
	}

	for url, expStatus := range expectedMap {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(h.Method(), url, nil)
		r.Body = ioutil.NopCloser(strings.NewReader(form.Encode()))
		m := mux.NewRouter()
		m.HandleFunc(h.Pattern(), h.CreateUser).Methods(h.Method())
		m.ServeHTTP(w, r)
		resp := w.Result()
		if resp.StatusCode != expStatus {
			t.Error(" expected %d, actual %d.\n and the resp is \n", expStatus, resp.StatusCode, resp)

		}
	}
}
func TestUserHandler_GetUser(t *testing.T) {
	testGetUser(
		t,
		&mockUserService{},
		http.StatusOK,
	)
}

func testGetUser(t *testing.T, service service.UserServiceInterface, expStatus int) {
	h := NewGetUserHandler(service)

	expectedMap := map[string]int{
		"/v0.2/internal_user/u_github_123": 200,
		"/v0.2/internal_user/u_github_1":   404,
	}

	for url, expStatus := range expectedMap {
		req, _ := http.NewRequest(h.Method(), url, nil)
		w := httptest.NewRecorder()
		m := mux.NewRouter()
		m.HandleFunc(h.Pattern(), h.GetUser).Methods(h.Method())
		m.ServeHTTP(w, req)

		resp := w.Result()
		if resp.StatusCode != expStatus {
			t.Error("expected %d, actual %d.\n", expStatus, resp.StatusCode)
		}

	}
}

func TestUserHandler_GetUserByLogin(t *testing.T) {
	testGetUserByLogin(
		t,
		&mockUserService{},
		http.StatusOK,
	)
}

func testGetUserByLogin(t *testing.T, service service.UserServiceInterface, expStatus int) {
	h := NewGetUserByLoginHandler(service)

	expectedMap := map[string]int{
		"/v0.2/internal_user?login_name=tutuwho520&auth=github.com": 200,
		"/v0.2/internal_user?login_name=test521&auth=github.com":    404,
		"/v0.2/internal_user?login_name=test521":                    400,
	}

	for url, expStatus := range expectedMap {
		req, _ := http.NewRequest(h.Method(), url, nil)
		w := httptest.NewRecorder()
		m := mux.NewRouter()
		m.HandleFunc(h.Pattern(), h.GetUserByLogin).Methods(h.Method())
		m.ServeHTTP(w, req)

		resp := w.Result()
		if resp.StatusCode != expStatus {
			t.Error("expected %d, actual %d.\n", expStatus, resp.StatusCode)
		}

	}
}
