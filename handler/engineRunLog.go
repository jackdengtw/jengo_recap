package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/qetuantuan/jengo_recap/api"
	"github.com/qetuantuan/jengo_recap/constant"
	"github.com/qetuantuan/jengo_recap/service"
	"github.com/qetuantuan/jengo_recap/util"
)

/**
-X GET /v0.1/run/run-8674665223082153551/run_log?offset=1&limit=12"
-H "Accept-Encoding:gzip"
parameter:
	offset: int, necessary
	limit: int, necessary
*/

type GetRunLogHandler struct {
	BaseHandler
	Service service.RunLogServiceInterface
}

func NewGetRunLogHandler(service service.RunLogServiceInterface) *GetRunLogHandler {
	h := &GetRunLogHandler{
		BaseHandler: BaseHandler{
			name:    "get_runlog",
			method:  "GET",
			pattern: "/v0.1/run/{run_id}/run_log"},
		Service: service,
	}
	h.handlerFunc = h.GetRunLog
	Register(h)
	return h
}

func (h *GetRunLogHandler) GetRunLog(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	r.ParseForm()
	offset, err1 := strconv.Atoi(r.Form.Get("offset"))
	limit, err2 := strconv.Atoi(r.Form.Get("limit"))
	if err1 != nil || err2 != nil {
		glog.Errorf("get offset error: %v, get limit error: %v", err1, err2)
		util.ReturnError(http.StatusBadRequest, "GetRunLog: parameters error", w)
		return
	}
	params := &api.GetRunLogParams{
		Offset: offset,
		Limit:  limit,
		RunId:  vars["run_id"],
	}

	logObj, ce := h.Service.GetRunLog(params)
	if ce != nil {
		glog.Errorf("Service.GetRunLog: %v", ce)
		util.ReturnError(http.StatusInternalServerError, "GetRunLog: Service.GetRunLog error", w)
		return
	}

	accept_encoding := r.Header.Get(constant.ACCEPT_ENCODING)
	if strings.Contains(accept_encoding, "gzip") || strings.Contains(accept_encoding, "*") {
		logGzip, ce := util.GzipEncode([]byte(logObj.Content))
		if ce != nil {
			glog.Errorf("Service.GetRunLog Gzip: %v", ce)
			util.ReturnError(http.StatusInternalServerError, "getRunLog gzip error", w)
			return
		}
		glog.Info("GetLog gzip: successful.")
		util.ReturnSuccessWithBytes(
			http.StatusOK,
			logGzip,
			constant.CONTENT_TYPE_TEXT_PLAIN,
			map[string]string{
				constant.HEADER_JENGO_RUNID:  logObj.RunId,
				constant.HEADER_JENGO_LENGTH: strconv.Itoa(logObj.Length),
				constant.CONTENT_ENCODING:    "gzip",
			},
			w)
	} else if accept_encoding == "" {
		glog.Info("GetLog raw: successful.")
		util.ReturnSuccessWithBytes(
			http.StatusOK,
			[]byte(logObj.Content),
			constant.CONTENT_TYPE_TEXT_PLAIN,
			map[string]string{
				constant.HEADER_JENGO_RUNID:  logObj.RunId,
				constant.HEADER_JENGO_LENGTH: strconv.Itoa(logObj.Length),
			},
			w)
	} else {
		util.ReturnError(http.StatusNotAcceptable, "getRunLog Accept-Encoding error", w)
	}

	return
}
