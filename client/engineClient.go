package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/qetuantuan/jengo_recap/api"
	"github.com/qetuantuan/jengo_recap/util"
)

const (
	createEngineRun = "/v0.1/project/%s/run"
	describeRuns    = "/v0.1/runs?project_id=%s&run_id=%s&user_id=%s"
	getRunLog       = "/v0.1/run/%s/run_log?offset=%d&limit=%d"
)

type EngineClient struct {
	HostPort string
}

func NewEngineClient(server string, port int) *EngineClient {
	return &EngineClient{
		HostPort: util.GetHostPort4Client(server, port),
	}
}

func (p *EngineClient) CreateRun(run *api.EngineCreateRunParams) (*http.Response, []byte, error) {
	url := "http://" + p.HostPort + fmt.Sprintf(createEngineRun, run.ProjectId)
	runBytes, err := json.Marshal(run)
	if err != nil {
		return nil, nil, err
	}
	return util.PutHttpRespWithData(url, nil, runBytes)
}

func (p *EngineClient) DescribeRuns(param *api.EngineDescribeRunsParams) (*http.Response, []byte, error) {
	url := "http://" + p.HostPort + fmt.Sprintf(describeRuns, param.ProjectId, param.RunId, param.UserId)
	return util.GetHttpResp(url, nil)
}

func (p *EngineClient) GetRunLog(param *api.GetRunLogParams) (*http.Response, []byte, error) {
	url := "http://" + p.HostPort + fmt.Sprintf(getRunLog, param.RunId, param.Offset, param.Limit)
	return util.GetHttpResp(url, nil)
}
