package dao

import (
	"errors"
	"fmt"

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
type BatchError struct {
	FailedIdx int
	RealErr   error
}

func (e BatchError) Error() string {
	return fmt.Sprintf("No. %v element failed. Real err: %v", e.FailedIdx, e.RealErr)
}

var (
	ErrorBuildNotFound           = errors.New("build not found")
	ErrorTypeNotMatch            = errors.New("type not match")
	ErrorAlreadyExisted          = errors.New("the same build id already existed")
	ErrorMoreThanOneBuildExisted = errors.New("more than one build existed for one id")
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
