package model

import (
	"github.com/qetuantuan/jengo_recap/api"
)

type BuildLog api.BuildLog

func (r *BuildLog) ToApiObj() *api.BuildLog {
	return (*api.BuildLog)(r)
}

// Note: Shadow Copy
func NewBuildLogFrom(r *api.BuildLog) *BuildLog {
	return (*BuildLog)(r)
}
