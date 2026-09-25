package actionfooter

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/randomvlad/trader-vlads/internal/appmod/keybind"
	"github.com/randomvlad/trader-vlads/internal/appmod/keybind/press"
	"github.com/randomvlad/trader-vlads/internal/appstyle"
	"github.com/randomvlad/trader-vlads/internal/util/stringutil"
)

type ActionFooter struct {
	footerType  FooterType
	keyBinder   *keybind.KeyBinder
	styleConfig actionFooterStyleConfig
}

type actionFooterStyleConfig struct {
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

func NewPanelFooter(keyBinder *keybind.KeyBinder, width int) *ActionFooter {
	footerType := FooterPanel
	return &ActionFooter{
		footerType: footerType,
		keyBinder:  keyBinder,
		styleConfig: actionFooterStyleConfig{
			width:             width,
			styleBorder:       getStyleBorder(footerType),
			styleFirstLetter:  appstyle.NewAppStyle().Bold(true).Underline(true),
			styleOtherLetters: appstyle.NewAppStyle(),
		},
	}
}

func (m *ActionFooter) Render(actionContextLeft, actionContextRight string) string {
	var viewParts []string

	viewLeft := stringutil.NewBuilder().Write("Actions: ")
	actionsLeft := m.getFooterVisibleActions(actionContextLeft)
	if len(actionsLeft) > 0 {
		viewLeft.Write(m.formatActions(actionsLeft))
	} else {
		viewLeft.Write("None")
	}
	viewParts = append(viewParts, viewLeft.String())

	actionsRight := m.getFooterVisibleActions(actionContextRight)
	if len(actionsRight) > 0 {
		viewRight := stringutil.NewBuilder().
			WriteStyle(" │ ", lipgloss.NewStyle().Foreground(appstyle.AppBorderColor)).
			Write(m.formatActions(actionsRight))

		spacerWidth := m.styleConfig.width - m.styleConfig.styleBorder.GetHorizontalFrameSize() - viewLeft.Width() - viewRight.Width()
		if spacerWidth > 0 {
			viewParts = append(viewParts, strings.Repeat(" ", spacerWidth))
		}

		viewParts = append(viewParts, viewRight.String())
	}

	footerContent := lipgloss.JoinHorizontal(lipgloss.Bottom, viewParts...)
	borderStyle := m.styleConfig.styleBorder.Width(m.styleConfig.width)
	return borderStyle.Render(footerContent)
}

func (m *ActionFooter) formatActions(actions []press.KeyAction) string {
	content := stringutil.NewBuilder()

	for index, action := range actions {
		content.WriteStyleRanges(
			action.Name,
			lipgloss.NewRange(0, 1, m.styleConfig.styleFirstLetter),
			lipgloss.NewRange(1, len(action.Name), appstyle.NewAppStyle()),
		)
		isLast := index == len(actions)-1
		if !isLast {
			content.Write(" • ")
		}
	}

	return content.String()
}

func (m *ActionFooter) getFooterVisibleActions(actionsContext string) []press.KeyAction {
	if m.keyBinder != nil && actionsContext != "" {
		return m.keyBinder.GetFooterVisible(actionsContext)
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
