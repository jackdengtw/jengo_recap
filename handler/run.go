package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"errors"
	"strings"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/qetuantuan/jengo_recap/api"
	"github.com/qetuantuan/jengo_recap/service"
	"github.com/qetuantuan/jengo_recap/util"
)

/**
GET /v1.0/user/{user_id}/project/{project_id}/runs
Function:
	get runs of project:project_id
Resposne:
	run list: run define in https://github.com/qetuantuan/jengo_recap/blob/main/model/run.go
*/

func getIntParam(r *http.Request, name string) (retValue int, err error) {

	retStr := r.Form.Get(name)
	retValue = 0
	if retStr != "" {
		fV, errT := strconv.Atoi(retStr)
		if errT != nil {
			err = errT
			return
		} else {
			retValue = fV
		}
	}
	return
}

type GetBuildsByFilterHandler struct {
	BaseHandler
	Service service.RunServiceInterface
}

func NewGetBuildsByFilterHandler(service service.RunServiceInterface) *GetBuildsByFilterHandler {
	h := &GetBuildsByFilterHandler{
		BaseHandler: BaseHandler{
			name:    "get_build_by_project",
			method:  "GET",
			pattern: "/v0.1/user/{user_id}/project/{project_id}/builds"},
		Service: service,
	}
	h.handlerFunc = h.GetBuildsByFilter
	Register(h)
	return h
}
func parseParams2GetBuildsByFilter(r *http.Request) (projectId string, offset, maxCount int, filter map[string]string, err error) {
	vars := mux.Vars(r)
	projectId = vars["project_id"]

	err = r.ParseForm()
	if err != nil {
		msg := fmt.Sprintf("get builds parse form failed! error:%v", err)
		glog.Warning(msg)
		err = errors.New(msg)
		return
	}
	maxCount, err = getIntParam(r, "max_count")
	if err != nil || maxCount < 0 {
		msg := "max_count must be integer and >= 0"
		glog.Warning(msg)
		err = errors.New(msg)
		return
	}
	offset, err = getIntParam(r, "offset")
	if err != nil || offset < 0 {
		msg := "offset must be integer and >=0 "
		glog.Warning(msg)
		err = errors.New(msg)
		return
	}

	glog.Infof("begin to get builds of project:%s, maxCount:%d, offset:%d", projectId, maxCount, offset)
	filter = map[string]string{}
	filter["projectid"] = projectId
	branch := r.Form.Get("branch")
	if branch != "" {
		filter["branch"] = branch
	}
	return
}
func (h *GetBuildsByFilterHandler) GetBuildsByFilter(w http.ResponseWriter, r *http.Request) {
	projectId, offset, maxCount, filter, err := parseParams2GetBuildsByFilter(r)
	if err != nil {
		msg := fmt.Sprintf("Params Error,BAD request! error:%v", err)
		glog.Warningf(msg)
		util.ReturnError(http.StatusBadRequest, msg, w)
	}
	builds, err := h.Service.GetBuildsByFilter(filter, maxCount, offset)
	if err != nil {
		msg := fmt.Sprintf("get builds failed! error:%v", err)
		glog.Warningf(msg)
		util.ReturnError(http.StatusInternalServerError, msg, w)
		return
	}

	util.ReturnSuccessWithObj(http.StatusOK, builds, w)

	glog.Info("get builds success! project_id:", projectId)
}

type InsertRunHandler struct {
	BaseHandler
	Service service.RunServiceInterface
}

func NewInsertRunHandler(service service.RunServiceInterface) *InsertRunHandler {
	h := &InsertRunHandler{
		BaseHandler: BaseHandler{
			name:    "insert_run",
			method:  "PUT",
			pattern: "/v0.1/user/{user_id}/project/{project_id}/run/{run_id}"},
		Service: service,
	}
	h.handlerFunc = h.InsertRun
	Register(h)
	return h
}
func parseParams2InsertRun(r *http.Request) (run api.Run, err error) {
	vars := mux.Vars(r)

	runId := vars["run_id"]
	if runId == "" {
		msg := "runId can't be empty"
		glog.Warning(msg)
		err = errors.New(msg)
		return
	}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := fmt.Sprintf("read data from request failed! error:%v", err)
		glog.Warning(msg)
		err = errors.New(msg)
		return
	}

	if err = json.Unmarshal(data, &run); err != nil {
		msg := fmt.Sprintf("Unmarshal data of run failed! error:%v", err)
		glog.Warning(msg)
		err = errors.New(msg)
		return
	}
	return
}
func (h *InsertRunHandler) InsertRun(w http.ResponseWriter, r *http.Request) {

	run, err := parseParams2InsertRun(r)
	if err != nil {
		msg := fmt.Sprintf("Params Error,BAD request!error: %v", err)
		glog.Warningf(msg)
		util.ReturnError(http.StatusBadRequest, msg, w)
		return
	}
	glog.Infof("Got run ready to save: %+v", run)

	build, err := h.Service.InsertRun(run)
	if err != nil {
		msg := fmt.Sprintf("Insert run failed! error: v%", err)
		glog.Warningf(msg)
		util.ReturnError(http.StatusInternalServerError, msg, w)
		return
	}
	glog.Infof("Insert run (id:%s) to build (id:%s) success!", run.Id, build.Id)
	util.ReturnSuccessWithMap(http.StatusOK, map[string]string{"build_id": build.Id}, w)

	return
}

type UpdateRunHandler struct {
	BaseHandler
	Service service.RunServiceInterface
}

func NewUpdateRunHandler(service service.RunServiceInterface) *UpdateRunHandler {
	h := &UpdateRunHandler{
		BaseHandler: BaseHandler{
			name:    "update_run",
			method:  "PUT",
			pattern: "/v0.1/user/{user_id}/project/{project_id}/build/{build_id}/run/{run_id}"},
		Service: service,
	}
	h.handlerFunc = h.UpdateRun
	Register(h)
	return h
}

/**
PUT /v1.0/user/{user_id}/project/{project_id}/run/{run_id}   -d "run_info"
Function:
	upsert run, if run_id already exist, update, otherwise insert new one
*/
func parseParams2UpRun(r *http.Request) (buildId, runId string, run api.Run, err error) {
	vars := mux.Vars(r)
	buildId = vars["build_id"]
	runId = vars["run_id"]
	if buildId == "" || runId == "" {
		msg := "buildId or runId can't be empty"
		glog.Warning(msg)
		err = errors.New(msg)
		return
	}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := fmt.Sprintf("read data from request failed! error:%v", err)
		glog.Warning(msg)
		err = errors.New(msg)
		return
	}

	if err = json.Unmarshal(data, &run); err != nil {
		msg := fmt.Sprintf("Unmarshal data of run failed! error:%v", err)
		glog.Warning(msg)
		err = errors.New(msg)
		return
	}
	return
}
func (h *UpdateRunHandler) UpdateRun(w http.ResponseWriter, r *http.Request) {
	buildId, runId, run, err := parseParams2UpRun(r)
	if err != nil {
		msg := fmt.Sprintf("Params Error,BAD request! error: %v", err)
		glog.Warningf(msg)
		util.ReturnError(http.StatusBadRequest, msg, w)
		return
	}

	err = h.Service.UpdateRun(buildId, run)
	if err != nil {
		msg := fmt.Sprintf("update run failed! error:%v", err)
		glog.Warningf(msg)
		util.ReturnError(http.StatusInternalServerError, msg, w)
		return
	}

	msg := fmt.Sprintf("update run (id:%s) to build (id:%s) success", runId, buildId)
	glog.Info(msg)
	util.ReturnSuccessWithMsg(http.StatusOK, msg, w)
	return
}

type PartialUpdateRunHandler struct {
	BaseHandler
	Service service.RunServiceInterface
}

func NewPartialUpdateRunHandler(service service.RunServiceInterface) *PartialUpdateRunHandler {
	h := &PartialUpdateRunHandler{
		BaseHandler: BaseHandler{
			name:    "partial_update_run",
			method:  "PATCH",
			pattern: "/v0.1/user/{user_id}/project/{project_id}/build/{build_id}/run/{run_id}"},
		Service: service,
	}
	h.handlerFunc = h.PartialUpdateRun
	Register(h)
	return h
}
func parseParams2PartialUpRun(r *http.Request) (buildId, runId string, run api.Run, runInterface map[string]interface{}, err error) {
	vars := mux.Vars(r)
	buildId = vars["build_id"]
	runId = vars["run_id"]
	if buildId == "" || runId == "" {
		msg := "buildId or runId can't be empty"
		glog.Warning(msg)
		err = errors.New(msg)
		return
	}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := fmt.Sprintf("read data from request failed! error:%v", err)
		glog.Warning(msg)
		err = errors.New(msg)
		return
	}

	// 1. validate the patch body
	// 2. get the run for UpdateProjectStateAndLatestBuild
	if err = json.Unmarshal(data, &run); err != nil {
		msg := fmt.Sprintf("Unmarshal data of run failed! error:%v", err)
		glog.Warning(msg)
		err = errors.New(msg)
		return
	}

	if err = json.Unmarshal(data, &runInterface); err != nil {
		msg := fmt.Sprintf("Unmarshal data of run failed! error:%v", err)
		glog.Warning(msg)
		err = errors.New(msg)
		return
	}
	return
}
func (h *PartialUpdateRunHandler) PartialUpdateRun(w http.ResponseWriter, r *http.Request) {
	buildId, runId, run, runInterface, err := parseParams2PartialUpRun(r)
	if err != nil {
		msg := fmt.Sprintf("Params Error,BAD request! error: %v", err)
		glog.Warningf(msg)
		util.ReturnError(http.StatusBadRequest, msg, w)
		return
	}

	err = h.Service.PartialUpdateRun(buildId, run, runInterface)
	if err != nil {
		msg := fmt.Sprintf("partial update run failed! error:%v", err)
		glog.Warningf(msg)
		util.ReturnError(http.StatusInternalServerError, msg, w)
		return
	}
	msg := fmt.Sprintf("update run (id:%s) to build (id:%s) success", runId, buildId)
	glog.Info(msg)
	util.ReturnSuccessWithMsg(http.StatusOK, msg, w)
	return
}

/*
   b : base builds id
   c : compare builds
   m : missing builds
*/
func getMissingBuilds(c []string, b []api.Build) (m []string) {
	bMap := make(map[string]string)
	for _, build := range b {
		bMap[build.Id] = ""
	}
	for _, id := range c {
		if _, ok := bMap[id]; !ok {
			m = append(m, id)
		}
	}
	return
}

type GetBuildsByIdsHandler struct {
	BaseHandler
	Service service.RunServiceInterface
}

func NewGetBuildsByIdsHandler(service service.RunServiceInterface) *GetBuildsByIdsHandler {
	h := &GetBuildsByIdsHandler{
		BaseHandler: BaseHandler{
			name:    "get_run_by_ids",
			method:  "GET",
			pattern: "/v0.1/builds",
			query:   "ids"},
		Service: service,
	}
	h.handlerFunc = h.GetBuildsByIds
	Register(h)
	return h
}
func (h *GetBuildsByIdsHandler) GetBuildsByIds(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	buildsId := vars["ids"]
	if buildsId == "" {
		msg := "ids can't be empty"
		glog.Warning(msg)
		util.ReturnError(http.StatusBadRequest, msg, w)
		return
	}
	idArray := strings.Split(buildsId, ",")
	builds, err := h.Service.GetBuildsByIds(idArray)
	if err != nil {
		msg := fmt.Sprintf("get builds failed! err :[%v]", err)
		glog.Warningf(msg)
		util.ReturnError(http.StatusInternalServerError, msg, w)
		return
	}
	mBuildsId := getMissingBuilds(idArray, builds)
	if len(mBuildsId) != 0 {
		msg := fmt.Sprintf("some build missing!%v", mBuildsId)
		util.ReturnError(http.StatusNotFound, msg, w)
		return
	}
	util.ReturnSuccessWithObj(http.StatusOK, builds, w)
	return

}
