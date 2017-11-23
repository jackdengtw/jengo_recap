package handler

import (
	"net/http"

	"github.com/google/go-github/github"
)

type Operator interface {
	WebHookType(r *http.Request) string
	ParseWebHook(messageType string, payload []byte) (interface{}, error)
	DeliveryID(r *http.Request) string
}

type GitHubOperator struct {
}

func (o *GitHubOperator) WebHookType(r *http.Request) string {
	return github.WebHookType(r)
}

func (o *GitHubOperator) ParseWebHook(messageType string, payload []byte) (interface{}, error) {
	return github.ParseWebHook(messageType, payload)
}

func (o *GitHubOperator) DeliveryID(r *http.Request) string {
	return github.DeliveryID(r)
}
