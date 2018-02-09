package vo

// TODO: move log struct out of API object.
// Log should be bytes and as payload of http response.
type BuildLog struct {
	Id string `json:"id" bson:"_id"`

	Content []byte `bson:"content"`

	// TODO:
	// use an object storage
	// url string
	// time stamp
}

type GetBuildLogParams struct {
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
	Uri    string `json:"uri"`
}
