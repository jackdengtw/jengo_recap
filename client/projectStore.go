package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/qetuantuan/jengo_recap/api"
	"github.com/qetuantuan/jengo_recap/util"
)

const (
	createRun  = "/v0.1/user/%s/project/%s/run/%s"
	updateRun  = "/v0.1/user/%s/project/%s/build/%s/run/%s"
	putLog     = "/v0.1/user/%s/project/%s/build/%s/run/%s/log"
	getProject = "/v0.1/user/%s/project/%s"
)

type ProjectStoreClientInterface interface {
	InsertRun(run *api.Run) (*http.Response, []byte, error)
	UpdatePatchRun(run *api.PatchRun, buildId string) (*http.Response, []byte, error)
	PutLog(userId string, projectId string, buildId string, runId string, runLog string) (*http.Response, []byte, error)
	GetProject(userId string, projectId string) (project api.Project, err error)
}

type ProjectStoreClient struct {
	HostPort string
}

func NewProjectStoreClient(host string, port int) (uc *ProjectStoreClient) {
	uc = &ProjectStoreClient{
		HostPort: util.GetHostPort4Client(host, port),
	}

	return
}

func (p *ProjectStoreClient) InsertRun(run *api.Run) (*http.Response, []byte, error) {
	url := "http://" + p.HostPort + fmt.Sprintf(createRun, run.UserId, run.ProjectId, run.Id)
	runBytes, err := json.Marshal(run)
	if err != nil {
		return nil, nil, err
	}
	return util.PutHttpRespWithData(url, nil, runBytes)
}

func (p *ProjectStoreClient) UpdatePatchRun(run *api.PatchRun, buildId string) (*http.Response, []byte, error) {
	url := "http://" + p.HostPort + fmt.Sprintf(updateRun, run.UserId, run.ProjectId, buildId, run.Id)
	runBytes, err := json.Marshal(run)
	if err != nil {
		return nil, nil, err
	}
	return util.PatchHttpRespWithData(url, nil, runBytes)
}

func (p *ProjectStoreClient) PutLog(userId, projectId, buildId, runId, runLog string) (*http.Response, []byte, error) {
	url := "http://" + p.HostPort + fmt.Sprintf(putLog, userId, projectId, buildId, runId)

	return util.PutHttpRespWithData(url, nil, []byte(runLog))
}

func (p *ProjectStoreClient) GetProject(userId string, projectId string) (project api.Project, err error) {
	url := "http://" + p.HostPort + fmt.Sprintf(getProject, userId, projectId)
	var resp *http.Response
	var raw []byte

	resp, raw, err = util.GetHttpResp(url, nil)
	if resp.StatusCode != http.StatusOK {
		err = errors.New("response is not 200!" + strconv.Itoa(resp.StatusCode))
		return
	}
	err = json.Unmarshal(raw, &project)
	return
}
