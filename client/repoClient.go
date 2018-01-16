package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/qetuantuan/jengo_recap/model"
	"github.com/qetuantuan/jengo_recap/util"
)

const (
	getrepo = "/v0.1/user/%s/repo/%s"
)

type HttpRepoService struct {
	HostPort string
}

func NewHttpRepoService(host string, port int) (uc *HttpRepoService) {
	uc = &HttpRepoService{
		HostPort: util.GetHostPort4Client(host, port),
	}

	return
}

func (p *HttpRepoService) GetRepo(userId string, repoId string) (repo model.Repo, err error) {
	url := "http://" + p.HostPort + fmt.Sprintf(getrepo, userId, repoId)
	var resp *http.Response
	var raw []byte

	resp, raw, err = util.GetHttpResp(url, nil)
	if resp.StatusCode != http.StatusOK {
		err = errors.New("response is not 200!" + strconv.Itoa(resp.StatusCode))
		return
	}
	err = json.Unmarshal(raw, &repo)
	return
}

func GetReposByFilter(filter map[string]interface{}, limitCount, offset int) (Repos []model.Repo, err error) {
	return
}
