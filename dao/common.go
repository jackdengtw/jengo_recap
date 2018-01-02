package dao

import (
	"errors"

	"github.com/golang/glog"
)

const (
	repoDbName string = "repos"

	repoCol          string = "repo"
	semanticBuildCol string = "semanticBuild"
	logCol           string = "log"
	hookCol          string = "hook"
)

// User DB
const (
	userDbName string = "users"
	userCol02  string = "user02"
)

const (
	SCM_COLUMN  = "scms"
	AUTH_COLUMN = "auths"
)

// Engine DB
const (
	engineDbName string = "engine"
	buildCol     string = "build"
)

// Errors
var (
	ErrorBuildNotFind = errors.New("build not found")
	ErrorTypeNotMatch = errors.New("type not match")
)

// MgoLog
type MgoLog int

func (m MgoLog) Output(calldepth int, s string) error {
	// if glog.V(glog.Level(calldepth)) {
	// TODO: Not sure why user service don't have glog info file log
	//       use error instead for now
	glog.Error(s)
	// }
	return nil
}
