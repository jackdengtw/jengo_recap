package model

import (
	"github.com/qetuantuan/jengo_recap/api"
)

// Build
type Build api.Build

func (b *Build) ToApiObj() *api.Build {
	return (*api.Build)(b)
}

// Note: Shadow Copy
func NewBuildFrom(b *api.Build) *Build {
	return (*Build)(b)
}

// Builds
type Builds []Build

func (bs Builds) ToApiObj() (builds api.Builds) {
	for _, b := range bs {
		builds = append(builds, *b.ToApiObj())
	}
	return
}

// Note: Shadow Copy
func NewBuildsFrom(bs api.Builds) (builds Builds) {
	for _, b := range bs {
		builds = append(builds, *NewBuildFrom(&b))
	}
	return
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

func (bs SemanticBuilds) ToApiObj() (builds api.SemanticBuilds) {
	for _, b := range bs {
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
