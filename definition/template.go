package definition

type Template struct {
	Steps       []Step
	Env         map[string]string
	Name        string
	Language    string
	BuildScript string
}
