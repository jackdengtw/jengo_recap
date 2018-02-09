package model

import (
	"time"

	"github.com/qetuantuan/jengo_recap/vo"
)

// Build
// type Build vo.Build
type Build struct {
	Id     string `json:"id" bson:"_id"`
	Status string `json:"status"`
	Index  int
	Result string `json:"result"`

	// duplicate info in SemanticBuil for now. Later on remove Semantic Build
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

func (b *Build) ToApiObj() *vo.Build {
	return &vo.Build{
		Id:         b.Id,
		Status:     b.Status,
		Result:     b.Result,
		RepoId:     b.RepoId,
		CommitId:   b.CommitId,
		Branch:     b.Branch,
		UserId:     b.UserId,
		EventId:    b.EventId,
		Commits:    b.Commits.ToApiObj(),
		HeadCommit: b.HeadCommit.ToApiObj(),
		LogUri:     b.LogUri,
		UpdatedAt:  b.UpdatedAt,
		CreatedAt:  b.CreatedAt,
		StartTime:  b.StartTime,
		Duration:   b.Duration,
	}
}

// Note: Shadow Copy
func NewBuildFrom(b *vo.Build) *Build {
	return &Build{
		Id:         b.Id,
		Status:     b.Status,
		Result:     b.Result,
		RepoId:     b.RepoId,
		CommitId:   b.CommitId,
		Branch:     b.Branch,
		UserId:     b.UserId,
		EventId:    b.EventId,
		Commits:    NewPushEventCommits(b.Commits),
		HeadCommit: NewPushEventCommit(*b.HeadCommit),
		LogUri:     b.LogUri,
		UpdatedAt:  b.UpdatedAt,
		CreatedAt:  b.CreatedAt,
		StartTime:  b.StartTime,
		Duration:   b.Duration,
	}
}

// Builds
type Builds []Build

func (bs Builds) ToApiObj() (builds vo.Builds) {
	for _, b := range bs {
		builds = append(builds, *b.ToApiObj())
	}
	return
}

// Note: Shadow Copy
func NewBuildsFrom(bs vo.Builds) (builds Builds) {
	for _, b := range bs {
		builds = append(builds, *NewBuildFrom(&b))
	}
	return
}

// SemanticBuild
// a copy of vo.SemanticBuild
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

func (b *SemanticBuild) ToApiObj() *vo.SemanticBuild {
	voObj := vo.SemanticBuild{
		Id:       b.Id,
		RepoId:   b.RepoId,
		CommitId: b.CommitId,
		Branch:   b.Branch,
		UserId:   b.UserId,
		Numero:   b.Numero,
		Builds:   b.Builds.ToApiObj(),
	}
	return &voObj
}

// Note: Shadow Copy
func NewSemanticBuildFrom(b *vo.SemanticBuild) *SemanticBuild {
	sbuild := SemanticBuild{
		Id:       b.Id,
		RepoId:   b.RepoId,
		CommitId: b.CommitId,
		Branch:   b.Branch,
		UserId:   b.UserId,
		Numero:   b.Numero,
		Builds:   NewBuildsFrom(b.Builds),
	}
	return &sbuild
}

// Semantic
type SemanticBuilds []SemanticBuild

func (bs SemanticBuilds) ToApiObj() (builds vo.SemanticBuilds) {
	for _, b := range bs {
		builds = append(builds, *b.ToApiObj())
	}
	return
}

// Note: Shadow Copy
func NewSemanticBuildsFrom(bs vo.SemanticBuilds) (builds SemanticBuilds) {
	for _, b := range bs {
		builds = append(builds, *NewSemanticBuildFrom(&b))
	}
	return
}
