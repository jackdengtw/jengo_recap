package model

type Hook interface {
	GetProjectId() string
	SetProjectId(string)
	GetUrl() string
}

// Github Hook
//Todo : one project may have multiple hooks related to our system, but now we restrict only one
type GithubHook struct {
	Id        string         `json:"id" bson:"_id"`
	ProjectId string         `json:"project_id"`
	OriginId  int            `json:"origin_id"`
	Events    []string       `json:"events"`
	Config    GithubHookConf `json:"config"`
	Type      string         `json:"type"`
	Name      string         `json:"name"`
	Active    bool           `json:"active"`
	Url       string         `json:"url"`
	TestUrl   string         `json:"test_url"`
	PingUrl   string         `json:"ping_url"`
}

// compile time check of implementation
var _ Hook = &GithubHook{}

func (gh *GithubHook) GetProjectId() string {
	return gh.ProjectId
}
func (gh *GithubHook) SetProjectId(projectId string) {
	gh.ProjectId = projectId
}
func (gh *GithubHook) GetUrl() string {
	return gh.Url
}

// GithubHook Config
type GithubHookConf struct {
	ContentType string `json:"content_type"`
	Url         string `json:"url"`
	InsecureSsl string `json:"insecure_ssl"`
}
