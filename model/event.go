package model

import (
	"time"

	"github.com/qetuantuan/jengo_recap/api"
)

// PushEventCommit represents a git commit in a GitHub PushEvent.
type PushEventCommit struct {
	api.PushEventCommit

	// The following fields are only populated by Events API.
	SHA *string `json:"sha,omitempty"`

	// The following fields are only populated by Webhook events.
	Added    []string `json:"added,omitempty"`
	Removed  []string `json:"removed,omitempty"`
	Modified []string `json:"modified,omitempty"`
}

func (p *PushEventCommit) ToApiObj() *api.PushEventCommit {
	return &p.PushEventCommit
}

func NewPushEventCommit(p api.PushEventCommit) *PushEventCommit {
	return &PushEventCommit{PushEventCommit: p}
}

type PushEventCommits []PushEventCommit

func (ps *PushEventCommits) ToApiObj() (aps api.PushEventCommits) {
	for _, p := range *ps {
		aps = append(aps, *p.ToApiObj())
	}
	return
}

func NewPushEventCommits(aps api.PushEventCommits) (ps PushEventCommits) {
	for _, p := range aps {
		ps = append(ps, *NewPushEventCommit(p))
	}
	return
}

// Timestamp represents a time that can be unmarshalled from a JSON string
// formatted as either an RFC3339 or Unix timestamp. This is necessary for some
// fields since the GitHub API is inconsistent in how it represents times. All
// exported methods of time.Time can be called on Timestamp.
type Timestamp struct {
	time.Time
}

// PushEventRepository represents the repo object in a PushEvent payload.
type PushEventRepository struct {
	ID              *int                `json:"id,omitempty"`
	Name            *string             `json:"name,omitempty"`
	FullName        *string             `json:"full_name,omitempty"`
	Owner           *PushEventRepoOwner `json:"owner,omitempty"`
	Private         *bool               `json:"private,omitempty"`
	Description     *string             `json:"description,omitempty"`
	Fork            *bool               `json:"fork,omitempty"`
	CreatedAt       *Timestamp          `json:"created_at,omitempty"`
	PushedAt        *Timestamp          `json:"pushed_at,omitempty"`
	UpdatedAt       *Timestamp          `json:"updated_at,omitempty"`
	Homepage        *string             `json:"homepage,omitempty"`
	Size            *int                `json:"size,omitempty"`
	StargazersCount *int                `json:"stargazers_count,omitempty"`
	WatchersCount   *int                `json:"watchers_count,omitempty"`
	Language        *string             `json:"language,omitempty"`
	HasIssues       *bool               `json:"has_issues,omitempty"`
	HasDownloads    *bool               `json:"has_downloads,omitempty"`
	HasWiki         *bool               `json:"has_wiki,omitempty"`
	HasPages        *bool               `json:"has_pages,omitempty"`
	ForksCount      *int                `json:"forks_count,omitempty"`
	OpenIssuesCount *int                `json:"open_issues_count,omitempty"`
	DefaultBranch   *string             `json:"default_branch,omitempty"`
	MasterBranch    *string             `json:"master_branch,omitempty"`
	Organization    *string             `json:"organization,omitempty"`
	URL             *string             `json:"url,omitempty"`
	ArchiveURL      *string             `json:"archive_url,omitempty"`
	HTMLURL         *string             `json:"html_url,omitempty"`
	StatusesURL     *string             `json:"statuses_url,omitempty"`
	GitURL          *string             `json:"git_url,omitempty"`
	SSHURL          *string             `json:"ssh_url,omitempty"`
	CloneURL        *string             `json:"clone_url,omitempty"`
	SVNURL          *string             `json:"svn_url,omitempty"`
}

// PushEventRepoOwner is a basic representation of user/org in a PushEvent payload.
type PushEventRepoOwner struct {
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty"`
}
