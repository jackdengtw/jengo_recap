package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/qetuantuan/jengo_recap/service"
	"github.com/qetuantuan/jengo_recap/util"
)

type PutLogHandler struct {
	BaseHandler
	Service service.RunLogServiceInterface
}

func NewPutLogHandler(service service.RunLogServiceInterface) *PutLogHandler {
	h := &PutLogHandler{
		BaseHandler: BaseHandler{
			name:    "put_log",
			method:  "PUT",
			pattern: "/v0.1/user/{user}/project/{project_id}/build/{build_id}/run/{run_id}/log"},
		Service: service,
	}
	h.handlerFunc = h.PutLog
	Register(h)
	return h
}

/**
PUT /v1.0/user/{user}/project/{project_id}/build/{build_id}/run/{run_id}/log -d "log data"
Function:
	store log into mongodb and update logid in run
Response:
	{
		"run_id":"run1",
		"log_id":"597a9a58b390cb41933bee96"
	}
*/
func (h *PutLogHandler) PutLog(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	buildId := vars["build_id"]
	runId := vars["run_id"]
	glog.Info(fmt.Sprintf("begin to put log of run:%s", runId))
	logs, err := ioutil.ReadAll(r.Body)
	if buildId == "" || runId == "" || err != nil {
		msg := fmt.Sprintf("buildId :%s or runId: %s is empty; OR read logs from request failed! error:%v", buildId, runId, err)
		glog.Warning(msg)
		util.ReturnError(http.StatusBadRequest, msg, w)
		return
	}
	id, err := h.Service.PutLog(logs, buildId, runId)
	if err != nil {
		msg := fmt.Sprintf("Put logs failed! error:%v", err)
		glog.Warning(msg)
		util.ReturnError(http.StatusInternalServerError, msg, w)
		return
	}
	util.ReturnSuccessWithMap(http.StatusOK, map[string]string{"run_id": runId, "log_id": id}, w)
	glog.Info(fmt.Sprintf("put log of run:%s success, logid:%s", runId, id))
	return
}

type GetLogHandler struct {
	BaseHandler
	Service service.RunLogServiceInterface
}

func NewGetLogHandler(service service.RunLogServiceInterface) *GetLogHandler {
	h := &GetLogHandler{
		BaseHandler: BaseHandler{
			name:    "get_log",
			method:  "GET",
			pattern: "/v0.1/log/{log_id}"},
		Service: service,
	}
	h.handlerFunc = h.GetLog
	Register(h)
	return h
}

/**
GET /v1.0/log/{log_id}
*/
func (h *GetLogHandler) GetLog(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	logId := vars["log_id"]
	glog.Info(fmt.Sprintf("begin to get run log:%s", logId))

	runLog, err := h.Service.GetLog(logId)
	if err != nil {
		if err == service.NotFoundError {
			msg := fmt.Sprintf("log(id:%s)not found", logId)
			glog.Warning(msg)
			util.ReturnError(http.StatusNotFound, msg, w)
			return
		}
		msg := fmt.Sprintf("get log failed! error:%v", err)
		glog.Warning(msg)
		util.ReturnError(http.StatusInternalServerError, msg, w)
		return
	}

	util.ReturnSuccess(http.StatusOK, runLog.Content, "text/plain; charset=UTF-8", w)
	glog.Info(fmt.Sprintf("get run log:%s success", logId))

	return
}
