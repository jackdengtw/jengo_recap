package scm

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/golang/glog"
)

type baseScm struct {
	ApiLink string
}

func (bs *baseScm) httpRequest(method string, uri string, body string, header map[string]string) (response_str []byte, status_code int, err error) {
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	request, err := http.NewRequest(method, uri, strings.NewReader(body))

	glog.Info("uri is", uri)
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
		err = errors.New(fmt.Sprintf("responst status code: %s is not 2**", response.StatusCode))
		return
	}
	defer response.Body.Close()
	response_str, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	return
}
