package scm

import (
	"encoding/json"
	"fmt"

	"github.com/qetuantuan/jengo_recap/model"
)

type GithubScm struct {
	baseScm
	scm string
}

func NewGithubScm() *GithubScm {
	gs := GithubScm{scm: "github.com"}
	gs.ApiLink = "https://api.github.com"
	return &gs
}

func (gs *GithubScm) GetGithubUser(token string) (user model.GithubUser, err error) {
	uri := fmt.Sprintf("%s/user", gs.ApiLink)
	byte_response, _, err := gs.httpRequest("GET", uri, "", map[string]string{"Authorization": "token " + token})
	if err != nil {
		return
	}

	err = json.Unmarshal(byte_response, &user)
	return
}
