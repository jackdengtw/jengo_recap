package model

import (
	"strconv"
	"time"

	"github.com/golang/glog"
	"github.com/qetuantuan/jengo_recap/api"
	"github.com/qetuantuan/jengo_recap/util"
)

type ApiUser api.User02

/*
 */

type GithubUser struct {
	Login      string `json:"login"`
	Id         int    `json:"id"`
	AvatarUrl  string `json:"avatar_url"`
	GravatarId string `json:"gravatar_id"`
	Name       string `json:"name"`
	Location   string `json:"location"`
	Email      string `json:"email"`
	/*Url                       string `json:"url"`
	HtmlUrl                  string `json:"html_url"`
	GistsUrl                 string `json:"gists_url"`
	OrganizationsUrl         string `json:"organizations_url"`
	ReposUrl                 string `json:"repos_url"`
	EventsUrl                string `json:"events_url"`
	Type                      string `json:"type"`
	SiteAdmin                bool `json:"site_admin"`
	Company                   string `json:"company"`
	Blog                      string `json:"blog"`

	Hireable                  string `json:"hireable"`
	Bio                       string `json:"bio"`
	PublicRepos              int `json:"private_repos"`
	PublicGists              int `json:"public_gists"`
	Followers                 int `json:"followers"`
	Following                 int    `json:"following"`
	CreatedAt                time.Time `json:"created_at"`
	UpdatedAt                time.Time `json:"updated_at"`
	TotalPrivateRepos       int `json:"total_private_repos"`
	OwnedPrivateRepos       int `json:"owned_private_repos"`
	PrivateGists             int `json:"private_gists"`
	DiskUsage                int `json:"disk_usage"`
	Collaborators             string `json:"collaborators"`
	TwoFactorAuthentication bool `json:"two_factor_authentication"`*/
}

func (u *GithubUser) ToMongoUser() (mongoUser *User) {
	mongoUser = &User{}
	// TODO: A new way to generate user id
	timeNow := time.Now().UTC()
	mongoUser.UserId = "u_github_" + strconv.Itoa(u.Id)
	mongoUser.CreatedAt = timeNow
	mongoUser.UpdatedAt = timeNow
	mongoUser.Auths = []Auth{
		Auth{
			Id:           "u_github_" + strconv.Itoa(u.Id),
			OriginId:     u.Id,
			AuthSourceId: api.AUTH_SOURCE_GITHUB,
			// Ignore token here
			Primary:     true,
			DisplayName: u.Name,
			LoginName:   u.Login,
			Email:       u.Email,
			AvatarUrl:   u.AvatarUrl,
			Locale:      u.Location,
		},
	}
	mongoUser.Scms = []Scm{
		Scm{
			Id:           mongoUser.UserId,
			OriginId:     u.Id,
			AuthSourceId: api.AUTH_SOURCE_GITHUB,
			ScmSourceId:  api.AUTH_SOURCE_GITHUB,
			// Ignore token here
			DisplayName: u.Name,
			LoginName:   u.Login,
			Email:       u.Email,
			AvatarUrl:   u.AvatarUrl,
			Locale:      u.Location,
			SyncedAt:    timeNow, //Todo: update synced time after update repo/project info from scm
			BoundAt:     timeNow,
		},
	}
	return mongoUser
}

// Mongo Object
type Auth struct {
	Id           string      `bson:"id"`
	OriginId     interface{} `bson:"origin_id"`
	AuthSourceId string      `bson:"auth_source_id"` // refer to AuthSource.Id
	Token        []byte      `bson:"token"`
	Primary      bool        `bson:"primary"`
	DisplayName  string      `bson:"display_name"`
	LoginName    string      `bson:"login_name"`
	Email        string      `bson:"email"`
	AvatarUrl    string      `bson:"avartar_url"`
	Locale       string      `bson:"locale"`
}

func (a *Auth) ToApiAuth() *api.Auth {
	auth := &api.Auth{
		Id:       a.Id,
		OriginId: a.OriginId,
		// Ignore token,
		Primary:     a.Primary,
		DisplayName: a.DisplayName,
		LoginName:   a.LoginName,
		Email:       a.Email,
		AvatarUrl:   a.AvatarUrl,
		Locale:      a.Locale,
		AuthSource:  api.AUTH_SOURCES[a.AuthSourceId],
	}
	return auth
}

// Mongo Object
type Scm struct {
	Id           string      `bson:"id"`
	OriginId     interface{} `bson:"origin_id"`
	AuthSourceId string      `bson:"auth_source_id"` // refer to AuthSource.Id
	ScmSourceId  string      `bson:"scm_source_id"`  // refer to ScmSource.Id
	Token        []byte      `bson:"token"`
	DisplayName  string      `bson:"display_name"`
	LoginName    string      `bson:"login_name"`
	Email        string      `bson:"email"`
	AvatarUrl    string      `bson:"avartar_url"`
	Locale       string      `bson:"locale"`
	SyncedAt     time.Time   `bson:"synced_at"`
	BoundAt      time.Time   `bson:"bound_at"`
}

func (s *Scm) ToApiScm() *api.Scm {
	scm := &api.Scm{
		Id:       s.Id,
		OriginId: s.OriginId,
		// Ignore Token,
		DisplayName: s.DisplayName,
		LoginName:   s.LoginName,
		Email:       s.Email,
		AvatarUrl:   s.AvatarUrl,
		Locale:      s.Locale,
		SyncedAt:    s.SyncedAt,
		BoundAt:     s.BoundAt,
		AuthSource:  api.AUTH_SOURCES[s.AuthSourceId],
		ScmSource:   api.SCM_SOURCES[s.ScmSourceId],
	}
	return scm
}

// Mongo Object User v0.2
type User struct {
	UserId    string    `bson:"_id"` // Temporarily
	UpdatedAt time.Time `bson:"updated_at"`
	CreatedAt time.Time `bson:"created_at"`
	Auths     []Auth    `bson:"auths"`
	Scms      []Scm     `bson:"scms"`
}

func (u *User) SetTokenEncrypted(id string, key string, token string) (err error) {
	var t []byte
	t, err = util.AESEncode([]byte(key), []byte(token))
	if err != nil {
		glog.Errorf("Encoding token error: %v", err)
		return
	}

	for i, _ := range u.Auths {
		if u.Auths[i].AuthSourceId == id {
			u.Auths[i].Token = t
		}
	}
	for i, _ := range u.Scms {
		if u.Scms[i].AuthSourceId == id {
			u.Scms[i].Token = t
		}
	}
	return
}

func (u *User) ToApiUser() (*ApiUser, error) {
	user := &ApiUser{
		UserId:    u.UserId,
		UpdatedAt: u.UpdatedAt,
		CreatedAt: u.CreatedAt,
	}
	for i, _ := range u.Auths {
		if plainText, err := util.AESDecode([]byte(util.KeyCoder), u.Auths[i].Token); err != nil {
			glog.Errorf("De/Encoding token error: %v", err)
			return nil, err
		} else {
			if u.Auths[i].Primary {
				user.Auth = *u.Auths[i].ToApiAuth()
				user.Auth.Token = string(plainText)
				break
			}
		}
	}
	for i, _ := range u.Scms {
		if plainText, err := util.AESDecode([]byte(util.KeyCoder), u.Scms[i].Token); err != nil {
			glog.Errorf("De/Encoding token error: %v", err)
			return nil, err
		} else {
			user.Scms = append(user.Scms, *u.Scms[i].ToApiScm())
			// same index as u.Scms
			user.Scms[i].Token = string(plainText)
		}
	}
	return user, nil
}
