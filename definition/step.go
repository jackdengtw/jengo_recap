package definition

import "github.com/qetuantuan/jengo_recap/action"

type Step struct {
	Actions     []action.Action
	UserVisible bool
	Name        string
}
