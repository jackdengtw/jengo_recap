package api

import (
	"time"
)

// distincted by branch + commitId for a repo
// in favor of end user
type SemanticBuild struct {
	// Id is hash value from repoId, branch and commitId
	Id       string `json:"id" bson:"_id"`
	RepoId   string `json:"repo_id"`
	CommitId string `json:"commit_id"`
	Branch   string `json:"branch"`

	UserId string `json:"user_id"`

	Numero *int `json:"numero,omitempty"`

	Builds Builds `json:"builds"`
}

type SemanticBuilds []SemanticBuild

type Build struct {
	Id     string `json:"id" bson:"_id"`
	Status string `json:"status"`
	Result string `json:"result"`

	// duplicate info in SemanticBuild for now. Later on remove Semantic Build
	RepoId   string `json:"repo_id"`
	CommitId string `json:"commit_id"`
	Branch   string `json:"branch"`
	UserId   string `json:"user_id"`

	EventId    *string          `json:"event_id,omitempty"`
	Commits    PushEventCommits `json:"commits"`
	HeadCommit *PushEventCommit `json:"head_commit,omitempty"`

	LogUri    *string        `json:"log_url,omitempty"`
	UpdatedAt *time.Time     `json:"updated_at,omitempty"`
	CreatedAt *time.Time     `json:"created_at,omitempty"`
	StartTime *time.Time     `json:"start_time,omitempty"`
	Duration  *time.Duration `json:"duration,omitempty"`
}

type Builds []Build
