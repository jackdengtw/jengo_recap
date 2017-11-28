package handler

import (
	"fmt"

	"net/http"
	"strconv"

	"reflect"
	"strings"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/qetuantuan/jengo_recap/api"
	"github.com/qetuantuan/jengo_recap/service"
	"github.com/qetuantuan/jengo_recap/util"

	"errors"
)

//temporary
const (
	VERSION_V1 = "v0.1"
	VERSION_V2 = "v0.2"
)

const ERROR_WRITE_BACK = "Error in writing back response!"

type UpdateProjectHandler struct {
	BaseHandler
	Service service.ProjectServiceInterface
}

func NewUpdateProjectHandler(service service.ProjectServiceInterface) *UpdateProjectHandler {
	h := &UpdateProjectHandler{
		BaseHandler: BaseHandler{
			name:    "update_project",
			method:  "GET",
			pattern: "/v0.1/user/{user_id}/projects/action",
			query:   "method"},
		Service: service,
	}
	h.handlerFunc = h.UpdateProject

	Register(h)
	return h
}

/**
GET /v1.0/user/{user_id}/projects/action?method={method}
Function:
	action of projects, only support update right now
Response:
	project list, project define in https://github.com/qetuantuan/jengo_recap/blob/main/model/project.go
*/
func (h *UpdateProjectHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["user_id"]
	method := vars["method"]
	glog.Info(fmt.Sprintf("begin action:%s for user:%s", method, userId))
	if method == "update" {

		_, err := h.Service.UpdateProjects(userId)
		if err != nil {
			glog.Warningf("update projects failed! user: %v, error: %v", userId, err)
			util.ReturnError(
				http.StatusInternalServerError,
				err.Error(),
				w)
			return
		}
		msg := fmt.Sprintf("update projects success! user_id: %v", userId)
		util.ReturnSuccessWithMsg(
			http.StatusOK,
			msg,
			w)
		glog.Info(msg)
	} else {
		msg := fmt.Sprintf("Method %v is not supported", method)
		glog.Error(msg)
		util.ReturnError(
			http.StatusBadRequest,
			msg,
			w)
	}
	glog.Info(fmt.Sprintf("action:%s for user:%s success", method, userId))
}

type GetProjectHandler struct {
	BaseHandler
	Service service.ProjectServiceInterface
}

func NewGetProjectHandler(service service.ProjectServiceInterface) *GetProjectHandler {
	h := &GetProjectHandler{
		BaseHandler: BaseHandler{
			name:    "get_project",
			method:  "GET",
			pattern: "/v0.1/user/{user_id}/project/{project_id}"},
		Service: service,
	}
	h.handlerFunc = h.GetProject
	Register(h)
	return h
}

/**
GET /v1.0/user/{user_id}/project/{project_id}
Response:
	project query, project define in https://github.com/qetuantuan/jengo_recap/blob/main/model/project.go
*/
func (h *GetProjectHandler) GetProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectId := vars["project_id"]
	project, err := h.Service.GetProject(projectId)
	if err != nil {
		glog.Warningf("get project:%s failed error:[%v]", projectId, err)
		if err == service.NotFoundError {
			util.ReturnError(http.StatusNotFound, err.Error(), w)
			return
		}
		util.ReturnError(http.StatusInternalServerError, err.Error(), w)
		return
	}
	glog.Infof("get project:%s success", projectId)
	util.ReturnSuccessWithObj(http.StatusOK, project, w)
	return
}

type GetProjectsHandler struct {
	BaseHandler
	Service service.ProjectServiceInterface
}

func NewGetProjectsHandler(service service.ProjectServiceInterface) *GetProjectsHandler {
	h := &GetProjectsHandler{
		BaseHandler: BaseHandler{
			name:    "get_projects",
			method:  "GET",
			pattern: "/v0.1/user/{user_id}/projects"},
		Service: service,
	}
	h.handlerFunc = h.GetProjects
	Register(h)
	return h
}

func parseParams2GetProjects(r *http.Request) (userId, groupBy string, offset, maxCount int, filterMap map[string]interface{}, err error) {
	vars := mux.Vars(r)
	userId = vars["user_id"]
	if userId == "" {
		msg := "userId can not be null"
		err = errors.New(msg)
		glog.Warning(msg)
		return
	}
	err = r.ParseForm()
	if err != nil {
		msg := fmt.Sprintf("get runs parse form failed! error: %v", err)
		glog.Warning(msg)
		return
	}
	maxCount, err = getIntParam(r, "max_count")
	if err != nil || maxCount < 0 {
		msg := "max_count must be integer and >= 0"
		err = errors.New(msg)
		glog.Warning(msg)
		return
	}
	offset, err = getIntParam(r, "offset")
	if err != nil || offset < 0 {
		msg := "offset must be integer and >= 0"
		err = errors.New(msg)
		glog.Warning(msg)
		return
	}
	groupBy = r.Form.Get("group_by")
	if groupBy != "" && groupBy != "scm" {

		msg := "not support group_by:" + groupBy
		err = errors.New(msg)
		glog.Warning(msg)
		return
	}
	filterMap = map[string]interface{}{"users": userId}
	scm := r.Form.Get("scms")
	if scm != "" {
		scmArr := strings.Split(scm, ",")
		filterMap["meta.scm"] = scmArr
	}

	enable := r.Form.Get("enable")
	if enable != "" {
		var enableB = true
		enableB, err = strconv.ParseBool(enable)
		if err != nil {
			msg := "enable must be true or false:" + enable
			glog.Warning(msg)
			return
		}
		filterMap["enable"] = enableB
	}
	return
}

/**
GET /v1.0/user/{user_id}/projects
Response:
	project list, project define in https://github.com/qetuantuan/jengo_recap/blob/main/model/project.go
*/

func (h *GetProjectsHandler) GetProjects(w http.ResponseWriter, r *http.Request) {

	userId, groupBy, offset, maxCount, filterMap, err := parseParams2GetProjects(r)

	if err != nil {
		msg := fmt.Sprintf("get projects params failed! error: %v", err)
		glog.Warning(msg)
		util.ReturnError(
			http.StatusBadRequest,
			msg,
			w)
		return
	}

	glog.Info(fmt.Sprintf("begin to get projects for user:%s, offset:%d, max_count:%d", userId, offset, maxCount))
	projects, err := h.Service.GetProjectsByFilter(filterMap, maxCount, offset)
	if err != nil {
		glog.Warningf("get projects failed! error: %v", err)
		util.ReturnError(
			http.StatusInternalServerError,
			err.Error(),
			w)
		return
	}
	// groupby only support scm, generic groupby using reflect will impact performance

	var retObj interface{} = projects
	if groupBy != "" {
		projectMap := groupProjectByScm(projects)
		retObj = projectMap
	}

	//return retObj,err
	util.ReturnSuccessWithObj(
		http.StatusOK,
		retObj,
		w)

	glog.Info("get projects successed! user_id:", userId)
}
func groupProjectByScm(oriProjects []api.Project) (projectsMap map[string]*[]api.Project) {
	projectsMap = make(map[string]*[]api.Project)

	for _, project := range oriProjects {
		if projectArray, ok := projectsMap[project.Meta.Scm]; !ok {
			projectsMap[project.Meta.Scm] = &[]api.Project{project}
			continue
		} else {
			*projectArray = append(*projectArray, project)
		}
	}
	return
}

// example group use reflect, unused
func groupProjects(oriProjects []api.Project, key string) (projectsMap map[string]*[]api.Project) {
	projectsMap = make(map[string]*[]api.Project)
	for _, project := range oriProjects {
		value := reflect.Indirect(reflect.ValueOf(project.Meta)).FieldByName(key).String()
		if projectArray, ok := projectsMap[value]; !ok {
			projectsMap[value] = &[]api.Project{project}
			continue
		} else {
			*projectArray = append(*projectArray, project)
		}
	}
	return
}

type SwitchProjectHandler struct {
	BaseHandler
	Service service.ProjectServiceInterface
}

func NewSwitchProjectHandler(service service.ProjectServiceInterface) *SwitchProjectHandler {
	h := &SwitchProjectHandler{
		BaseHandler: BaseHandler{
			name:    "switch_project",
			method:  "GET",
			pattern: "/v0.1/user/{user_id}/project/{project_id}/action",
			query:   "method"},
		Service: service,
	}
	h.handlerFunc = h.SwitchProject
	Register(h)
	return h
}

/**
GET /v1.0/user/{user_id}/project/{project_id}/action?method=enable
*/

func (h *SwitchProjectHandler) SwitchProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	method := vars["method"]
	projectId := vars["project_id"]
	userId := vars["user_id"]
	glog.Info(fmt.Sprintf("begin to swith project %s to %s", projectId, method))

	enableStatus := true
	if userId == "" || projectId == "" || (method != "enable" && method != "disable") {
		msg := "Check !!!!!uerId / projectId is empty, or method must be enable or disable"
		glog.Warning(msg)
		util.ReturnError(
			http.StatusBadRequest,
			msg,
			w)
		return
	} else {
		if method == "disable" {
			enableStatus = false
		}
	}
	err := h.Service.SwitchProject(userId, projectId, enableStatus)
	if enableStatus {
		if err != nil {
			glog.Warningf("enable hook or switch project status failed! error:%v", err)
			util.ReturnError(
				http.StatusInternalServerError,
				err.Error(),
				w)
			return
		}
		util.ReturnSuccessWithMsg(http.StatusOK, "enabled", w)
		return
	} else { //do disable the hook
		if err != nil {
			glog.Warningf("disable hook or switch project status failed! error:%v", err)
			util.ReturnError(
				http.StatusInternalServerError,
				err.Error(),
				w)
			return
		}
		util.ReturnSuccessWithMsg(http.StatusOK, "disabled", w)
	}
	return

}
