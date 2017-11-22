package api

import "time"

// API Object
type Auth struct {
	Id          string      `json:"id"`
	OriginId    interface{} `json:"origin_id"`
	AuthSource  AuthSource  `json:"auth_source"`
	Token       string      `json:"token,omitempty"`
	Primary     bool        `json:"primary"`
	DisplayName string      `json:"display_name"`
	LoginName   string      `json:"login_name"`
	Email       string      `json:"email"`
	AvatarUrl   string      `json:"avartar_url"`
	Locale      string      `json:"locale"`
}

// API Object
type Scm struct {
	Id          string      `json:"id"`
	OriginId    interface{} `json:"origin_id"`
	AuthSource  AuthSource  `json:"auth_source"`
	ScmSource   ScmSource   `json:"scm_source"`
	Token       string      `json:"token,omitempty"`
	DisplayName string      `json:"display_name"`
	LoginName   string      `json:"login_name"`
	Email       string      `json:"email"`
	AvatarUrl   string      `json:"avartar_url"`
	Locale      string      `json:"locale"`
	SyncedAt    time.Time   `json:"synced_at"`
	BoundAt     time.Time   `json:"bound_at"`
}

// API Object User v0.2
type User02 struct {
	UserId    string    `json:"user_id"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
	Auth      Auth      `json:"auth"`
	Scms      []Scm     `json:"scms,omitempty"`
}

// API Object
type AuthSource struct {
	Id       string `json:"id"`
	OauthUrl string `json:"oauth_url"`
	TokenUrl string `json:"token_url"`
}

// API Object
type ScmSource struct {
	AuthSource
	ApiUrl string `json:"api_url"`
}

// Facts
// Alternatively they can be to Mongo table or config file

const (
	AUTH_SOURCE_GITHUB = "github.com"
)

var AUTH_SOURCES map[string]AuthSource
var SCM_SOURCES map[string]ScmSource

func init() {
	github := AuthSource{
		Id:       AUTH_SOURCE_GITHUB,
		OauthUrl: "http://github.com/login/oauth/authorize",
		TokenUrl: "https://github.com/login/oauth/access_token",
	}
	AUTH_SOURCES = map[string]AuthSource{
		github.Id: github,
	}

	SCM_SOURCES = map[string]ScmSource{
		AUTH_SOURCE_GITHUB: ScmSource{
			AuthSource: github,
			ApiUrl:     "https://api.gitbub.com",
		},
	}
}
