package handler

import (
	"encoding/json"
	"net/http"

	"github.com/qetuantuan/jengo_recap/api"
	"github.com/qetuantuan/jengo_recap/client"

	"github.com/golang/glog"
)

type Trigger interface {
	TriggerRun(runReq *api.EngineCreateRunParams) (
		httpResp *http.Response, runResp *api.Run, err error)
}

type GitHubTrigger struct {
	EngineClient *client.EngineClient
}

func NewGitHubTrigger(server string, port int) *GitHubTrigger {
	return &GitHubTrigger{
		EngineClient: client.NewEngineClient(server, port),
	}
}
func (t *GitHubTrigger) TriggerRun(runReq *api.EngineCreateRunParams) (
	httpResp *http.Response, runResp *api.Run, err error) {
	glog.Infof("starting\n")

	var raw []byte
	httpResp, raw, err = t.EngineClient.CreateRun(runReq)

	if err != nil {
		glog.Infof("CreateRun error: err(%s)",
			err.Error())
		return httpResp, nil, err
	} else if err = json.Unmarshal(raw, runResp); err != nil {
		glog.Infof("unmarshal error: err(%s)",
			err.Error())
		return httpResp, nil, err
	} else {
		glog.Infof("succeed: httpResp(%v), "+
			"raw(%v), runResp(%v)",
			httpResp, raw, runResp)
		return
	}
	glog.Infof("ending\n")
	return
}
