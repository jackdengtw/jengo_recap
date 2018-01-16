package api

import "time"

// PushEventCommit is a subset of model.PushEventCommit exposing to end user
type PushEventCommit struct {
	Message  *string       `json:"message,omitempty"`
	Author   *CommitAuthor `json:"author,omitempty"`
	URL      *string       `json:"url,omitempty"`
	Distinct *bool         `json:"distinct,omitempty"`

	// The following fields are only populated by Webhook events.
	ID        *string       `json:"id,omitempty"`
	TreeID    *string       `json:"tree_id,omitempty"`
	Timestamp *time.Time    `json:"timestamp,omitempty"`
	Committer *CommitAuthor `json:"committer,omitempty"`
}

type PushEventCommits []PushEventCommit

// CommitAuthor is a subset of model.CommitAuther exposing to end user
type CommitAuthor struct {
	Date  *time.Time `json:"date,omitempty"`
	Name  *string    `json:"name,omitempty"`
	Email *string    `json:"email,omitempty"`

	// The following fields are only populated by Webhook events.
	Login *string `json:"username,omitempty"` // Renamed for go-github consistency.
}

/*
// PushEventRepository represents the repo object in a PushEvent payload.
type PushEventRepository struct {
	ID           *int                `json:"id,omitempty"`
	Name         *string             `json:"name,omitempty"`
	FullName     *string             `json:"full_name,omitempty"`
	Owner        *PushEventRepoOwner `json:"owner,omitempty"`
	Organization *string             `json:"organization,omitempty"`
	URL          *string             `json:"url,omitempty"`
	HTMLURL      *string             `json:"html_url,omitempty"`
}

// PushEventRepoOwner is a basic representation of user/org in a PushEvent payload.
type PushEventRepoOwner struct {
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty"`
}
*/
