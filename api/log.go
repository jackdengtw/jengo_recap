package api

type BuildLog struct {
	BuildId string `json:"build_id"`

	FileName string `json:"file_name"`
	Content  string `bson:"content"`

	// TODO:
	// url string
	// time stampe
}

type GetBuildLogParams struct {
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
	RunId  string `json:"run_id"`
}
