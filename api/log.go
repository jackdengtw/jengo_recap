package api

// TODO: move log struct out of API object.
// Log should be bytes and as payload of http response.
type BuildLog struct {
	Id string `json:"id" bson:"_id"`

	FileName string `json:"file_name"`
	Content  string `bson:"content"`

	// TODO:
	// url string
	// time stampe
}

type GetBuildLogParams struct {
	Offset  int    `json:"offset"`
	Limit   int    `json:"limit"`
	BuildId string `json:"build_id"`
}
