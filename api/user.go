package api

import "time"

type AuthBase struct {
	Id          string `json:"id"`        // github_$originId
	OriginId    string `json:"origin_id"` // origin user id for auth
	Primary     bool   `json:"primary"`
	DisplayName string `json:"display_name"`
	LoginName   string `json:"login_name"`
	Email       string `json:"email"`
	AvatarUrl   string `json:"avartar_url"`
	Locale      string `json:"locale"`
}

type Auth struct {
	AuthBase
	AuthSource AuthSource `json:"auth_source"`
}

type ScmBase struct {
	Id          string    `json:"id"`        // github_$originId
	OriginId    string    `json:"origin_id"` // origin user id for scm
	DisplayName string    `json:"display_name"`
	LoginName   string    `json:"login_name"`
	Email       string    `json:"email"`
	AvatarUrl   string    `json:"avartar_url"`
	Locale      string    `json:"locale"`
	SyncedAt    time.Time `json:"synced_at"`
	BoundAt     time.Time `json:"bound_at"`
}

type Scm struct {
	ScmBase
	AuthSource AuthSource `json:"auth_source"`
	ScmSource  ScmSource  `json:"scm_source"`
}

type User02 struct {
	Id        string     `json:"id"` // userId in Jengo
	UpdatedAt *time.Time `json:"updated_at"`
	CreatedAt *time.Time `json:"created_at"`
	Auth      Auth       `json:"auth"`

	Scms []Scm `json:"scms,omitempty"`
}

type AuthSource struct {
	Name     string `json:"id"`
	OauthUrl string `json:"oauth_url"`
	TokenUrl string `json:"token_url"`
}

type ScmSource struct {
	AuthSource
	ApiUrl string `json:"api_url"`
}

// Facts
// Alternatively they can be to Mongo table or config file

const (
	AUTH_SOURCE_GITHUB = "github"
)

var (
	AUTH_SOURCES map[string]AuthSource
	SCM_SOURCES  map[string]ScmSource
)

func init() {
	github := AuthSource{
		Name:     AUTH_SOURCE_GITHUB,
		OauthUrl: "http://github.com/login/oauth/authorize",
		TokenUrl: "https://github.com/login/oauth/access_token",
	}
	AUTH_SOURCES = map[string]AuthSource{
		github.Name: github,
	}

	SCM_SOURCES = map[string]ScmSource{
		AUTH_SOURCE_GITHUB: ScmSource{
			AuthSource: github,
			ApiUrl:     "https://api.gitbub.com",
		},
	}
}
