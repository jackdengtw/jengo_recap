package handler

import (
	"net/http"
	"testing"

	"github.com/qetuantuan/jengo_recap/api"

	"fmt"

	"github.com/golang/glog"
	"github.com/google/go-github/github"
)

type MockGitHubTrigger struct {
}

func (t *MockGitHubTrigger) TriggerRun(runReq *api.EngineCreateRunParams) (
	httpResp *http.Response, runResp *api.Run, err error) {
	glog.Infof("starting\n")

	httpResp = nil
	runResp = nil
	err = nil

	glog.Infof("ending\n")
	return
}

func TestGitHubHandlerProcessWithPushEvent(t *testing.T) {

	var gHandler *GitHubHandler = NewGitHubHandler("127.0.0.1", 8088)
	gHandler.Trigger = &MockGitHubTrigger{}

	repoId := 11111111
	branch := "branch_name"
	compare := "event_compare_info"
	ownerId := 22222222
	eventId := "event_33333333"

	r := &github.PushEventRepository{
		ID: &repoId,
	}
	ref := "refs/heads/" + branch
	e := &github.PushEvent{
		Ref:     &ref,
		Repo:    r,
		Compare: &compare,
	}
	payload := `{
	"repository": {
		"owner": {
			"name": "jengo_owner",
			"email": "jengo@jengo.com",
			"login": "jengo_owner",
			"id": %d
		}
	}
}`
	payload = fmt.Sprintf(payload, ownerId)

	run, err := gHandler.ProcessWithPushEvent(e, ([]byte)(payload), eventId)
	if err != nil {
		t.Fatalf("failed to call ProcessWithPushEvent: err(%v)", err)
	}
	projectId := fmt.Sprintf("p_github_%d", repoId)
	if run.ProjectId != projectId {
		t.Fatalf("Project not equal: expected(%s), got(%s)", projectId, run.ProjectId)
	}
	userId := fmt.Sprintf("u_github_%d", ownerId)
	if run.UserId != userId {
		t.Fatalf("UserId not equal: expected(%s), got(%s)", userId, run.UserId)
	}
	if run.Branch != branch {
		t.Fatalf("Branch not equal: expected(%s), got(%s)", branch, run.Branch)
	}
	if run.Compare != compare {
		t.Fatalf("Compare not equal: expected(%s), got(%s)", compare, run.Compare)
	}

}
