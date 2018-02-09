package scm

import (
	"strconv"
	"time"

	"github.com/qetuantuan/jengo_recap/vo"
	"github.com/qetuantuan/jengo_recap/model"
)

type GithubUser struct {
	Login      string `json:"login"`
	Id         int    `json:"id"`
	AvatarUrl  string `json:"avatar_url"`
	GravatarId string `json:"gravatar_id"`
	Name       string `json:"name"`
	Location   string `json:"location"`
	Email      string `json:"email"`
	/*
		Url              string `json:"url"`
		HtmlUrl          string `json:"html_url"`
		GistsUrl         string `json:"gists_url"`
		OrganizationsUrl string `json:"organizations_url"`
		ReposUrl         string `json:"repos_url"`
		EventsUrl        string `json:"events_url"`
		Type             string `json:"type"`
		SiteAdmin        bool   `json:"site_admin"`
		Company          string `json:"company"`
		Blog             string `json:"blog"`

		Hireable                string    `json:"hireable"`
		Bio                     string    `json:"bio"`
		PublicRepos             int       `json:"private_repos"`
		PublicGists             int       `json:"public_gists"`
		Followers               int       `json:"followers"`
		Following               int       `json:"following"`
		CreatedAt               time.Time `json:"created_at"`
		UpdatedAt               time.Time `json:"updated_at"`
		TotalPrivateRepos       int       `json:"total_private_repos"`
		OwnedPrivateRepos       int       `json:"owned_private_repos"`
		PrivateGists            int       `json:"private_gists"`
		DiskUsage               int       `json:"disk_usage"`
		Collaborators           string    `json:"collaborators"`
		TwoFactorAuthentication bool      `json:"two_factor_authentication"`
	*/
}

func (u *GithubUser) ToMongoUser() (mongoUser *model.User) {
	mongoUser = &model.User{}
	// TODO: A new way to generate user id
	timeNow := time.Now().UTC()
	mongoUser.Id = "u_github_" + strconv.Itoa(u.Id)
	mongoUser.CreatedAt = &timeNow
	mongoUser.UpdatedAt = &timeNow
	mongoUser.Auths = []model.Auth{
		model.Auth{
			AuthBase: vo.AuthBase{
				Id:          "u_github_" + strconv.Itoa(u.Id),
				OriginId:    strconv.Itoa(u.Id),
				Primary:     true,
				DisplayName: u.Name,
				LoginName:   u.Login,
				Email:       u.Email,
				AvatarUrl:   u.AvatarUrl,
				Locale:      u.Location,
			},
			AuthSourceId: vo.AUTH_SOURCE_GITHUB,
			// Ignore token here
		},
	}
	mongoUser.Scms = []model.Scm{
		model.Scm{
			ScmBase: vo.ScmBase{
				Id:          mongoUser.Id,
				OriginId:    strconv.Itoa(u.Id),
				DisplayName: u.Name,
				LoginName:   u.Login,
				Email:       u.Email,
				AvatarUrl:   u.AvatarUrl,
				Locale:      u.Location,
				SyncedAt:    timeNow, //Todo: update synced time after update repo/project info from scm
				BoundAt:     timeNow,
			},
			AuthSourceId: vo.AUTH_SOURCE_GITHUB,
			ScmSourceId:  vo.AUTH_SOURCE_GITHUB,
			// Ignore token here
		},
	}
	return mongoUser
}
