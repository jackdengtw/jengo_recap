package util

import (
	"encoding/json"
	"net/http"

	"github.com/golang/glog"
	"github.com/qetuantuan/jengo_recap/constant"
)

type logMethods struct {
	Errorf func(string, ...interface{})
}

var logger = logMethods{
	Errorf: glog.Errorf,
}

type CommonResponse struct {
	Id  string `json:"id,omitempty"`
	Msg string `json:"msg,omitempty"`
	Url string `json:"url,omitempty"`
}

func ReturnError(httpCode int, msg string, w http.ResponseWriter) {
	ReturnErrorWithObj(httpCode, CommonResponse{Msg: msg}, w)
}

func ReturnErrorWithObj(httpCode int, obj interface{}, w http.ResponseWriter) {
	// logger.Errorf("%v [%v]", logMsg, appErr)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var httpBody []byte
	var err error
	if httpBody, err = json.Marshal(obj); err != nil {
		glog.Errorf("Marshal Error: %v", err)
		ReturnErrorWithObj(http.StatusInternalServerError, CommonResponse{Msg: "Marshal Error"}, w)
		return
	}

	w.WriteHeader(httpCode)
	_, err = w.Write([]byte(httpBody))
	if err != nil {
		glog.Errorf("Error sending http response: %v", err)
	}
	return
}

func ReturnSuccess(httpCode int, httpBody string, contentType string, w http.ResponseWriter) {
	ReturnSuccessWithBytes(httpCode, []byte(httpBody), contentType, map[string]string{}, w)
}

func ReturnSuccessWithBytes(httpCode int, httpBody []byte, contentType string, headerMap map[string]string, w http.ResponseWriter) {
	if httpCode != http.StatusNoContent && contentType == "" {
		w.Header().Set(constant.CONTENT_TYPE, "application/json; charset=UTF-8")
		// Assuming utf8 encoding. if !utf8.Valid(httpBody) {
	} else if contentType != "" {
		w.Header().Set(constant.CONTENT_TYPE, contentType)
	}

	for k, v := range headerMap {
		w.Header().Set(k, v)
	}

	w.WriteHeader(httpCode)
	if httpCode != http.StatusNoContent {
		_, err := w.Write(httpBody)
		if err != nil {
			glog.Errorf("Error sending http response: %v", err)
		}
	}
}

// ReturnSuccessWithMap
// TODO: deprecated. remove this whenever no one using it.
func ReturnSuccessWithMap(httpCode int, amap map[string]string, w http.ResponseWriter) {
	ReturnSuccessWithMaps(httpCode, amap, map[string]string{}, w)
}

func ReturnSuccessWithMsg(httpCode int, msg string, w http.ResponseWriter) {
	ReturnSuccessWithObj(httpCode, CommonResponse{Msg: msg}, w)
}

func ReturnSuccessWithObj(httpCode int, obj interface{}, w http.ResponseWriter) {
	ReturnSuccessWithObjAndHeaders(httpCode, obj, map[string]string{}, w)
}

func ReturnSuccessWithObjAndHeaders(httpCode int, obj interface{}, headerMap map[string]string, w http.ResponseWriter) {
	w.Header().Set(constant.CONTENT_TYPE, "application/json; charset=UTF-8")

	for k, v := range headerMap {
		w.Header().Set(k, v)
	}

	var httpBody []byte
	var err error
	if httpBody, err = json.Marshal(obj); err != nil {
		glog.Errorf("Marshal Error: %v", err)
		ReturnErrorWithObj(http.StatusInternalServerError, CommonResponse{Msg: "Marshal Error"}, w)
		return
	} else {
		w.WriteHeader(httpCode)
		_, err = w.Write(httpBody)
		if err != nil {
			glog.Errorf("Error sending http response: %v", err)
		}
		return
	}
}

// ReturnSuccessWithMaps
// TODO: remove this whenever no one using it.
// Encourage to use ReturnSuccessWithObjAndHeaders
func ReturnSuccessWithMaps(httpCode int, amap map[string]string, headerMap map[string]string, w http.ResponseWriter) {
	for k, v := range headerMap {
		w.Header().Set(k, v)
	}

	var httpBody []byte
	var err error
	if httpBody, err = json.Marshal(amap); err != nil {
		glog.Errorf("Marshal Error: %v", err)
		ReturnErrorWithObj(http.StatusInternalServerError, CommonResponse{Msg: "Marshal Error"}, w)
		return
	} else {
		w.WriteHeader(httpCode)
		_, err = w.Write(httpBody)
		if err != nil {
			glog.Errorf("Error sending http response: %v", err)
		}
		return
	}
}
