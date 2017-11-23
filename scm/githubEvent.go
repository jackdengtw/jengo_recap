package scm

import (
	. "github.com/google/go-github/github"
)

// PushEventRepoOwner is a basic representation of user/org in a PushEvent payload.
type GithubPushEventRepoOwner struct {
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty"`
	Login *string `json:"login,omitempty"`
	Id    *int    `json:"id,omitempty"`
}

// PushEventRepository represents the repo object in a PushEvent payload.
type GithubPushEventRepository struct {
	ID              *int                      `json:"id,omitempty"`
	Name            *string                   `json:"name,omitempty"`
	FullName        *string                   `json:"full_name,omitempty"`
	Owner           *GithubPushEventRepoOwner `json:"owner,omitempty"`
	Private         *bool                     `json:"private,omitempty"`
	Description     *string                   `json:"description,omitempty"`
	Fork            *bool                     `json:"fork,omitempty"`
	CreatedAt       *Timestamp                `json:"created_at,omitempty"`
	PushedAt        *Timestamp                `json:"pushed_at,omitempty"`
	UpdatedAt       *Timestamp                `json:"updated_at,omitempty"`
	Homepage        *string                   `json:"homepage,omitempty"`
	Size            *int                      `json:"size,omitempty"`
	StargazersCount *int                      `json:"stargazers_count,omitempty"`
	WatchersCount   *int                      `json:"watchers_count,omitempty"`
	Language        *string                   `json:"language,omitempty"`
	HasIssues       *bool                     `json:"has_issues,omitempty"`
	HasDownloads    *bool                     `json:"has_downloads,omitempty"`
	HasWiki         *bool                     `json:"has_wiki,omitempty"`
	HasPages        *bool                     `json:"has_pages,omitempty"`
	ForksCount      *int                      `json:"forks_count,omitempty"`
	OpenIssuesCount *int                      `json:"open_issues_count,omitempty"`
	DefaultBranch   *string                   `json:"default_branch,omitempty"`
	MasterBranch    *string                   `json:"master_branch,omitempty"`
	Organization    *string                   `json:"organization,omitempty"`
	URL             *string                   `json:"url,omitempty"`
	ArchiveURL      *string                   `json:"archive_url,omitempty"`
	HTMLURL         *string                   `json:"html_url,omitempty"`
	StatusesURL     *string                   `json:"statuses_url,omitempty"`
	GitURL          *string                   `json:"git_url,omitempty"`
	SSHURL          *string                   `json:"ssh_url,omitempty"`
	CloneURL        *string                   `json:"clone_url,omitempty"`
	SVNURL          *string                   `json:"svn_url,omitempty"`
}

// PushEvent represents a git push to a GitHub repository.
//
// GitHub API docs: https://developer.github.com/v3/activity/events/types/#pushevent
type GithubPushEvent struct {
	PushID       *int              `json:"push_id,omitempty"`
	Head         *string           `json:"head,omitempty"`
	Ref          *string           `json:"ref,omitempty"`
	Size         *int              `json:"size,omitempty"`
	Commits      []PushEventCommit `json:"commits,omitempty"`
	Before       *string           `json:"before,omitempty"`
	DistinctSize *int              `json:"distinct_size,omitempty"`

	// The following fields are only populated by Webhook events.
	After        *string                    `json:"after,omitempty"`
	Created      *bool                      `json:"created,omitempty"`
	Deleted      *bool                      `json:"deleted,omitempty"`
	Forced       *bool                      `json:"forced,omitempty"`
	BaseRef      *string                    `json:"base_ref,omitempty"`
	Compare      *string                    `json:"compare,omitempty"`
	Repo         *GithubPushEventRepository `json:"repository,omitempty"`
	HeadCommit   *PushEventCommit           `json:"head_commit,omitempty"`
	Pusher       *User                      `json:"pusher,omitempty"`
	Sender       *User                      `json:"sender,omitempty"`
	Installation *Installation              `json:"installation,omitempty"`
}

type GithubEventWrapper struct {
	Event     interface{}
	EventId   string
	Payload   []byte
	Signature string
}
