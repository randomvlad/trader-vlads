package gimme

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
)

const defaultWidth = 80

// Internal ID management. Used during animating to ensure that frame messages
// are received only by spinner components that sent them.
var (
	lastID int
	idMtx  sync.Mutex
)

// Return the next ID we should use on the Model.
func nextID() int {
	idMtx.Lock()
	defer idMtx.Unlock()
	lastID++
	return lastID
}

// FormState represents the current state of the form.
type FormState int

const (
	// StateNormal is when the user is completing the form.
	StateNormal FormState = iota

	// StateCompleted is when the user has completed the form.
	StateCompleted

	// StateAborted is when the user has aborted the form.
	StateAborted
)

// ErrUserAborted is the error returned when a user exits the form before submitting.
var ErrUserAborted = errors.New("user aborted")

// ErrTimeout is the error returned when the timeout is reached.
var ErrTimeout = errors.New("timeout")

// Form is a collection of groups that are displayed one at a time on a "page".
//
// The form can navigate between groups and is complete once all the groups are
// complete.
type Form struct {
	// collection of groups
	selector *Selector[*Group]

	results map[string]any

	// callbacks
	SubmitCmd tea.Cmd
	CancelCmd tea.Cmd

	State FormState

	// options
	width      int
	height     int
	theme      huh.Theme
	hasDarkBg  bool
	keymap     *huh.KeyMap
	timeout    time.Duration
	teaOptions []tea.ProgramOption

	layout Layout
}

// NewForm returns a form with the given groups and default themes and
// keybindings.
//
// Use With* methods to customize the form with options, such as setting
// different themes and keybindings.
func NewForm(groups ...*Group) *Form {

	f := &Form{
		selector: NewSelector(groups),
		keymap:   huh.NewDefaultKeyMap(),
		results:  make(map[string]any),
		layout:   LayoutDefault,
		teaOptions: []tea.ProgramOption{
			tea.WithOutput(os.Stderr),
		},
	}

	f.WithKeyMap(f.keymap)
	f.WithWidth(f.width)
	f.WithHeight(f.height)
	f.UpdateFieldPositions()

	return f
}

// FieldPosition is positional information about the given field and form.
type FieldPosition struct {
	Group      int
	Field      int
	FirstField int
	LastField  int
	GroupCount int
	FirstGroup int
	LastGroup  int
}

// IsFirst returns whether a field is the form's first field.
func (p FieldPosition) IsFirst() bool {
	return p.Field == p.FirstField && p.Group == p.FirstGroup
}

// IsLast returns whether a field is the form's last field.
func (p FieldPosition) IsLast() bool {
	return p.Field == p.LastField && p.Group == p.LastGroup
}

// nextGroupMsg is a message to move to the next group.
type nextGroupMsg struct{}

// prevGroupMsg is a message to move to the previous group.
type prevGroupMsg struct{}

// nextGroup is the command to move to the next group.
func nextGroup() tea.Msg {
	return nextGroupMsg{}
}

// prevGroup is the command to move to the previous group.
func prevGroup() tea.Msg {
	return prevGroupMsg{}
}

// WithShowHelp sets whether or not the form should show help.
//
// This allows the form groups and field to show what keybindings are available
// to the user.
func (f *Form) WithShowHelp(v bool) *Form {
	f.selector.Range(func(_ int, group *Group) bool {
		group.WithShowHelp(v)
		return true
	})
	return f
}

// WithShowErrors sets whether or not the form should show errors.
//
// This allows the form groups and fields to show errors when the Validate
// function returns an error.
func (f *Form) WithShowErrors(v bool) *Form {
	f.selector.Range(func(_ int, group *Group) bool {
		group.WithShowErrors(v)
		return true
	})
	return f
}

// WithTheme sets the theme on a form.
//
// This allows all groups and fields to be themed consistently, however themes
// can be applied to each group and field individually for more granular
// control.
func (f *Form) WithTheme(theme huh.Theme) *Form {
	if theme == nil {
		return f
	}
	f.theme = theme
	f.selector.Range(func(_ int, group *Group) bool {
		group.WithTheme(theme)
		return true
	})
	return f
}

// WithKeyMap sets the keymap on a form.
//
// This allows customization of the form key bindings.
func (f *Form) WithKeyMap(keymap *huh.KeyMap) *Form {
	if keymap == nil {
		return f
	}
	f.keymap = keymap
	f.selector.Range(func(_ int, group *Group) bool {
		group.WithKeyMap(keymap)
		return true
	})
	f.UpdateFieldPositions()
	return f
}

// WithWidth sets the width of a form.
//
// This allows all groups and fields to be sized consistently, however width
// can be applied to each group and field individually for more granular
// control.
func (f *Form) WithWidth(width int) *Form {
	if width <= 0 {
		return f
	}
	f.width = width
	f.selector.Range(func(_ int, group *Group) bool {
		width := f.layout.GroupWidth(f, group, width)
		group.WithWidth(width)
		return true
	})
	return f
}

// WithHeight sets the height of a form.
func (f *Form) WithHeight(height int) *Form {
	if height <= 0 {
		return f
	}
	f.height = height
	f.selector.Range(func(_ int, group *Group) bool {
		group.WithHeight(height)
		return true
	})
	return f
}

// WithTimeout sets the duration for the form to be killed.
func (f *Form) WithTimeout(t time.Duration) *Form {
	f.timeout = t
	return f
}

// WithProgramOptions sets the tea options of the form.
func (f *Form) WithProgramOptions(opts ...tea.ProgramOption) *Form {
	f.teaOptions = opts
	return f
}

// WithLayout sets the layout on a form.
//
// This allows customization of the form group layout.
func (f *Form) WithLayout(layout Layout) *Form {
	f.layout = layout
	return f
}

// UpdateFieldPositions sets the position on all the fields.
func (f *Form) UpdateFieldPositions() *Form {
	firstGroup := 0
	lastGroup := f.selector.Total() - 1

	// determine the first non-hidden group.
	f.selector.Range(func(_ int, g *Group) bool {
		if !g.IsHidden() {
			return false
		}
		firstGroup++
		return true
	})

	// determine the last non-hidden group.
	f.selector.ReverseRange(func(_ int, g *Group) bool {
		if !g.IsHidden() {
			return false
		}
		lastGroup--
		return true
	})

	f.selector.Range(func(g int, group *Group) bool {
		// determine the first non-skippable field.
		var firstField int
		group.selector.Range(func(_ int, field Field) bool {
			if !field.Skip() || group.selector.Total() == 1 {
				return false
			}
			firstField++
			return true
		})

		// determine the last non-skippable field.
		var lastField int
		group.selector.ReverseRange(func(i int, field Field) bool {
			lastField = i
			if !field.Skip() || group.selector.Total() == 1 {
				return false
			}
			return true
		})

		group.selector.Range(func(i int, field Field) bool {
			field.WithPosition(FieldPosition{
				Group:      g,
				Field:      i,
				FirstField: firstField,
				LastField:  lastField,
				FirstGroup: firstGroup,
				LastGroup:  lastGroup,
			})
			return true
		})

		return true
	})
	return f
}

// Errors returns the current groups' errors.
func (f *Form) Errors() []error {
	return f.selector.Selected().Errors()
}

// Help returns the current groups' help.
func (f *Form) Help() help.Model {
	return f.selector.Selected().help
}

// KeyBinds returns the current fields' keybinds.
func (f *Form) KeyBinds() []key.Binding {
	group := f.selector.Selected()
	return group.selector.Selected().KeyBinds()
}

// GetFieldKeys Get list of form field keys
func (f *Form) GetFieldKeys() []string {
	var fieldKeys []string
	for _, group := range f.selector.items {
		for _, field := range group.selector.items {
			fieldKeys = append(fieldKeys, field.GetKey())
		}
	}
	return fieldKeys
}

// Get returns a result from the form.
func (f *Form) Get(key string) any {
	return f.results[key]
}

// GetString returns a result as a string from the form.
func (f *Form) GetString(key string) string {
	v, ok := f.results[key].(string)
	if !ok {
		return ""
	}
	return v
}

// GetValuesInt input field values converted to integers
func (f *Form) GetValuesInt(fieldKeys ...string) map[string]int {

	if len(fieldKeys) == 0 {
		fieldKeys = f.GetFieldKeys()
	}

	result := make(map[string]int, len(fieldKeys))

	for _, fieldKey := range fieldKeys {
		integer, _ := strconv.Atoi(f.GetString(fieldKey))
		result[fieldKey] = integer
	}

	return result
}

// GetInt returns a result as a int from the form.
func (f *Form) GetInt(key string) int {
	v, ok := f.results[key].(int)
	if !ok {
		return 0
	}
	return v
}

// GetBool returns a result as a string from the form.
func (f *Form) GetBool(key string) bool {
	v, ok := f.results[key].(bool)
	if !ok {
		return false
	}
	return v
}

// NextGroup moves the form to the next group.
func (f *Form) NextGroup() tea.Cmd {
	_, cmd := f.Update(nextGroup())
	return cmd
}

// PrevGroup moves the form to the next group.
func (f *Form) PrevGroup() tea.Cmd {
	_, cmd := f.Update(prevGroup())
	return cmd
}

// NextField moves the form to the next field.
func (f *Form) NextField() tea.Cmd {
	_, cmd := f.Update(NextField())
	return cmd
}

// PrevField moves the form to the next field.
func (f *Form) PrevField() tea.Cmd {
	_, cmd := f.Update(PrevField())
	return cmd
}

// GetFocusedField returns the focused form field.
func (f *Form) GetFocusedField() Field {
	return f.selector.Selected().selector.Selected()
}

// Init initializes the form.
func (f *Form) Init() tea.Cmd {
	var cmds []tea.Cmd
	f.selector.Range(func(i int, group *Group) bool {
		if i == 0 {
			group.active = true
		}
		cmds = append(cmds, group.Init())
		return true
	})

	if f.selector.Selected().IsHidden() {
		cmds = append(cmds, nextGroup)
	}

	cmds = append(cmds, tea.RequestWindowSize)
	return tea.Sequence(cmds...)
}

// Update updates the form.
func (f *Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// If the form is aborted or completed there's no need to update it.
	if f.State != StateNormal {
		return f, nil
	}

	group := f.selector.Selected()

	switch msg := msg.(type) {
	case tea.BackgroundColorMsg:
		f.hasDarkBg = msg.IsDark()
	case tea.WindowSizeMsg:
		if f.width == 0 {
			f.selector.Range(func(_ int, group *Group) bool {
				width := f.layout.GroupWidth(f, group, msg.Width)
				group.WithWidth(width)
				return true
			})
		}
		if f.height == 0 {
			// calculate the needed height, which is the height of the
			// highest group, accounting for the width, wraps, etc.
			neededHeight := 0
			f.selector.Range(func(_ int, group *Group) bool {
				neededHeight = max(neededHeight, group.rawHeight())
				return true
			})

			f.selector.Range(func(_ int, group *Group) bool {
				group.WithHeight(min(neededHeight, msg.Height))
				return true
			})
		}

	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, f.keymap.Quit):
			f.State = StateAborted
			return f, f.CancelCmd
		}

	case nextFieldMsg:
		// Form is progressing to the next field, let's save the value of the current field.
		field := group.selector.Selected()
		f.results[field.GetKey()] = field.GetValue()

	case nextGroupMsg:
		if len(group.Errors()) > 0 {
			return f, nil
		}

		if f.selector.OnLast() {
			f.State = StateCompleted
			return f, f.SubmitCmd
		}

		for i := f.selector.Index() + 1; i < f.selector.Total(); i++ {
			if !f.selector.Get(i).IsHidden() {
				f.selector.SetIndex(i)
				break
			}
			// all subsequent groups are hidden, so we must act as
			// if we were in the last one.
			if i == f.selector.Total()-1 {
				f.State = StateCompleted
				return f, f.SubmitCmd
			}
		}
		f.selector.Selected().active = true
		return f, f.selector.Selected().Init()

	case prevGroupMsg:
		if len(group.Errors()) > 0 {
			return f, nil
		}

		for i := f.selector.Index() - 1; i >= 0; i-- {
			if !f.selector.Get(i).IsHidden() {
				f.selector.SetIndex(i)
				break
			}
		}

		f.selector.Selected().active = true
		return f, f.selector.Selected().Init()
	}

	m, cmd := group.Update(msg)
	if len(group.Errors()) > 0 {
		return f, cmd
	}

	f.selector.Set(f.selector.Index(), m.(*Group))

	// A user input a key, this could hide or show other groups,
	// let's update all of their positions.
	switch msg.(type) {
	case tea.KeyPressMsg:
		f.UpdateFieldPositions()
	}

	return f, cmd
}

func (f *Form) getTheme() *huh.Styles {
	if f.theme != nil {
		return f.theme.Theme(f.hasDarkBg)
	}
	return huh.ThemeCharm(f.hasDarkBg)
}

func (f *Form) styles() huh.FormStyles {
	return f.getTheme().Form
}

// View renders the form.
func (f *Form) View() tea.View {
	switch f.State {
	case StateAborted, StateCompleted:
		return tea.NewView("")
	case StateNormal:
		fallthrough
	default:
		return tea.NewView(f.styles().Base.Render(f.layout.View(f)))
	}
}

// Run runs the form.
func (f *Form) Run() error {
	return f.RunWithContext(context.Background())
}

// RunWithContext runs the form with the given context.
func (f *Form) RunWithContext(ctx context.Context) error {
	f.SubmitCmd = tea.Quit
	f.CancelCmd = tea.Interrupt

	if f.selector.Total() == 0 {
		return nil
	}

	return f.run(ctx)
}

// run runs the form in normal mode.
func (f *Form) run(ctx context.Context) error {
	var cancel context.CancelFunc
	if f.timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, f.timeout)
		defer cancel()
	}

	f.teaOptions = append(f.teaOptions, tea.WithContext(ctx))
	_, err := tea.NewProgram(f, f.teaOptions...).Run()
	if f.State == StateAborted || errors.Is(err, tea.ErrInterrupted) {
		return ErrUserAborted
	}
	if errors.Is(err, tea.ErrProgramKilled) {
		return ErrTimeout
	}
	if err != nil {
		return fmt.Errorf("huh: %w", err)
	}
	return nil
}
