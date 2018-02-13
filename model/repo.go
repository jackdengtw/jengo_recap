package model

import (
	"github.com/qetuantuan/jengo_recap/vo"
)

type RepoMeta struct {
	OriginId string `json: origin_id`
	ScmName  string `json:"scm_name"`

	Name     *string `json:"name,omitempty"`
	FullName *string `json:"full_name,omitempty"`
	GitUrl   *string `json:"git_url,omitempty"`
	HtmlUrl  *string `json:"html_url,omitempty"`

	HooksUrl *string `json:"hooks_url,omitempty"`
}

type Repo struct {
	RepoMeta `json:"meta"`

	Id      string `json:"id" bson:"_id"` // repoId in Jengo
	Enabled bool   `json:"enabled"`

	OwnerIds []string `json:"owner_ids" bson:"owner_ids"` // userId in Jengo
	UserIds  []string `json:"user_ids" bson:"user_ids"`   // userId in Jengo
	Branches []string `json:"branches"`

	BuildIndex int
}

func (r *Repo) ToViewObj() *vo.Repo {
	return &vo.Repo{
		RepoMeta: vo.RepoMeta{
			OriginId: r.OriginId,
			ScmName:  r.ScmName,

			Name:     r.Name,
			FullName: r.FullName,
			GitUrl:   r.GitUrl,
			HtmlUrl:  r.HtmlUrl,

			HooksUrl: r.HooksUrl,
		},
		Id:      r.Id,
		Enabled: r.Enabled,

		OwnerIds: r.OwnerIds,
		UserIds:  r.UserIds,
		Branches: r.Branches,

		State:         "",
		LatestBuildId: "",
	}
}

// Note: Shadow Copy
func NewRepoFrom(r *vo.Repo) *Repo {
	return &Repo{
		RepoMeta: RepoMeta{
			OriginId: r.OriginId,
			ScmName:  r.ScmName,

			Name:     r.Name,
			FullName: r.FullName,
			GitUrl:   r.GitUrl,
			HtmlUrl:  r.HtmlUrl,

			HooksUrl: r.HooksUrl,
		},
		Id:      r.Id,
		Enabled: r.Enabled,

		OwnerIds: r.OwnerIds,
		UserIds:  r.UserIds,
		Branches: r.Branches,
	}
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
