package model

import (
	"testing"

	"github.com/qetuantuan/jengo_recap/api"
	"github.com/qetuantuan/jengo_recap/util"
)

func TestSetToken(t *testing.T) {
	token := "iamatoken"
	user := &User{
		Auths: []Auth{
			Auth{
				AuthSourceId: "myid",
			},
		},
		Scms: []Scm{
			Scm{
				AuthSourceId: "myid",
			},
		},
	}
	user.SetTokenEncrypted("myid", util.KeyCoder, token)
	plainText, err := util.AESDecode([]byte(util.KeyCoder), user.Auths[0].Token)
	if err != nil || string(plainText) != token {
		t.Logf("plainText: %v, Error: %v", plainText, err)
		t.Error("Wrong in Auth.Token")
	}
	plainText, err = util.AESDecode([]byte(util.KeyCoder), user.Scms[0].Token)
	if err != nil || string(plainText) != token {
		t.Logf("plainText: %v, Error: %v", plainText, err)
		t.Error("Wrong in Scm.Token")
	}
}

func TestToApiUser(t *testing.T) {
	token := "iamatoken"
	user := &User{
		UserId: "123",
		Auths: []Auth{
			Auth{
				AuthSourceId: api.AUTH_SOURCE_GITHUB,
				Primary:      true,
			},
			Auth{
				AuthSourceId: api.AUTH_SOURCE_GITHUB,
			},
		},
		Scms: []Scm{
			Scm{
				AuthSourceId: api.AUTH_SOURCE_GITHUB,
			},
		},
	}
	user.SetTokenEncrypted(api.AUTH_SOURCE_GITHUB, util.KeyCoder, token)
	if u, err := user.ToApiUser(); err != nil {
		t.Fatalf("ToApiUser failed: %v", err)
	} else if u.Auth.Token != token {
		t.Errorf("token not decrypted. %v", u)
	} else if u.Auth.AuthSource.Id != api.AUTH_SOURCE_GITHUB {
		t.Error("auth source not initialized. %v", u)
	}
}
