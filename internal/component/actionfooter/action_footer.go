package actionfooter

import (
	"charm.land/lipgloss/v2"
	"github.com/randomvlad/trader-vlads/internal/appmod/keybind/appaction"
	"github.com/randomvlad/trader-vlads/internal/appstyle"
	"github.com/randomvlad/trader-vlads/internal/util/stringutil"
)

type Model struct {
	footerType FooterType
	appActions []appaction.AppAction
	visual     *modelVisual
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

func NewModel(footerType FooterType, actions ...appaction.AppAction) *Model {
	return &Model{
		footerType: footerType,
		appActions: actions,
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

func (m *Model) Render() string {
	view := stringutil.NewBuilder().Write("Actions: ")

	if len(m.appActions) > 0 {
		for index, appAction := range m.appActions {
			view.WriteStyleRanges(
				appAction.GetName(),
				lipgloss.NewRange(0, 1, m.visual.styleFirstLetter),
				lipgloss.NewRange(1, len(appAction.GetName()), appstyle.NewAppStyle()),
			)

			isLast := index == len(m.appActions)-1
			if !isLast {
				view.Write(" • ")
			}
		}
	} else {
		view.Write("None")
	}

	borderStyle := m.visual.styleBorder.Width(m.visual.width)

	return view.StringStyle(borderStyle)
}

func getStyleBorder(footerType FooterType) lipgloss.Style {
	switch footerType {
	case FooterStandalone:
		return lipgloss.NewStyle().
			Padding(0, 2).
			Border(lipgloss.RoundedBorder(), true).
			BorderForeground(appstyle.AppBorderColor)
	case FooterPanel:
		return lipgloss.NewStyle().
			Padding(0, 2).
			Border(lipgloss.RoundedBorder(), false, true, true, true).
			BorderForeground(appstyle.AppBorderColor)
	case FooterNoStyle:
		fallthrough
	default:
		return lipgloss.NewStyle()
	}
}
