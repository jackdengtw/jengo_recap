package api

type RunLog struct {
	Id       string `bson:"_id"`
	RunId    string `json:"run_id"`
	FileName string `json:"file_name"`
	Content  string `bson:"content"`
	Length   int    `json:"length"`
}

type GetRunLogParams struct {
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
	RunId  string `json:"run_id"`
}
