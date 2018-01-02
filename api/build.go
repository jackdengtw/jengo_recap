package api

import (
	"time"
)

// distincted by branch + commitId for a repo
// in favor of end user
type SemanticBuild struct {
	// Id is hash value from userId, repoId, branch and commitId
	Id        string `json:"id" bson:"_id"`
	ProjectId string `json:"project_id"`
	UserId    string `json:"user_id"`
	CommitId  string `json:"commit_id"`
	Branch    string `json:"branch"`

	Numero *int `json:"numero,omitempty"`

	Builds []Build `json:"builds"`
}

type SemanticBuilds []SemanticBuild

type Build struct {
	Id     string `json:"id" bson:"_id"`
	Status string `json:"status"`
	Result string `json:"result"`

	EventId    *string           `json:"event_id,omitempty"`
	Commits    []PushEventCommit `json:"commits"`
	HeadCommit *PushEventCommit  `json:"head_commit,omitempty"`

	LogUri    *string        `json:"log_url,omitempty"`
	UpdatedAt *time.Time     `json:"updated_at,omitempty"`
	CreatedAt *time.Time     `json:"created_at,omitempty"`
	StartTime *time.Time     `json:"start_time,omitempty"`
	Duration  *time.Duration `json:"duration,omitempty"`
}

type Builds []Build
