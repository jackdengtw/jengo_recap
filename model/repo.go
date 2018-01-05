package model

import (
	"github.com/qetuantuan/jengo_recap/api"
)

type Repo api.Repo

func (r *Repo) ToApiObj() *api.Repo {
	return (*api.Repo)(r)
}

// Note: Shadow Copy
func NewRepoFrom(r *api.Repo) *Repo {
	return (*Repo)(r)
}

type ById []Repo

func (s ById) Len() int {
	return len(s)
}
func (s ById) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ById) Less(i, j int) bool {
	return s[i].Id < s[j].Id
}
