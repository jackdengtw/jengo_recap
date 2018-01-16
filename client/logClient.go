package client

import (
	"encoding/json"
	"fmt"

	"github.com/golang/glog"

	"github.com/qetuantuan/jengo_recap/model"
	"github.com/qetuantuan/jengo_recap/util"
)

const (
	BuildLogUriPathPattern = "%s?offset=%d&limit=%d"
)

type HttpBuildLogService struct {
	HostPort string
}

func NewHttpBuildLogService(server string, port int) *HttpBuildLogService {
	return &HttpBuildLogService{
		HostPort: util.GetHostPort4Client(server, port),
	}
}

func (p *HttpBuildLogService) GetLog(param *model.GetBuildLogParams) (model.BuildLog, error) {
	url := "http://" + p.HostPort + fmt.Sprintf(BuildLogUriPathPattern, param.LogUri, param.Offset, param.Limit)
	var log model.BuildLog
	if resp, data, err := util.GetHttpResp(url, nil); err != nil {
		glog.Errorf("Http request failed: %v", err)
		return log, err
	} else if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		glog.Errorf("Http request is not 2xx: %v", *resp)
		return log, ErrStatusCodeNotSuccess
	} else if err := json.Unmarshal(data, &log); err != nil {
		glog.Errorf("Unmarshal response data failed: %v", err)
		return log, err
	} else {
		return log, nil
	}
}
