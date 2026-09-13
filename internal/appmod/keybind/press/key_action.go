package press

import (
	tea "charm.land/bubbletea/v2"
	"github.com/randomvlad/trader-vlads/internal/util/stringutil"
)

type KeyAction struct {
	Name       string
	KeyPress   string
	ActionFunc func() tea.Cmd
	Permanent  bool
}

func NewKeyAction(name string, actionFunc func()) KeyAction {
	return KeyAction{
		Name:     name,
		KeyPress: stringutil.FirstCharLowercase(name),
		ActionFunc: func() tea.Cmd {
			actionFunc()
			return nil
		},
		Permanent: false,
	}
}

func NewAction(name string) *KeyAction {
	return &KeyAction{
		Name:     name,
		KeyPress: stringutil.FirstCharLowercase(name),
	}
}

func NewKeyActionPermanent(name string, actionFunc func() tea.Cmd) KeyAction {
	return KeyAction{
		Name:       name,
		KeyPress:   stringutil.FirstCharLowercase(name),
		ActionFunc: actionFunc,
		Permanent:  true,
	}
}

// TODO: temp while coding is being migrated and refactored
func NewKeyActionNoOp(name string) KeyAction {
	return NewKeyAction(name, func() {})
}

func (k *KeyAction) WithKeyPress(keyPress string) *KeyAction {
	k.KeyPress = keyPress
	return k
}

func (k *KeyAction) WithPermanent(permanent bool) *KeyAction {
	k.Permanent = permanent
	return k
}

func (k *KeyAction) WithFuncCmd(funcCmd func() tea.Cmd) *KeyAction {
	k.ActionFunc = funcCmd
	return k
}

func (k *KeyAction) WithFunc(actionFunc func()) *KeyAction {
	k.ActionFunc = func() tea.Cmd {
		actionFunc()
		return nil
	}
	return k
}
