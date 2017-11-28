package api

import "time"

const (
	STATE_SUCCESS = "success"
	STATE_FAILED  = "failed"
	STATE_UNKNOW  = "unknown"
)

type BranchCommit struct {
	Sha string `json:"sha"`
	Url string `json:"url"`
}

type ProjectMeta struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	FullName  string    `json:"full_name"`
	OwnerId   string    `json:"owner_id"`
	HtmlUrl   string    `json:"html_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	PushedAt  time.Time `json:"pushed_at"`
	GitUrl    string    `json:"git_url"`
	SshUrl    string    `json:"ssh_url"`
	CloneUrl  string    `json:"clone_url"`
	HooksUrl  string    `json:"hooks_url"`
	Scm       string    `json:"scm"`
	Private   bool      `json:"private"`
	Language  string    `json:"language"`
}

type Project struct {
	Meta          ProjectMeta `json:"meta"`
	Enable        bool        `json:"enable"`
	State         string      `json:"state"`
	LatestBuildId string      `json:"latest_build_id"`
	RunIndex      int         `json:"-"`
	Users         []string    `json:"users"`
	Branches      []string    `json:"branches"`
}
