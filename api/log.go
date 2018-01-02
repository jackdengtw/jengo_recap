package api

type BuildLog struct {
	Id string `json:"id" bson:"_id"`

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
