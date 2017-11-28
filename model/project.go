package model

import (
	"github.com/qetuantuan/jengo_recap/api"
)

// Project
type Project struct {
	api.Project
	HashId []byte `bson:"_id"`
}

func (p *Project) ToApiObj() *api.Project {
	return &p.Project
}

func NewProjectFrom(p *api.Project) *Project {
	return &Project{
		Project: *p,
	}
}

type ById []Project

func (s ById) Len() int {
	return len(s)
}
func (s ById) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ById) Less(i, j int) bool {
	return s[i].Meta.Id < s[j].Meta.Id
}

// Run
type Run api.Run

func (r *Run) ToApiObj() *api.Run {
	return (*api.Run)(r)
}

func NewRunFrom(r *api.Run) *Run {
	return (*Run)(r)
}

// PatchRun
type PatchRun api.PatchRun

func (r *PatchRun) ToApiObj() *api.PatchRun {
	return (*api.PatchRun)(r)
}

func NewPatchRunFrom(r *api.PatchRun) *PatchRun {
	return (*PatchRun)(r)
}

// Build
type Build api.Build

func (b *Build) ToApiObj() *api.Build {
	return (*api.Build)(b)
}

func NewBuildFrom(b *api.Build) *Build {
	return (*Build)(b)
}

type Builds []Build

func (bs *Builds) ToApiObj() (builds api.Builds) {
	for _, b := range *bs {
		builds = append(builds, *b.ToApiObj())
	}
	return
}

func NewBuildsFrom(bs api.Builds) (builds Builds) {
	for _, b := range bs {
		builds = append(builds, *NewBuildFrom(&b))
	}
	return
}

// RunLog
type RunLog api.RunLog

func (r *RunLog) ToApiObj() *api.RunLog {
	return (*api.RunLog)(r)
}

func NewRunLogFrom(r *api.RunLog) *RunLog {
	return (*RunLog)(r)
}
