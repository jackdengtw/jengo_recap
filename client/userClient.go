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
	GET_USER = "/v0.2/internal_user/%s"
)

type HttpUserService struct {
	HostPort string
}

func NewHttpUserService(host string, port int) (uc *HttpUserService) {
	uc = &HttpUserService{
		HostPort: util.GetHostPort4Client(host, port),
	}
	return
}

func (c *HttpUserService) GetUser(userId string) (user model.User, err error) {
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
