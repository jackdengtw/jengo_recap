package model

type CreateRunParams struct {
	ProjectId  string `json:"project_id"`
	UserId     string `json:"user_id"`
	ScmUrl     string `json:"scm_url"`
	Branch     string `json:"branch"`
	CommitId   string `json:"commit_id"`
	CommitUser string `json:"commit_user"`
	CommitLink string `json:"commit_link"`
	DiffLink   string `json:"diff_link"`
}

type DescribeRunsParams struct {
	ProjectId string `json:"project_id"`
	RunId     string `json:"run_id"`
	UserId    string `json:"user_id"`
	Offset    int    `json:"offset"`
	Limit     int    `json:"limit"`
}

type GetRunLogParams struct {
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
	RunId  string `json:"run_id"`
}
