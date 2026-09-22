package actionfooter

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/randomvlad/trader-vlads/internal/appmod/keybind"
	"github.com/randomvlad/trader-vlads/internal/appmod/keybind/press"
	"github.com/randomvlad/trader-vlads/internal/appstyle"
	"github.com/randomvlad/trader-vlads/internal/util/stringutil"
)

type Model struct {
	footerType                FooterType
	keyBinder                 *keybind.KeyBinder
	displayRightActionContext string
	actions                   []press.KeyAction
	visual                    *modelVisual
}

type modelVisual struct {
	width             int
	styleBorder       lipgloss.Style
	styleFirstLetter  lipgloss.Style
	styleOtherLetters lipgloss.Style
}

type FooterType int

const (
	FooterStandalone FooterType = iota
	FooterPanel
	FooterNoStyle
)

func NewModel(footerType FooterType, actions ...press.KeyAction) *Model {
	return &Model{
		footerType: footerType,
		actions:    actions,
		visual: &modelVisual{
			width:             appstyle.AppWidth,
			styleBorder:       getStyleBorder(footerType),
			styleFirstLetter:  appstyle.NewAppStyle().Bold(true).Underline(true),
			styleOtherLetters: appstyle.NewAppStyle(),
		},
	}
}

func (m *Model) WithWidth(value int) *Model {
	m.visual.width = value
	return m
}

func (m *Model) WithStyle(style lipgloss.Style) *Model {
	m.visual.styleBorder = style
	return m
}

func (m *Model) WithKeyBinder(keyBinder *keybind.KeyBinder) *Model {
	m.keyBinder = keyBinder
	return m
}

func (m *Model) WithDisplayRightActionContext(actionContext string) *Model {
	m.displayRightActionContext = actionContext
	return m
}

func (m *Model) SetActions(actions ...press.KeyAction) {
	m.actions = actions
}

func (m *Model) Render(actionsContext string) string {
	var viewParts []string

	viewLeft := stringutil.NewBuilder().Write("Actions: ")
	actionsLeft := m.getActions(actionsContext)
	if len(actionsLeft) > 0 {
		viewLeft.Write(m.formatActions(actionsLeft))
	} else {
		viewLeft.Write("None")
	}
	viewParts = append(viewParts, viewLeft.String())

	actionsRight := m.getDisplayRightActions()
	if len(actionsRight) > 0 {
		viewRight := stringutil.NewBuilder().
			WriteStyle(" │ ", lipgloss.NewStyle().Foreground(appstyle.AppBorderColor)).
			Write(m.formatActions(actionsRight))

		spacerWidth := m.visual.width - m.visual.styleBorder.GetHorizontalFrameSize() - viewLeft.Width() - viewRight.Width()
		if spacerWidth > 0 {
			viewParts = append(viewParts, strings.Repeat(" ", spacerWidth))
		}

		viewParts = append(viewParts, viewRight.String())
	}

	footerContent := lipgloss.JoinHorizontal(lipgloss.Bottom, viewParts...)
	borderStyle := m.visual.styleBorder.Width(m.visual.width)
	return borderStyle.Render(footerContent)
}

func (m *Model) formatActions(actions []press.KeyAction) string {
	content := stringutil.NewBuilder()

	for index, action := range actions {
		content.WriteStyleRanges(
			action.Name,
			lipgloss.NewRange(0, 1, m.visual.styleFirstLetter),
			lipgloss.NewRange(1, len(action.Name), appstyle.NewAppStyle()),
		)
		isLast := index == len(actions)-1
		if !isLast {
			content.Write(" • ")
		}
	}

	return content.String()
}

func (m *Model) getActions(actionsContext string) []press.KeyAction {
	if m.keyBinder != nil && actionsContext != "" {
		return m.keyBinder.GetFooterVisible(actionsContext)
	} else {
		return m.actions // TODO: once refactored, m.actions should no longer be needed
	}
}

func (m *Model) getDisplayRightActions() []press.KeyAction {
	if m.keyBinder != nil && m.displayRightActionContext != "" {
		return m.keyBinder.GetFooterVisible(m.displayRightActionContext)
	} else {
		return []press.KeyAction{}
	}
}

func getStyleBorder(footerType FooterType) lipgloss.Style {
	switch footerType {
	case FooterStandalone:
		return lipgloss.NewStyle().
			Padding(0, 1).
			Border(lipgloss.RoundedBorder(), true).
			BorderForeground(appstyle.AppBorderColor)
	case FooterPanel:
		return lipgloss.NewStyle().
			Padding(0, 1).
			Border(lipgloss.RoundedBorder(), false, true, true, true).
			BorderForeground(appstyle.AppBorderColor)
	case FooterNoStyle:
		fallthrough
	default:
		return lipgloss.NewStyle()
	}
}
