package runner

import (
	_ "github.com/qetuantuan/jengo_recap/context"
)

type dict map[string]string

// actionId: KeyVal Map
type actionResult map[string]dict

type GenericContext struct {
	Id string

	runner        *Generic
	actionResults []actionResult
}

func (r *GenericContext) Version() string {
	return "0.1"
}

func (r *GenericContext) Cwd() string {
	return r.runner.cwd
}

func (r *GenericContext) ContextId() string {
	return r.Id
}

func (r *GenericContext) LogDir() string {
	return r.runner.cwd + "/" + r.runner.logRelative
}

func (r *GenericContext) ActionResult(key string, asc bool) (value string) {
	return "TODO"
}

func (r *GenericContext) ActionResultAll(key string) (values []string) {
	return []string{"TODO"}
}

func (r *GenericContext) SetCurrentResult(map[string]string) error {
	return nil
}
