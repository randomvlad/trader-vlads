package press

import "github.com/randomvlad/trader-vlads/internal/util/stringutil"

type KeyAction struct {
	Name       string
	KeyPress   string
	ActionFunc func()
}

func NewKeyAction(name string, actionFunc func()) KeyAction {
	return KeyAction{
		Name:       name,
		KeyPress:   stringutil.FirstCharLowercase(name),
		ActionFunc: actionFunc,
	}
}

// TODO: temp while coding is being migrated and refactored
func NewKeyActionNoOp(name string) KeyAction {
	return NewKeyAction(name, func() {})
}
