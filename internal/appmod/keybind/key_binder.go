package keybind

import (
	"cmp"
	"slices"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/randomvlad/trader-vlads/internal/appmod/keybind/press"
)

type KeyBinder struct {
	bindings map[KeyContext]press.KeyAction
}

func NewKeyBinder() *KeyBinder {
	return &KeyBinder{
		bindings: make(map[KeyContext]press.KeyAction),
	}
}

type KeyContext struct {
	ContextId string
	KeyPress  string
}

func NewKeyContext(contextId string, key string) KeyContext {
	return KeyContext{
		ContextId: contextId,
		KeyPress:  strings.ToLower(key),
	}
}

func (b *KeyBinder) Remove(k KeyContext) *KeyBinder {
	delete(b.bindings, k)
	return b
}

func (b *KeyBinder) AddAction(contextId string, action press.KeyAction) *KeyBinder {
	b.bindings[NewKeyContext(contextId, action.KeyPress)] = action
	return b
}

func (b *KeyBinder) AddActions(contextId string, actions ...press.KeyAction) *KeyBinder {
	for _, action := range actions {
		b.AddAction(contextId, action)
	}
	return b
}

func (b *KeyBinder) AddActionsMap(contextActions map[string][]press.KeyAction) *KeyBinder {
	for contextId, actions := range contextActions {
		for _, action := range actions {
			b.AddAction(contextId, action)
		}
	}
	return b
}

func (b *KeyBinder) IsBound(contextId string, msg tea.Msg) bool {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		_, ok := b.Get(contextId, msg.String())
		return ok
	default:
		return false
	}
}

func (b *KeyBinder) Get(contextId string, keyPress string) (press.KeyAction, bool) {
	key := NewKeyContext(contextId, keyPress)
	keyAction, ok := b.bindings[key]
	return keyAction, ok
}

func (b *KeyBinder) GetFooterVisible(contextId string) []press.KeyAction {
	var contextActions []press.KeyAction
	for key, action := range b.bindings {
		if key.ContextId == contextId && action.IsFooterVisible() {
			contextActions = append(contextActions, action)
		}
	}

	slices.SortFunc(contextActions, func(a, b press.KeyAction) int {
		orderDiff := cmp.Compare(a.SortOrder, b.SortOrder)
		if orderDiff == 0 {
			return cmp.Compare(a.Name, b.Name)
		} else {
			return orderDiff
		}
	})

	return contextActions
}

func (b *KeyBinder) Execute(contextId string, msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		key := NewKeyContext(contextId, msg.String())
		keyAction, ok := b.bindings[key]
		if ok {
			if !keyAction.Permanent {
				b.Remove(key)
			}
			return keyAction.ActionFunc()
		}
	}
	return nil
}
