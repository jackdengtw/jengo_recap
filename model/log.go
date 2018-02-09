package model

import (
	"github.com/qetuantuan/jengo_recap/vo"
)

type BuildLog vo.BuildLog

func (r *BuildLog) ToApiObj() *vo.BuildLog {
	return (*vo.BuildLog)(r)
}

// Note: Shadow Copy
func NewBuildLogFrom(r *vo.BuildLog) *BuildLog {
	return (*BuildLog)(r)
}
