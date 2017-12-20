package api

type RepoMeta struct {
	Id       string   `json:"id"` // repoId in Jengo
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
	Enable   bool `json:"enable"`

	UserIds  []string `json:"user_ids"` // userId in Jengo
	Branches []string `json:"branches"`
}
