package registry

import (
	"github.com/qetuantuan/jengo_recap/action"
	"github.com/qetuantuan/jengo_recap/definition"
)

type TemplateRegistry struct {
	// a map from language to []*Template
	// a map from buildscript to []*Template
}

func (r *TemplateRegistry) Match(m definition.Manifest) *definition.Template {
	// select from language map
	// select from buildscript map
	return &DefaultTemplate
}

func (r *TemplateRegistry) Register(*definition.Template) {
	// add a template to registry
}

// TODO: Replace it with a templateService later on
var GlobalTemplateRegistry TemplateRegistry = TemplateRegistry{}

// temporarily with global var
// TODO: use yml parser to init templates after parse logic done
var DefaultTemplate definition.Template = definition.Template{
	Name:        "Default",
	Language:    "*",
	BuildScript: "*",
	Steps: []definition.Step{
		definition.Step{
			Name:        "Get Source",
			UserVisible: true,
			Actions: []action.Action{
				&action.Bash{
					Env:           []string{},
					Timeout:       300,
					ScriptContent: "git clone https://github.com/qetuantuan/project1",
				},
				&action.Bash{
					Env:           []string{},
					Timeout:       300,
					ScriptContent: "go test",
				},
			},
		},
	},
}
