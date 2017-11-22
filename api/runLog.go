package api

type RunLog struct {
	RunId    string `json:"run_id"`
	FileName string `json:"file_name"`
	Log      string `json:"log"` //text?
	Length   int    `json:"length"`
}

type GetRunLogParams struct {
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
	RunId  string `json:"run_id"`
}
