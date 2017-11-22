package util

import (
	"fmt"
)

type EntitySymbol string

type ScmSymbol string

const (
	UserSymbol    EntitySymbol = "u"
	ProjectSymbol EntitySymbol = "p"

	GithubSymbol ScmSymbol = "github"

	Separator = "_"
)

func GetUserId(scm ScmSymbol, rawId interface{}) string {
	return fmt.Sprintf("%v%v%v%v%v", UserSymbol, Separator, scm, Separator, rawId)
}

func GetProjectId(scm ScmSymbol, rawId interface{}) string {
	return fmt.Sprintf("%v%v%v%v%v", ProjectSymbol, Separator, scm, Separator, rawId)
}
