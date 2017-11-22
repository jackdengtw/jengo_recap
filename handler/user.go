package handler

import (
	"net/http"
	"strings"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/qetuantuan/jengo_recap/api"
	"github.com/qetuantuan/jengo_recap/service"
	"github.com/qetuantuan/jengo_recap/util"
)

type CreateUserHandler struct {
	BaseHandler
	Service service.UserServiceInterface
}

func NewCreateUserHandler(service service.UserServiceInterface) *CreateUserHandler {
	h := &CreateUserHandler{
		BaseHandler: BaseHandler{
			name:    "create_user",
			method:  "POST",
			pattern: "/v0.2/user",
		},
		Service: service,
	}
	h.handlerFunc = h.CreateUser
	Register(h)
	return h
}

func processCreateUserParams(r *http.Request) (string, string, string, string) {
	loginName := r.FormValue("login_name")
	userToken := r.FormValue("token")
	auth := r.FormValue("auth")
	scm := r.FormValue("scm")
	glog.Info("user login name is ", loginName)
	var token string
	if len(userToken) < 5 {
		token = userToken
	} else {
		token = userToken[:5]
	}
	glog.Info("user token is ", token, "***")
	glog.Infof("auth and scm are %s %s", auth, scm)
	return loginName, userToken, auth, scm
}

// CreateUser implements the same API of createUser but v0.2
// valid auth:
//   * github.com
/*
params:
login_name: tutuwho520
token: github_token
auth: user_auth
scm: place_holder
------------------------------------
response:
Msg:       success or others info
id:	       Value uniquely identifying the user.
*/
func (h *CreateUserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	loginName, userToken, auth, _ := processCreateUserParams(r)
	// FIXME: only support github.com
	if !strings.EqualFold(auth, api.AUTH_SOURCE_GITHUB) {
		util.ReturnError(http.StatusBadRequest, "Unsupported auth: "+auth, w)
		return
	}
	id, err := h.Service.CreateUser(loginName, auth, userToken)
	if err == nil {
		util.ReturnSuccessWithObj(http.StatusCreated, util.CommonResponse{Id: id}, w)
		return
	} else if err == service.CreateConflictError {
		util.ReturnError(http.StatusConflict, err.Error(), w)
		return
	} else if err == service.NotSupportedAuthError {
		util.ReturnError(http.StatusBadRequest, err.Error(), w)
		return
	} else {
		util.ReturnError(http.StatusInternalServerError, err.Error(), w)
		return
	}
}

type GetUserHandler struct {
	BaseHandler
	Service service.UserServiceInterface
}

func NewGetUserHandler(service service.UserServiceInterface) *GetUserHandler {
	h := &GetUserHandler{
		BaseHandler: BaseHandler{
			name:    "get_user",
			method:  "GET",
			pattern: "/v0.2/internal_user/{user}",
		},
		Service: service,
	}
	h.handlerFunc = h.GetUser
	Register(h)
	return h
}

// getUser implements the same API of getUser but v0.2
/*
params:
user:userId
------------------------------------
response:
user_id:	Value uniquely identifying the user.
login_name:		Login set on scm.
display_name:	Name set on scm.
email :         user email
token:          Scm token
local:          location from scm
avatar_url:		Avatar URL set on scm.
is_syncing:		Whether or not the user is currently being synced with scm.
synced_at:		The last time the user was synced with scm.
created_at :    user created time
*/

func (h *GetUserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId := params["user"]
	glog.Info("user is", userId)

	user, err := h.Service.GetUser(userId)
	if err == nil {
		glog.Info("user found")
		util.ReturnSuccessWithObj(http.StatusOK, user, w)
		return
	} else if err == service.UserNotFoundError {
		util.ReturnError(http.StatusNotFound, err.Error(), w)
		return
	} else {
		util.ReturnError(http.StatusInternalServerError, err.Error(), w)
		return
	}
}

type GetUserByLoginHandler struct {
	BaseHandler
	Service service.UserServiceInterface
}

func NewGetUserByLoginHandler(service service.UserServiceInterface) *GetUserByLoginHandler {
	h := &GetUserByLoginHandler{
		BaseHandler: BaseHandler{
			name:    "getUserByLogin",
			method:  "GET",
			pattern: "/v0.2/internal_user",
		},
		Service: service,
	}
	h.handlerFunc = h.GetUserByLogin
	Register(h)
	return h
}

// getUserByLogin implements the same API of getUserByLogin but v0.2
/*
getUserByLogin
params:
  login_name: login name of auth_source
  auth: enum {"github.com"}
------------------------------------
response:
  the same as getUser
*/

func (h *GetUserByLoginHandler) GetUserByLogin(w http.ResponseWriter, r *http.Request) {
	loginName := r.FormValue("login_name")
	auth := r.FormValue("auth")
	if loginName == "" || !strings.EqualFold(auth, api.AUTH_SOURCE_GITHUB) {
		glog.Warningf("Invalid login name or auth: login_name=%v, auth=%v", loginName, auth)
		util.ReturnError(http.StatusBadRequest, "Invalid login name or auth", w)
		return
	}
	glog.Infof("user is %v/%v", auth, loginName)

	user, err := h.Service.GetUserByLogin(loginName, auth)
	if err == nil {
		util.ReturnSuccessWithObj(http.StatusOK, user, w)
		return
	} else if err == service.UserNotFoundError {
		util.ReturnError(http.StatusNotFound, err.Error(), w)
		return
	} else {
		util.ReturnError(http.StatusInternalServerError, err.Error(), w)
		return
	}

	glog.Info("get user by login success!")
}

// token can be updated by user actively
// other scm properties should be synced by service itself
type UpdateScmTokenHandler struct {
	BaseHandler
	Service service.UserServiceInterface
}

func NewUpdateScmTokenHandler(service service.UserServiceInterface) *UpdateScmTokenHandler {
	h := &UpdateScmTokenHandler{
		BaseHandler: BaseHandler{
			name:    "UpdateScmToken",
			method:  "POST",
			pattern: "/v0.2/user/{user_id}/scm/{scm_id}/token",
		},
		Service: service,
	}
	h.handlerFunc = h.UpdateScmToken
	Register(h)
	return h
}

// UpdateScmToken update token for scm of specific user
func (h *UpdateScmTokenHandler) UpdateScmToken(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.FormValue("token")
	params := mux.Vars(r)
	userId := params["user_id"]
	scmId := params["scm_id"]
	if tokenStr == "" || userId == "" || scmId == "" {
		glog.Warningf("userId/scmId/token should NOT be empty!")
		util.ReturnError(http.StatusBadRequest, "Empty userId/scmId/token", w)
		return
	}
	glog.Infof("update token for %v/%v", userId, scmId)

	err := h.Service.UpdateScmToken(userId, scmId, tokenStr)
	if err == service.UserNotFoundError {
		util.ReturnError(http.StatusNotFound, "User not found: "+err.Error(), w)
		return
	} else if err == service.ScmNotFoundError {
		util.ReturnError(http.StatusNotFound, "Scm not found: "+err.Error(), w)
		return
	} else if err != nil {
		util.ReturnError(http.StatusInternalServerError, err.Error(), w)
		return
	}
	util.ReturnSuccessWithMsg(http.StatusOK, "Token updated", w)
	glog.Info("update scm token success!")
	return
}

// Reserve
// 1. /v?.?/user/$user_id
// 2. /v?.?/user
//   for external api without token ...
