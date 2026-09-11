package appaction

import "github.com/randomvlad/trader-vlads/internal/util/stringutil"

// rename package to press
// struct to KeyPressAction ?

type AppAction struct {
	name       string
	keyPress   string
	actionFunc func()
}

func NewAppAction(name string, actionFunc func()) AppAction {
	return AppAction{
		name:       name,
		keyPress:   stringutil.FirstCharLowercase(name),
		actionFunc: actionFunc,
	}
}

// TODO: temp while coding is being migrated and refactored
func NewAppActionNoOp(name string) AppAction {
	return NewAppAction(name, func() {})
}

func (a AppAction) GetName() string {
	return a.name
}

func (a AppAction) GetKeyPress() string {
	return a.keyPress
}

func (a AppAction) GetActionFunc() func() {
	return a.actionFunc
}
