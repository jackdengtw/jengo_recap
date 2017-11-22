package util

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/golang/glog"
)

func PostHttpRespWithData(req string, headers map[string]string, data []byte) (*http.Response, []byte, error) {
	return httpRespWithData("POST", req, headers, data)
}

func PutHttpRespWithData(req string, headers map[string]string, data []byte) (*http.Response, []byte, error) {
	return httpRespWithData("PUT", req, headers, data)
}

func GetHttpResp(req string, headers map[string]string) (*http.Response, []byte, error) {
	return httpResp("GET", req, headers)
}

func DeleteHttp(req string, headers map[string]string) (*http.Response, []byte, error) {
	return httpResp("DELETE", req, headers)
}

func httpRespWithData(method string, req string, headers map[string]string, data []byte) (*http.Response, []byte, error) {
	// data := []byte(fmt.Sprintf(`{"reboot" : {"type" : "SOFT"}}`))
	client := &http.Client{Timeout: 10 * time.Second}
	glog.Info("[" + method + " HttpRespWithHeaders] url: " + req + ", post_data: " + string(data))
	if req, err := http.NewRequest(method, req, bytes.NewReader(data)); err != nil {
		glog.Errorf("["+method+"HttpRespWithHeaders] http.NewRequest err %v", err)
		return nil, nil, err
	} else {
		for k, v := range headers {
			req.Header.Add(k, string(v))
		}

		if resp, err := client.Do(req); err != nil {
			glog.Errorf("["+method+"HttpRespWithHeaders] client.Do err %v", err)
			return nil, nil, err
		} else {
			defer resp.Body.Close()
			if raw, err := ioutil.ReadAll(resp.Body); err != nil {
				glog.Errorf("["+method+"HttpRespWithHeaders] ioutils.ReadAll err %v", err)
				return resp, nil, err
			} else {
				glog.Info("INFO][" + method + "HttpRespWithHeaders] raw: " + string(raw))
				return resp, raw, nil
			}
		}
	}
}

func httpResp(method string, req string, headers map[string]string) (*http.Response, []byte, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	glog.Infof("[%s HttpRespWithHeaders] url: %s", method, req)
	if req, err := http.NewRequest(method, req, nil); err != nil {
		glog.Errorf("[%s HttpRespWithHeaders] http.NewRequest err %v", method, err)
		return nil, nil, err
	} else {
		for k, v := range headers {
			req.Header.Add(k, string(v))
		}

		if resp, err := client.Do(req); err != nil {
			glog.Errorf("[%s HttpRespWithHeaders] client.Do err %v", method, err)
			return nil, nil, err
		} else {
			defer resp.Body.Close()
			if raw, err := ioutil.ReadAll(resp.Body); err != nil {
				return resp, nil, err
			} else {
				//glog.Infof("[GetHttpRespWithHeaders] raw: " + string(raw))
				return resp, raw, nil
			}
		}
	}
}
