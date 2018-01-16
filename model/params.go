package model

import (
	"time"

	"github.com/qetuantuan/jengo_recap/api"
)

type CreateBuildParams struct {
	RepoId     string `json:"Repo_id"`
	UserId     string `json:"user_id"`
	ScmUrl     string `json:"scm_url"`
	Branch     string `json:"branch"`
	CommitId   string `json:"commit_id"`
	CommitUser string `json:"commit_user"`
	CommitLink string `json:"commit_link"`
	DiffLink   string `json:"diff_link"`
}

type CreateBuildResponse struct {
	BuildId string `json:"build_id"`
}

type ListBuildsParams struct {
	RepoId  string `json:"Repo_id"`
	BuildId string `json:"Build_id"`
	UserId  string `json:"user_id"`
	Offset  int    `json:"offset"`
	Limit   int    `json:"limit"`
}

type GetBuildLogParams struct {
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
	LogUri string `json:"loguri"`
}

type EngineCreateBuildParams struct {
	RepoId     string               `json:"Repo_id"`
	UserId     string               `json:"user_id"`
	EventId    string               `json:"event_id"`
	Branch     string               `json:"branch"`
	Commits    []PushEventCommit    `json:"commits"`
	HeadCommit *PushEventCommit     `json:"head_commit"`
	Repo       *PushEventRepository `json:"repository"`
	Compare    string               `json:"compare"`
}

type EngineListBuildsParams struct {
	RepoId  string `json:"Repo_id"`
	BuildId string `json:"Build_id"`
	UserId  string `json:"user_id"`
	EventId string `json:"event_id"`
	Offset  int    `json:"offset"`
	Limit   int    `json:"limit"`
}

type UpdateBuildParams struct {
	Id         string            `bson:"_id" json:"id"`
	EventId    string            `json:"event_id,omitempty"`
	UserId     string            `json:"user_id,omitempty"`
	RepoId     string            `json:"Repo_id,omitempty"`
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
	BuildEnv   string            `json:"Build_env,omitempty"`
	UpdatedAt  *time.Time        `json:"updated_at,omitempty"`
	CreatedAt  *time.Time        `json:"created_at,omitempty"`
}

type SemanticBuildIndicator struct {
	BuildId string
	UserId  string
	RepoId  string
}

type PatchBuildParams api.Build
