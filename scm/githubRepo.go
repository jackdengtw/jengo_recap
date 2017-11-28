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

type GithubProject struct {
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

func (p *GithubProject) CopyMetaTo(project *model.Project) {
	project.Meta.Id = "p_github_" + strconv.Itoa(p.Id)
	project.Meta.Name = p.Name
	project.Meta.Scm = "github"
	project.Meta.OwnerId = "u_github_" + strconv.Itoa(p.Owner.Id)
	project.Meta.FullName = p.FullName
	project.Meta.HtmlUrl = p.HtmlUrl
	project.Meta.CreatedAt = p.CreatedAt
	project.Meta.UpdatedAt = p.UpdatedAt
	project.Meta.PushedAt = p.PushedAt
	project.Meta.GitUrl = p.GitUrl
	project.Meta.SshUrl = p.SshUrl
	project.Meta.CloneUrl = p.CloneUrl
	project.Meta.Private = p.Private
	project.Meta.Language = p.Language
	project.Meta.HooksUrl = p.HooksUrl
}

func (p *GithubProject) CopyTo(project *model.Project) {
	p.CopyMetaTo(project)
}
