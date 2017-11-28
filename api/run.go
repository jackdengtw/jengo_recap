package api

import (
	"time"
)

type Run struct {
	Id         string            `bson:"_id" json:"id"`
	EventId    string            `json:"event_id"`
	UserId     string            `json:"user_id"`
	ProjectId  string            `json:"project_id"`
	Branch     string            `json:"branch"`
	Status     string            `json:"status"`
	State      string            `json:"state"`
	Commits    []PushEventCommit `json:"commits"`
	HeadCommit *PushEventCommit  `json:"head_commit"`
	Compare    string            `json:"compare"`
	StartTime  *time.Time        `json:"start_time"`
	Duration   time.Duration     `json:"duration"`
	Config     string            `json:"config"`
	LogId      string            `json:"log_id"`
	RunEnv     string            `json:"run_env"`
	UpdatedAt  *time.Time        `json:"updated_at"`
	CreatedAt  *time.Time        `json:"created_at"`
}

type Runs []Run

type PatchRun struct {
	Id         string            `bson:"_id" json:"id"`
	EventId    string            `json:"event_id,omitempty"`
	UserId     string            `json:"user_id,omitempty"`
	ProjectId  string            `json:"project_id,omitempty"`
	Branch     string            `json:"branch,omitempty"`
	Status     string            `json:"status,omitempty"`
	State      string            `json:"state,omitempty"`
	Commits    []PushEventCommit `json:"commits,omitempty"`
	HeadCommit *PushEventCommit  `json:"head_commit,omitempty"`
	Compare    string            `json:"compare,omitempty"`
	StartTime  *time.Time        `json:"start_time,omitempty"`
	Duration   time.Duration     `json:"duration,omitempty"`
	Config     string            `json:"config,omitempty"`
	LogId      string            `json:"log_id,omitempty"`
	RunEnv     string            `json:"run_env,omitempty"`
	UpdatedAt  *time.Time        `json:"updated_at,omitempty"`
	CreatedAt  *time.Time        `json:"created_at,omitempty"`
}

type EngineCreateRunParams struct {
	ProjectId  string               `json:"project_id"`
	UserId     string               `json:"user_id"`
	EventId    string               `json:"event_id"`
	Branch     string               `json:"branch"`
	Commits    []PushEventCommit    `json:"commits"`
	HeadCommit *PushEventCommit     `json:"head_commit"`
	Repo       *PushEventRepository `json:"repository"`
	Compare    string               `json:"compare"`
}

type EngineDescribeRunsParams struct {
	ProjectId string `json:"project_id"`
	RunId     string `json:"run_id"`
	UserId    string `json:"user_id"`
	EventId   string `json:"event_id"`
	Offset    int    `json:"offset"`
	Limit     int    `json:"limit"`
}

type ProjectStoreCreateRunResponse struct {
	BuildId string `json:"build_id"`
}
