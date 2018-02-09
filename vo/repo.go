// TODO: rename to ViewObject
package vo

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

	OwnerIds []string `json:"owner_ids"`                // userId in Jengo
	UserIds  []string `json:"user_ids" bson:"user_ids"` // userId in Jengo
	Branches []string `json:"branches"`

	// from Build
	State         string `json:"state"`
	LatestBuildId string
}
