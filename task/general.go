package task

import (
	"github.com/qetuantuan/jengo_recap/definition"
	"github.com/qetuantuan/jengo_recap/model"
)

type General struct {
	Version uint

	// Start // executor update Run.Run.State
	Build *model.Build

	// Parser
	User *model.User
	Repo *model.Repo

	Template *definition.Template
	Manifest *definition.Manifest

	// Finalize
	Status string
}
