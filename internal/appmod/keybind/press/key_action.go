package press

import (
	tea "charm.land/bubbletea/v2"
	"github.com/randomvlad/trader-vlads/internal/util/stringutil"
)

type KeyAction struct {
	Name            string
	KeyPress        string
	ActionFunc      func() tea.Cmd
	Permanent       bool
	SortOrder       int
	IsFooterVisible func() bool
}

type KeyActionBuilder struct {
	name              string
	keyPress          string
	actionFunc        func() tea.Cmd
	permanent         bool
	sortOrder         int
	footerVisibleFunc func() bool
}

func NewActionBuilder() *KeyActionBuilder {
	return &KeyActionBuilder{}
}

func (b *KeyActionBuilder) Name(name string) *KeyActionBuilder {
	return b.NameKeyPress(name, stringutil.FirstCharLowercase(name))
}

func (b *KeyActionBuilder) NameKeyPress(name string, keyPress string) *KeyActionBuilder {
	b.name = name
	b.keyPress = keyPress
	return b
}

func (b *KeyActionBuilder) Permanent() *KeyActionBuilder {
	b.permanent = true
	return b
}

func (b *KeyActionBuilder) FooterHidden() *KeyActionBuilder {
	return b.FooterVisible(func() bool {
		return false
	})
}

func (b *KeyActionBuilder) FooterVisible(visible func() bool) *KeyActionBuilder {
	b.footerVisibleFunc = visible
	return b
}

func (b *KeyActionBuilder) SortOrder(sortOrder int) *KeyActionBuilder {
	b.sortOrder = sortOrder
	return b
}

func (b *KeyActionBuilder) ActionCmd(funcCmd func() tea.Cmd) *KeyActionBuilder {
	b.actionFunc = funcCmd
	return b
}

func (b *KeyActionBuilder) Action(funcNoCmd func()) *KeyActionBuilder {
	b.actionFunc = func() tea.Cmd {
		funcNoCmd()
		return nil
	}
	return b
}

func (b *KeyActionBuilder) Build() KeyAction {
	if b.footerVisibleFunc == nil {
		b.footerVisibleFunc = func() bool {
			return true
		}
	}

	return KeyAction{
		Name:            b.name,
		KeyPress:        b.keyPress,
		ActionFunc:      b.actionFunc,
		Permanent:       b.permanent,
		SortOrder:       b.sortOrder,
		IsFooterVisible: b.footerVisibleFunc,
	}
}
