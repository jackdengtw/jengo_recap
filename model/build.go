package model

import (
	"github.com/qetuantuan/jengo_recap/api"
)

// Build
type Build api.Build

func (r *Build) ToApiObj() *api.Build {
	return (*api.Build)(r)
}

// Note: Shadow Copy
func NewBuildFrom(r *api.Build) *Build {
	return (*Build)(r)
}

// SemanticBuild
type SemanticBuild api.SemanticBuild

func (b *SemanticBuild) ToApiObj() *api.SemanticBuild {
	return (*api.SemanticBuild)(b)
}

// Note: Shadow Copy
func NewSemanticBuildFrom(b *api.SemanticBuild) *SemanticBuild {
	return (*SemanticBuild)(b)
}

// Semantic
type SemanticBuilds []SemanticBuild

func (bs *SemanticBuilds) ToApiObj() (builds api.SemanticBuilds) {
	for _, b := range *bs {
		builds = append(builds, *b.ToApiObj())
	}
	return
}

// Note: Shadow Copy
func NewSemanticBuildsFrom(bs api.SemanticBuilds) (builds SemanticBuilds) {
	for _, b := range bs {
		builds = append(builds, *NewSemanticBuildFrom(&b))
	}
	return
}
