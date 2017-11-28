package action

import (
	"github.com/qetuantuan/jengo_recap/constant"
	"github.com/qetuantuan/jengo_recap/context"
)

type ActionState int

const (
	UNKNOWN ActionState = iota
	SUCCESS
	FAILED
)

var stateMap = map[ActionState]string{
	UNKNOWN: constant.RUN_STATE_UNKNOWN,
	SUCCESS: constant.RUN_STATE_SUCCESS,
	FAILED:  constant.RUN_STATE_FAILED,
}

func ActionStateToString(state ActionState) string {
	return stateMap[state]
}

type Action interface {
	Do()

	GetId() string

	State() ActionState

	SetContext(context.RunnerContext)
}
