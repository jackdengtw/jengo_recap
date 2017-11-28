package action

import (
	"github.com/qetuantuan/jengo_recap/context"
)

// Noop embedded to real actions
// providing common behavior
type Noop struct {
	Id string

	ctx   context.RunnerContext
	state ActionState
}

func (a *Noop) Do() {
}

func (a *Noop) GetId() string {
	return a.Id
}

func (a *Noop) State() ActionState {
	return a.state
}

func (a *Noop) SetContext(c context.RunnerContext) {
	a.ctx = c
}
