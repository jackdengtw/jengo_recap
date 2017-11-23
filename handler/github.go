package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/qetuantuan/jengo_recap/api"
	"github.com/qetuantuan/jengo_recap/scm"
	"github.com/qetuantuan/jengo_recap/util"

	"github.com/golang/glog"
	"github.com/google/go-github/github"
)

// GitHubWebHookHandler defines the how to handle received webhook payload
// from github server.
//

type GitHubHandler struct {
	BaseHandler
	Operator
	Trigger
}

func NewGitHubHandler(server string, port int) *GitHubHandler {
	h := &GitHubHandler{
		BaseHandler: BaseHandler{
			name:    "github_webhook",
			method:  "POST",
			pattern: "/github_webhook",
		},
		Operator: &GitHubOperator{},
		Trigger:  NewGitHubTrigger(server, port),
	}
	h.handlerFunc = h.Handle
	Register(h)
	return h
}

func (h *GitHubHandler) ParseRequest(w http.ResponseWriter, r *http.Request) (
	eventWrapper *scm.GithubEventWrapper, err error) {
	glog.Infof("starting\n")

	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		glog.Errorf("error reading request body: err(%v)\n", err)
		return
	}
	glog.Infof("succeeded to read request body: payload(%v)\n", payload)

	event, err := h.ParseWebHook(h.WebHookType(r), payload)
	if err != nil {
		appErr := "could not parse webhook to get X-GitHub-Event: " + err.Error()
		glog.Error(appErr)
		util.ReturnError(http.StatusBadRequest, appErr, w)
		return
	}
	glog.Infof("get X-Github-Event(%v)", event)

	deliveryId := h.DeliveryID(r)
	if deliveryId == "" {
		appErr := "could not get X-GitHub-Delivery"
		glog.Errorf(appErr)
		util.ReturnError(http.StatusBadRequest, appErr, w)
		return
	}
	glog.Infof("get X-Github-Delivery(%s)", deliveryId)
	// TODO: id generation may move to util later
	eventId := fmt.Sprintf("e_github_%s", deliveryId)

	signature := r.Header.Get("X-Hub-Signature")
	if signature == "" {
		appErr := "TODO: could not get X-Hub-Signature"
		glog.Errorf(appErr)
		// TODO: validate signature?
	}
	glog.Infof("get X-Hub-Signature(%s)", signature)

	glog.Infof("ending: payload(\n%v\n), event(\n%v\n), eventId(\n%v\n), "+
		"signature(\n%v\n)\n", payload, event, eventId, signature)

	eventWrapper = &scm.GithubEventWrapper{}
	eventWrapper.Event = event
	eventWrapper.EventId = eventId
	eventWrapper.Payload = payload
	eventWrapper.Signature = signature

	return
}

func (h *GitHubHandler) Handle(w http.ResponseWriter, r *http.Request) {
	eventWrapper, err := h.ParseRequest(w, r)
	if err != nil {
		glog.Errorf("failed to ParseRequest: err(%v)\n", err.Error())
		return
	}

	switch e := eventWrapper.Event.(type) {
	case *github.PushEvent:
		glog.Info("event type PushEvent\n")
		glog.Infof("event(\n%v\n)\n", eventWrapper.Event)
		util.ReturnSuccessWithMap(200, map[string]string{"event_id": eventWrapper.EventId}, w)
		if _, err := h.ProcessWithPushEvent(e, eventWrapper.Payload, eventWrapper.EventId); err != nil {
			glog.Errorf("failed to call processWithPushEvent: err(%v)\n", err)
			return
		}
		return
	case *github.PullRequestEvent:
		glog.Info("event type PullRequestEvent\n")
		glog.Infof("event(\n%v\n)\n", eventWrapper.Event)
		util.ReturnError(http.StatusBadRequest, eventWrapper.EventId, w)
		return
	case *github.PingEvent:
		glog.Info("event type PingEvent\n")
		glog.Infof("event(\n%v\n)\n", eventWrapper.Event)
		util.ReturnSuccess(http.StatusBadRequest, eventWrapper.EventId, "", w)
		return
	case *github.WatchEvent:
		if e.Action != nil && *e.Action == "starred" {
			glog.Errorf("%s starred repository %s\n",
				*e.Sender.Login, *e.Repo.FullName)
		}
		glog.Infof("event(\n%v\n)\n", eventWrapper.Event)
		util.ReturnError(http.StatusBadRequest, eventWrapper.EventId, w)
		return
	default:
		glog.Infof("unknown event type(%s)\n",
			h.WebHookType(r))
		appErr := "received unknown event type"
		glog.Infof(appErr)
		util.ReturnError(http.StatusBadRequest, appErr, w)
		return
	}
}

func (h *GitHubHandler) ProcessWithPushEvent(e *github.PushEvent, payload []byte, eventId string) (
	run *api.EngineCreateRunParams, err error) {
	glog.Infof("starting\n")

	if err = h.validatePushEvent(e); err != nil {
		glog.Infof("validate error: err(%s)\n", err.Error())
		return
	}

	// init run
	run = &api.EngineCreateRunParams{}

	// assign EventId
	run.EventId = eventId

	// assign ProjectId
	run.ProjectId, _ = h.getProjectId(*e.Repo.ID)

	// assign UserId
	run.UserId, err = h.getUserId(payload)
	if err != nil {
		glog.Errorf("failed to call getUserId: err(%s)\n", err)
		return
	}

	// assign Branch
	run.Branch, err = h.getBranch(e)
	if err != nil {
		glog.Errorf("failed to call getBranch: err(%s)\n", err)
		return
	}

	// assign Commits
	run.Commits, err = h.getCommits(e)
	if err != nil {
		glog.Errorf("failed to get commits: err(%s)\n", err)
	}

	// assign Head Commit
	run.HeadCommit = h.getHeadCommit(e)

	// assign Repo
	run.Repo, err = h.getRepo(e)
	if err != nil {
		glog.Errorf("failed to get repo: err(%s)\n", err)
	}

	// assign Compare
	run.Compare = *e.Compare

	// trigger create run request
	httpResp, runResp, err := h.TriggerRun(run)
	if err != nil {
		glog.Errorf("failed to get repo owner id: err(%s)\n", err)
	}

	glog.Infof("ending: httpResp(%v), runResp(%v)\n", httpResp, runResp)
	return
}

func (h *GitHubHandler) validatePushEvent(e *github.PushEvent) (err error) {
	glog.Infof("starting\n")

	if e != nil && e.Repo != nil && e.Repo.ID != nil && e.Compare != nil {
		glog.Infof("looks good: *e.Repo.ID(%s), *e.Compare\n",
			*e.Repo.ID, *e.Compare)
		err = nil
	} else {
		err = fmt.Errorf("nil error: either e is nil or e.Repo is nil or e.Repo.ID is nil")
		glog.Info("nil error: err(%s)\n", err.Error())
	}

	// TODO: continue to validate other used params
	// TODO: validate project and user info from service

	glog.Infof("ending\n")
	return
}

func (h *GitHubHandler) getRepoOwnerId(payload []byte) (repoOwnerId *int, err error) {
	glog.Infof("starting\n")

	if len(payload) == 0 {
		err = fmt.Errorf("0 length error: the length of payload is 0")
		glog.Info("0 length error: err(%s)\n", err.Error())
	}

	githubEvent := &scm.GithubPushEvent{}

	err = json.Unmarshal(payload, githubEvent)
	if err != nil {
		glog.Errorf("failed to call json.Unmarshal: err(%s), payload(%s)\n",
			err, payload)
		return
	}

	repoOwnerId = githubEvent.Repo.Owner.Id
	if repoOwnerId == nil {
		err = fmt.Errorf("nil error: get nil repoOwnerId")
		glog.Errorf("failed to get repo owner id: err(%s), payload(%s)\n",
			err, payload)
		return
	}

	glog.Infof("ending: repoOwnerId(%v)\n", repoOwnerId)
	return
}

func (h *GitHubHandler) getBranch(e *github.PushEvent) (branch string, err error) {
	glog.Infof("starting\n")

	if e != nil && e.Ref != nil {
		glog.Infof("looks good: *e.Ref(%s)\n", *e.Ref)
	} else {
		err = fmt.Errorf("nil error: either e is nil or e.Ref is nil")
		glog.Info("nil error: err(%s)\n", err.Error())
	}

	strList := strings.Split(*e.Ref, "/")
	branch = strings.Trim(strList[len(strList)-1], "")

	glog.Infof("ending: branch(%s)\n", branch)
	return
}

func (h *GitHubHandler) generateEventId(deliveryId string) (eventId string, err error) {
	glog.Infof("starting\n")

	eventId = fmt.Sprintf("event_%s", deliveryId)

	glog.Infof("ending: eventId(%s)\n", eventId)
	return
}

func (h *GitHubHandler) getCommits(e *github.PushEvent) (commits []api.PushEventCommit, err error) {
	glog.Infof("starting\n")

	// TODO: data structure convert needed instead of Marshal and Unmarshal
	raw, err := json.Marshal(e.Commits)
	if err != nil {
		glog.Errorf("Marshal error: err(%s)\n", err.Error())
	}
	glog.Infof("raw(%s)\n", string(raw))

	err = json.Unmarshal(raw, &commits)
	if err != nil {
		glog.Errorf("Unmarshal error: err(%s)\n", err.Error())
	}

	glog.Infof("ending: commits(%v)\n", commits)
	return
}

func (h *GitHubHandler) getHeadCommit(e *github.PushEvent) (headCommit *api.PushEventCommit) {
	glog.Infof("starting\n")

	if e.HeadCommit == nil {
		return
	}

	headCommit = &api.PushEventCommit{
		Message: e.HeadCommit.Message,
		Author: &api.CommitAuthor{
			Date:  e.HeadCommit.Author.Date,
			Name:  e.HeadCommit.Author.Name,
			Email: e.HeadCommit.Author.Email,
			Login: e.HeadCommit.Author.Login,
		},
		URL:      e.HeadCommit.URL,
		Distinct: e.HeadCommit.Distinct,
		SHA:      e.HeadCommit.SHA,
		ID:       e.HeadCommit.ID,
		TreeID:   e.HeadCommit.TreeID,
		Timestamp: &api.Timestamp{
			Time: e.HeadCommit.Timestamp.Time,
		},
		Committer: &api.CommitAuthor{
			Date:  e.HeadCommit.Committer.Date,
			Name:  e.HeadCommit.Committer.Name,
			Email: e.HeadCommit.Committer.Email,
			Login: e.HeadCommit.Committer.Login,
		},
		Added:    e.HeadCommit.Added,
		Removed:  e.HeadCommit.Removed,
		Modified: e.HeadCommit.Modified,
	}

	glog.Infof("ending: head_commit(%v)\n", headCommit)
	return
}

func (h *GitHubHandler) getRepo(e *github.PushEvent) (repo *api.PushEventRepository, err error) {
	glog.Infof("starting\n")

	// TODO: data structure convert needed instead of Marshal and Unmarshal
	raw, err := json.Marshal(e.Repo)
	if err != nil {
		glog.Errorf("Marshal error: err(%s)\n", err.Error())
	}
	glog.Infof("raw(%s)\n", string(raw))

	err = json.Unmarshal(raw, &repo)
	if err != nil {
		glog.Errorf("Unmarshal error: err(%s)\n", err.Error())
	}

	glog.Infof("ending: repo(%v)\n", repo)
	return
}

func (h *GitHubHandler) getProjectId(repoId int) (projectId string, err error) {
	glog.Infof("starting\n")

	projectId = "p_github_" + strconv.Itoa(repoId)

	glog.Infof("ending: projectId(%v)\n", projectId)
	return
}

func (h *GitHubHandler) getUserId(payload []byte) (userId string, err error) {
	glog.Infof("starting\n")

	repoOwnerId, err := h.getRepoOwnerId(payload)
	if err != nil {
		glog.Errorf("failed to call getRepoOwnerId: err(%s)\n", err.Error())
		return
	}

	if repoOwnerId == nil {
		err = fmt.Errorf("repoOwnerId is nil")
		glog.Errorf("failed to call getRepoOwnerId: err(%s)\n", err.Error())
		return
	}

	userId = "u_github_" + strconv.Itoa(*repoOwnerId)

	glog.Infof("ending: userId(%v)\n", userId)
	return
}
