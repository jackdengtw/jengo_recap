package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/qetuantuan/jengo_recap/api"
	"github.com/qetuantuan/jengo_recap/service"
	"github.com/qetuantuan/jengo_recap/util"
)

/**
-X POST /v0.1/project/project_123/run -d '{"project_id":"project_1234", "user_id":"user", "scm_url":"scm_url", "branch":"master", "commit_id":"commit_id_1234", "commit_user":"", "commit_link": "", "diff_link": ""}'
parameter: EngineCreateRunParams
	user_id: string, necessary
	scm_url: string, necessary
	branch: string, necessary
	commits: []PushEventCommit, necessary
	repository: *PushEventRepository, necessary
*/

type CreateRunHandler struct {
	BaseHandler
	Service service.EngineRunServiceInterface
}

func NewCreateRunHandler(service service.EngineRunServiceInterface) *CreateRunHandler {
	h := &CreateRunHandler{
		BaseHandler: BaseHandler{
			name:    "create_run",
			method:  "POST",
			pattern: "/v0.1/project/{project_id}/run"},
		Service: service,
	}
	h.handlerFunc = h.CreateRun
	Register(h)
	return h
}

func (h *CreateRunHandler) CreateRun(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		glog.Errorf("io read params error: %v", err)
		util.ReturnError(http.StatusBadRequest, "CreateRun: read params error", w)
		return
	}
	params := &api.EngineCreateRunParams{}
	if err := json.Unmarshal(data, params); err != nil {
		glog.Errorf("json.Unmarshal get params error: %v", err)
		util.ReturnError(http.StatusBadRequest, "CreateRun: json.Unmarshal get params error", w)
		return
	}
	params.ProjectId = vars["project_id"]

	runObj, err := h.Service.CreateRun(params)
	if err != nil {
		glog.Errorf("Service.CreateRun: %v", err)
		util.ReturnError(http.StatusInternalServerError, "CreateRun: Service.CreateRun error", w)
		return
	}

	body, err := json.Marshal(runObj)
	if err != nil {
		glog.Errorf("json.Marshal error: %v", err)
		util.ReturnError(http.StatusInternalServerError, "CreateRun: json.Marshal response error", w)
		return
	}
	glog.Infof("Service.CreateRun successfully! params: %v", params)
	util.ReturnSuccessWithBytes(http.StatusOK, body, "", map[string]string{}, w)
	return
}

/**
-X GET /v0.1/runs?project_id=project_1234567&run_id=run_123&user_id=user_123"
parameters; EngineDescribeRunsParams
	project_id: string, not necessary, for filter
	run_id: string, not necessary, for filter
	user_id: string, not necessary, for filter
*/

type DescribeRunsHandler struct {
	BaseHandler
	Service service.EngineRunServiceInterface
}

func NewDescribeRunsHandler(service service.EngineRunServiceInterface) *DescribeRunsHandler {
	h := &DescribeRunsHandler{
		BaseHandler: BaseHandler{
			name:    "describe_runs",
			method:  "GET",
			pattern: "/v0.1/runs"},
		Service: service,
	}
	h.handlerFunc = h.DescribeRuns
	Register(h)
	return h
}

func (h *DescribeRunsHandler) DescribeRuns(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	offset, _ := strconv.Atoi(vars["offset"])
	limit, _ := strconv.Atoi(vars["limit"])

	r.ParseForm()
	params := &api.EngineDescribeRunsParams{
		ProjectId: r.Form.Get("project_id"),
		RunId:     r.Form.Get("run_id"),
		UserId:    r.Form.Get("user_id"),
		EventId:   r.Form.Get("event_id"),
		Offset:    offset,
		Limit:     limit,
	}

	runs, err := h.Service.DescribeRuns(params)
	if err != nil {
		glog.Infof("DescribeRuns: %v", err)
		util.ReturnError(http.StatusInternalServerError, "DescribeRuns: service.DescribeRuns error", w)
		return
	}

	body, err := json.Marshal(runs)
	if err != nil {
		glog.Errorf("Service.DescribeRuns json: %v", err)
		util.ReturnError(http.StatusInternalServerError, "DescribeRuns: json.Marshal error", w)
		return
	}
	glog.Infof("Service.DescribeRuns successfully! params: %v", params)
	if len(runs) == 0 {
		util.ReturnSuccess(http.StatusOK, "[]", "", w)
	} else {
		util.ReturnSuccessWithBytes(http.StatusOK, body, "", map[string]string{}, w)
	}
	return
}

/**
-X GET /v0.1/run/run_id"
*/

type DescribeRunHandler struct {
	BaseHandler
	Service service.EngineRunServiceInterface
}

func NewDescribeRunHandler(service service.EngineRunServiceInterface) *DescribeRunHandler {
	h := &DescribeRunHandler{
		BaseHandler: BaseHandler{
			name:    "describe_run",
			method:  "GET",
			pattern: "/v0.1/run/{run_id}"},
		Service: service,
	}
	h.handlerFunc = h.DescribeRun
	Register(h)
	return h
}

func (h *DescribeRunHandler) DescribeRun(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	runId := vars["run_id"]

	run, err := h.Service.DescribeRun(runId)
	if err != nil {
		glog.Infof("DescribeRun: %v", err)
		util.ReturnError(http.StatusInternalServerError, "DescribeRun: service.DescribeRun error", w)
		return
	}

	body, err := json.Marshal(run)
	if err != nil {
		glog.Errorf("Service.DescribeRun json: %v", err)
		util.ReturnError(http.StatusInternalServerError, "DescribeRun: json.Marshal error", w)
		return
	}
	glog.Infof("Service.DescribeRun successfully! runId: %s", runId)
	if run == nil {
		util.ReturnError(http.StatusNotFound, "{}", w)
	} else {
		util.ReturnSuccessWithBytes(http.StatusOK, body, "", map[string]string{}, w)
	}
	return
}
