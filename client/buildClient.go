package client

import (
	"encoding/json"
	"fmt"

	"github.com/golang/glog"

	"github.com/qetuantuan/jengo_recap/model"
	"github.com/qetuantuan/jengo_recap/util"
)

const (
	UpdateBuildUriPathPattern = "/v0.1/build/%s?sub_id=%s"
	InsertBuildUriPathPattern = "/v0.1/build"
)

type HttpBuildService struct {
	HostPort string
}

func NewHttpBuildService(host string, port int) (uc *HttpBuildService) {
	uc = &HttpBuildService{
		HostPort: util.GetHostPort4Client(host, port),
	}

	return
}

func (p *HttpBuildService) InsertBuild(build model.Build) (model.SemanticBuild, error) {
	var sbuild model.SemanticBuild
	url := "http://" + p.HostPort + fmt.Sprintf(InsertBuildUriPathPattern)
	buildBytes, err := json.Marshal(build)
	if err != nil {
		glog.Errorf("Marshal build request failed: %v", err)
		return sbuild, err
	}
	if resp, data, err := util.PostHttpRespWithData(url, nil, buildBytes); err != nil {
		glog.Errorf("Http request failed: %v", err)
		return sbuild, err
	} else if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		glog.Errorf("Http request is not 2xx: %v", *resp)
		return sbuild, ErrStatusCodeNotSuccess
	} else if err := json.Unmarshal(data, &sbuild); err != nil {
		glog.Errorf("Unmarshal response data failed: %v", err)
		return sbuild, err
	} else {
		return sbuild, nil
	}
}

func (p *HttpBuildService) PartialUpdateBuild(sbuildId, buildId string, updateKv map[string]interface{}) error {
	url := "http://" + p.HostPort + fmt.Sprintf(UpdateBuildUriPathPattern, sbuildId, buildId)
	buildBytes, err := json.Marshal(updateKv)
	if err != nil {
		return err
	}
	if resp, _, err := util.PatchHttpRespWithData(url, nil, buildBytes); err != nil {
		glog.Errorf("Http request failed: %v", err)
		return err
	} else if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		glog.Errorf("Http request is not 2xx: %v", *resp)
		return ErrStatusCodeNotSuccess
	} else {
		return nil
	}
}

func (r *HttpBuildService) GetSemanticBuildsByFilter(
	filter map[string]string, maxCount, offset int) (builds model.Builds, err error) {
	// TODO
	return
}

func (r *HttpBuildService) GetSemanticBuildsByIds(
	buildIds []string) (builds model.Builds, err error) {
	// TODO
	return
}
