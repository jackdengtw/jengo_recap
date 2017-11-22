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
	GET_USER = "/v0.2/internal_user/%s"
)

type UserStoreClientInterface interface {
	GetUser(userId string) (user api.User02, err error)
}

type UserStoreClient struct {
	HostPort string
}

func NewUserStoreClient(hostPort string) (uc *UserStoreClient) {
	uc = &UserStoreClient{
		HostPort: hostPort,
	}
	return
}

func (c *UserStoreClient) GetUser(userId string) (user api.User02, err error) {
	var resp *http.Response
	var raw []byte
	url := "http://" + c.HostPort + fmt.Sprintf(GET_USER, userId)

	resp, raw, err = util.GetHttpResp(url, nil)
	if resp.StatusCode != http.StatusOK {
		err = errors.New("response is not 200!" + strconv.Itoa(resp.StatusCode))
		return
	}
	err = json.Unmarshal(raw, &user)
	return
}
