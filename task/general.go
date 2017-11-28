package task

import (
	"github.com/qetuantuan/jengo_recap/api"
	"github.com/qetuantuan/jengo_recap/definition"
	"github.com/qetuantuan/jengo_recap/model"
)

type General struct {
	Version uint

	// Start // executor update Run.Run.State
	Run *model.InnerRun

	// Parser
	User    *api.User02
	Project *model.Project

	Template *definition.Template
	Manifest *definition.Manifest

	// Finalize
	Status string
}
