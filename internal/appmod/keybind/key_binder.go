package keybind

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/randomvlad/trader-vlads/internal/appmod/keybind/press"
)

type KeyBinder struct {
	bindings map[KeyContext]func() tea.Cmd
}

func NewKeyBinder() *KeyBinder {
	return &KeyBinder{
		bindings: make(map[KeyContext]func() tea.Cmd),
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

func (r *KeyBinder) Add(contextId string, key string, f func() tea.Cmd) *KeyBinder {
	r.bindings[NewKeyContext(contextId, key)] = f
	return r
}

func (r *KeyBinder) Remove(k KeyContext) *KeyBinder {
	delete(r.bindings, k)
	return r
}

func (b *KeyBinder) AddAction(contextId string, action press.KeyAction) *KeyBinder {
	return b.Add(contextId, action.KeyPress, func() tea.Cmd {
		action.ActionFunc()
		return nil
	})
}

func (b *KeyBinder) AddActions(contextActions map[string][]press.KeyAction) *KeyBinder {
	for contextId, actions := range contextActions {
		for _, action := range actions {
			b.AddAction(contextId, action)
		}
	}
	return b
}

func (r *KeyBinder) IsBound(contextId string, msg tea.Msg) bool {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		_, ok := r.Get(contextId, msg.String())
		return ok
	default:
		return false
	}
}

func (r *KeyBinder) Get(contextId string, keyPress string) (func() tea.Cmd, bool) {
	key := NewKeyContext(contextId, keyPress)
	executeFunc, ok := r.bindings[key]
	return executeFunc, ok
}

// TODO: fire once and remove vs permanent
func (r *KeyBinder) Execute(contextId string, msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		key := NewKeyContext(contextId, msg.String())
		executeFunc, ok := r.bindings[key]
		if ok {
			delete(r.bindings, key)
			return executeFunc()
		}
	}
	return nil
}
