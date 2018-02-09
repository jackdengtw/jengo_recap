package scm

import (
	"strconv"
	"time"

	"github.com/qetuantuan/jengo_recap/model"
)

type User struct {
	Login string
	Id    int
}

type GithubRepo struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	FullName  string    `json:"full_name"`
	Owner     User      `json:"owner"`
	HtmlUrl   string    `json:"html_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	PushedAt  time.Time `json:"pushed_at"`
	HooksUrl  string    `json:"hooks_url"`
	GitUrl    string    `json:"git_url"`
	SshUrl    string    `json:"ssh_url"`
	CloneUrl  string    `json:"clone_url"`
	Private   bool      `json:"private"`
	Language  string    `json:"language"`
}

func (p *GithubRepo) CopyRepoMetaTo(repo *model.Repo) {
	repo.Id = "p_github_" + strconv.Itoa(p.Id)
	repo.RepoMeta.Name = &p.Name
	repo.ScmName = "github"
	repo.OwnerIds = []string{"u_github_" + strconv.Itoa(p.Owner.Id)}
	repo.RepoMeta.FullName = &p.FullName
	repo.RepoMeta.HtmlUrl = &p.HtmlUrl
	// repo.RepoMeta.CreatedAt = p.CreatedAt
	// repo.RepoMeta.UpdatedAt = p.UpdatedAt
	// repo.RepoMeta.PushedAt = p.PushedAt
	repo.RepoMeta.GitUrl = &p.GitUrl
	repo.RepoMeta.HtmlUrl = &p.HtmlUrl
	// repo.RepoMeta.SshUrl = p.SshUrl
	// repo.RepoMeta.CloneUrl = p.CloneUrl
	// repo.RepoMeta.Private = p.Private
	// repo.RepoMeta.Language = p.Language
	repo.RepoMeta.HooksUrl = &p.HooksUrl
}

func (p *GithubRepo) CopyTo(Repo *model.Repo) {
	p.CopyRepoMetaTo(Repo)
}
