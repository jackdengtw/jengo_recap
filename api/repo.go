package api

type RepoMeta struct {
	OriginId string   `json: origin_id`
	OwnerIds []string `json:"owner_ids"` // userId in Jengo
	ScmName  string   `json:"scm_name"`

	Name     *string `json:"name,omitempty"`
	FullName *string `json:"full_name,omitempty"`
	Url      *string `json:"url,omitempty"`
	HtmlUrl  *string `json:"html_url,omitempty"`

	HooksUrl *string `json:"hooks_url,omitempty"`
}

type Repo struct {
	RepoMeta `json:"meta"`

	Id      string `json:"id" bson:"_id"` // repoId in Jengo
	Enabled bool   `json:"enabled"`

	UserIds       []string `json:"user_ids" bson:"user_ids"` // userId in Jengo
	Branches      []string `json:"branches"`
	State         string   `json:"state"`
	LatestBuildId string

	// TODO: move this out of API obj
	BuildIndex int
}
