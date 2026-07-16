package gimme

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
)

// Field is a primitive of a form.
//
// A field represents a single input control on a form such as a text input,
// confirm button, select option, etc...
//
// Each field implements the Bubble Tea Model interface.
type Field interface {
	tea.Model

	// Bubble Tea Events
	Blur() tea.Cmd
	Focus() tea.Cmd

	// Errors and Validation
	Error() error

	// Skip returns whether this input should be skipped or not.
	Skip() bool

	// Zoom returns whether this input should be zoomed or not.
	// Zoom allows the field to take focus of the group / form height.
	Zoom() bool

	// KeyBinds returns help keybindings.
	KeyBinds() []key.Binding

	// WithTheme sets the theme on a field.
	WithTheme(huh.Theme) Field

	// WithKeyMap sets the keymap on a field.
	WithKeyMap(*huh.KeyMap) Field

	// WithWidth sets the width of a field.
	WithWidth(int) Field

	// WithHeight sets the height of a field.
	WithHeight(int) Field

	// WithPosition tells the field the index of the group and position it is in.
	WithPosition(FieldPosition) Field

	// GetKey returns the field's key.
	GetKey() string

	// GetValue returns the field's value.
	GetValue() any
}
