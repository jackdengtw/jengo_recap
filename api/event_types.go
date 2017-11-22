package api

import "time"

// PushEventCommit represents a git commit in a GitHub PushEvent.
type PushEventCommit struct {
	Message  *string       `json:"message,omitempty"`
	Author   *CommitAuthor `json:"author,omitempty"`
	URL      *string       `json:"url,omitempty"`
	Distinct *bool         `json:"distinct,omitempty"`

	// The following fields are only populated by Events API.
	SHA *string `json:"sha,omitempty"`

	// The following fields are only populated by Webhook events.
	ID        *string       `json:"id,omitempty"`
	TreeID    *string       `json:"tree_id,omitempty"`
	Timestamp *Timestamp    `json:"timestamp,omitempty"`
	Committer *CommitAuthor `json:"committer,omitempty"`
	Added     []string      `json:"added,omitempty"`
	Removed   []string      `json:"removed,omitempty"`
	Modified  []string      `json:"modified,omitempty"`
}

// CommitAuthor represents the author or committer of a commit. The commit
// author may not correspond to a GitHub User.
type CommitAuthor struct {
	Date  *time.Time `json:"date,omitempty"`
	Name  *string    `json:"name,omitempty"`
	Email *string    `json:"email,omitempty"`

	// The following fields are only populated by Webhook events.
	Login *string `json:"username,omitempty"` // Renamed for go-github consistency.
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
