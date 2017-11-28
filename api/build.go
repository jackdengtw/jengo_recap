package api

type Build struct {
	Id        string `json:"id" bson:"_id"`
	ProjectId string `json:"project_id"`
	UserId    string `json:"user_id"`
	Numero    int    `json:"numero"`
	Runs      []Run  `json:"runs"`
	CommitId  string `json:"commit_id"`
	Branch    string `json:"branch"`
}

type Builds []Build
