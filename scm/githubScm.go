package scm

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/qetuantuan/jengo_recap/model"
)

type GithubScm struct {
	baseScm
	scm string
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

type GithubError struct {
	Resource string
	Code     string
	Message  string
}
type GithubHookError struct {
	Message string
	Errors  []GithubError
}

const WEBHOOKNAME string = "web"
const HOOK_EXIST_ERROR_MESSAGE string = "Hook already exists on this repository"

var HookExistError = errors.New("hook already exists")
var HookNonExistError = errors.New("hook does not exists")

func (gs *GithubScm) SetGatewayHookUrl(gatewayHookUrl string) {
	gs.HookURI = gatewayHookUrl
}

func NewGithubScm(hookUrl string) *GithubScm {
	gs := GithubScm{scm: "github"}
	gs.ApiLink = "https://api.github.com"
	gs.HookURI = hookUrl
	return &gs
}

const hookConfFormatString = "{\"name\": \"%s\", " +
	"\"active\": true, " +
	"\"events\": [\"push\"], " +
	"\"config\": {\"url\": \"%s\",\"content_type\": \"json\"}}"

const hookEditConfFormatString = "{\"name\": \"%s\", " +
	"\"active\": true, " +
	"\"events\": [\"push\"], " +
	"\"config\": {\"url\": \"%s\",\"content_type\": \"json\"}, " +
	"\"add_events\": [\"\"], " +
	"\"remove_events\": [\"\"]}"

func (gs *GithubScm) SetHook(projectName string) (hook model.GithubHook, err error) {
	hookConf := fmt.Sprintf(hookConfFormatString, WEBHOOKNAME, gs.HookURI)
	uri := fmt.Sprintf("%s/repos/%s/%s/hooks", gs.ApiLink, gs.User, projectName)

	byteResponses, _, err := gs.httpRequest(
		"POST",
		uri,
		hookConf,
		map[string]string{
			"Authorization": "token " + gs.Token,
			// https://developer.github.com/v3/#timezones
			"Time-Zone": "GMT",
		})
	if err != nil {
		githubHookError := GithubHookError{}
		errTmp := json.Unmarshal(byteResponses, &githubHookError)
		if errTmp != nil || len(githubHookError.Errors) < 1 || githubHookError.Errors[0].Message != HOOK_EXIST_ERROR_MESSAGE {
			return
		}
		err = HookExistError
		return
	}
	err = json.Unmarshal(byteResponses, &hook)
	return
}

// keep edit recently ,but not used now
func (gs *GithubScm) EditHook(hookUrl string) (hook model.GithubHook, err error) {
	hookConf := fmt.Sprintf(hookEditConfFormatString, WEBHOOKNAME, gs.HookURI)
	uri := fmt.Sprintf("%s", hookUrl)
	byteResponses, _, err := gs.httpRequest("PATCH", uri, hookConf, map[string]string{"Authorization": "token " + gs.Token})
	if err != nil {
		githubHookError := GithubHookError{}
		errTmp := json.Unmarshal(byteResponses, &githubHookError)
		if errTmp != nil || len(githubHookError.Errors) < 1 || githubHookError.Errors[0].Message != HOOK_EXIST_ERROR_MESSAGE {
			return
		}
		err = HookExistError
		return
	}
	err = json.Unmarshal(byteResponses, &hook)
	return
}

func (gs *GithubScm) GetHook(url string) (hook model.GithubHook, err error) {
	byteResponses, _, err := gs.httpRequest("GET", url, "", map[string]string{"Authorization": "token " + gs.Token})
	if err != nil {
		return
	}
	hooks := []model.GithubHook{}
	err = json.Unmarshal(byteResponses, &hooks)
	if err != nil {
		return
	}
	//Todo : one project may have multiple hooks related to our system, but now we restrict only one
	// should we define gitHook.Config.Url more clear? something like := gs.HookURI+"/webhook" or other...
	for _, gitHook := range hooks {
		if gitHook.Config.Url == gs.HookURI &&
			gitHook.Config.ContentType == "json" &&
			gitHook.Config.InsecureSsl == "0" &&
			len(gitHook.Events) == 1 && gitHook.Events[0] == "push" {
			//use config and event to recognize the hook ,see https://developer.github.com/v3/repos/hooks/
			hook = gitHook
			return
		}
	}
	err = HookNonExistError
	return
}

// for http test
func (gs *GithubScm) SetApiLink(apilink string) {
	gs.ApiLink = apilink
}

func (gs *GithubScm) DeleteHook(url string) (err error) {
	_, _, err = gs.httpRequest("DELETE", url, "", map[string]string{"Authorization": "token " + gs.Token})
	return
}

func (gs *GithubScm) GetProjectList() (projects []model.Project, err error) {
	//uri := fmt.Sprintf("%s/users/%s/repos", gs.ApiLink, gs.User)
	uri := fmt.Sprintf("%s/user/repos", gs.ApiLink)
	byteResponses, _, err := gs.httpRequest(
		"GET",
		uri,
		"",
		map[string]string{
			"Authorization": "token " + gs.Token,
			"Time-Zone":     "GMT",
		})
	if err != nil {
		return
	}
	githubProjects := []GithubProject{}
	err = json.Unmarshal(byteResponses, &githubProjects)
	if err != nil {
		return
	}
	projects = make([]model.Project, len(githubProjects))
	for i := range githubProjects {
		githubProjects[i].CopyTo(&projects[i])
	}
	return
}
