package client

import (
	"encoding/json"
	"fmt"

	"github.com/golang/glog"

	"github.com/qetuantuan/jengo_recap/model"
	"github.com/qetuantuan/jengo_recap/util"
)

const (
	CreateBuildUriPathPattern = "/v0.1/build"
	ListBuildsUriPathPattern  = "/v0.1/builds?repo_id=%s&build_id=%s&user_id=%s"
)

type HttpEngineService struct {
	HostPort string
}

func NewHttpEngineService(server string, port int) *HttpEngineService {
	return &HttpEngineService{
		HostPort: util.GetHostPort4Client(server, port),
	}
}

func (p *HttpEngineService) CreateBuild(params model.EngineCreateBuildParams) (string, error) {
	url := "http://" + p.HostPort + fmt.Sprintf(CreateBuildUriPathPattern)
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return "", err
	}

	var r model.CreateBuildResponse
	if resp, data, err := util.PostHttpRespWithData(url, nil, paramsBytes); err != nil {
		glog.Errorf("Http request failed: %v", err)
		return "", err
	} else if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		glog.Errorf("Http request is not 2xx: %v", *resp)
		return "", ErrStatusCodeNotSuccess
	} else if err := json.Unmarshal(data, &r); err != nil {
		glog.Errorf("Unmarshal response data failed: %v", err)
		return "", err
	} else {
		return r.BuildId, nil
	}
}

func (p *HttpEngineService) ListBuilds(param model.EngineListBuildsParams) (model.Builds, error) {
	// TODO
	var builds model.Builds
	/*
		url := "http://" + p.HostPort + fmt.Sprintf(ListBuildsUriPathPattern, param.RepoId, param.BuildId, param.UserId)
		return util.GetHttpResp(url, nil)
	*/
	return builds, nil
}

func GetBuild(buildId string) (model.Build, error) {
	// TODO
	return model.Build{}, nil
}
