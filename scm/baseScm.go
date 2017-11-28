package scm

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/golang/glog"
	"github.com/qetuantuan/jengo_recap/model"
)

type Scm interface {
	Init(string)
	SetHook(string) error
	EditHook(string) error
	DeleteHook(string) error
	GetProjectList() ([]model.Project, error)
}

type baseScm struct {
	ApiLink string
	User    string
	Token   string
	HookURI string
}

func (bs *baseScm) httpRequest(method string, uri string, data string, header map[string]string) (responseStrs []byte, statusCode int, err error) {
	glog.Infof("[httpRequest] uri:[%s];data:[%s];header:[%v]", uri, data, header)
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	request, err := http.NewRequest(method, uri, strings.NewReader(data))
	fmt.Println(uri)
	if err != nil {
		return
	}
	for key, value := range header {
		request.Header.Set(key, value)
	}
	response, err := client.Do(request)
	if err != nil {
		return
	}

	if response.StatusCode > 300 || response.StatusCode < 200 {
		resStr := ""
		defer response.Body.Close()
		responseStrs, err = ioutil.ReadAll(response.Body)
		if err == nil {
			resStr = string(responseStrs)
		}
		err = errors.New(fmt.Sprintf("responst status code: %s is not 2**, response:%s", response.StatusCode, resStr))
		return
	}
	defer response.Body.Close()
	responseStrs, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	return
}
