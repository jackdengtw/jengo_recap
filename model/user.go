package model

import (
	"time"

	"github.com/golang/glog"
	"github.com/qetuantuan/jengo_recap/api"
	"github.com/qetuantuan/jengo_recap/util"
)

type Auth struct {
	api.AuthBase
	AuthSourceId string `bson:"auth_source_id"` // refer to AuthSource.Id
	Token        []byte `bson:"token"`
}

func (a *Auth) ToApiAuth() *api.Auth {
	auth := &api.Auth{
		AuthBase: api.AuthBase{
			Id:          a.Id,
			OriginId:    a.OriginId,
			Primary:     a.Primary,
			DisplayName: a.DisplayName,
			LoginName:   a.LoginName,
			Email:       a.Email,
			AvatarUrl:   a.AvatarUrl,
			Locale:      a.Locale,
		},
		AuthSource: api.AUTH_SOURCES[a.AuthSourceId],
		// Ignore token,
	}
	return auth
}

// Mongo Object
type Scm struct {
	api.ScmBase
	AuthSourceId string `bson:"auth_source_id"` // refer to AuthSource.Id
	ScmSourceId  string `bson:"scm_source_id"`  // refer to ScmSource.Id
	Token        []byte `bson:"token"`
}

func (s *Scm) ToApiScm() *api.Scm {
	scm := &api.Scm{
		ScmBase: api.ScmBase{
			Id:          s.Id,
			OriginId:    s.OriginId,
			DisplayName: s.DisplayName,
			LoginName:   s.LoginName,
			Email:       s.Email,
			AvatarUrl:   s.AvatarUrl,
			Locale:      s.Locale,
			SyncedAt:    s.SyncedAt,
			BoundAt:     s.BoundAt,
		},
		AuthSource: api.AUTH_SOURCES[s.AuthSourceId],
		ScmSource:  api.SCM_SOURCES[s.ScmSourceId],
		// Ignore Token,
	}
	return scm
}

// Mongo Object User v0.2
type User struct {
	UserId    string     `bson:"_id"` // Temporarily
	UpdatedAt *time.Time `bson:"updated_at"`
	CreatedAt *time.Time `bson:"created_at"`
	Auths     []Auth     `bson:"auths"`
	Scms      []Scm      `bson:"scms"`
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

func (u *User) ToApiUser() (*api.User02, error) {
	user := &api.User02{
		UserId:    u.UserId,
		UpdatedAt: u.UpdatedAt,
		CreatedAt: u.CreatedAt,
	}
	for i, _ := range u.Auths {
		if u.Auths[i].Primary {
			user.Auth = *u.Auths[i].ToApiAuth()
			// user.Auth.Token = string(plainText)
			break
		}
	}
	for i, _ := range u.Scms {
		user.Scms = append(user.Scms, *u.Scms[i].ToApiScm())
		// user.Scms[i].Token = string(plainText)
	}
	return user, nil
}
