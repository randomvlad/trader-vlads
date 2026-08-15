package actionfooter

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/randomvlad/trader-vlads/internal/appstyle"
)

type Model struct {
	footerType FooterType
	Actions    []string
}

type FooterType int

const (
	FooterStandalone FooterType = iota
	FooterTab
	FooterNoStyle
)

func NewModel(footerType FooterType, actions ...string) *Model {
	return &Model{footerType, actions}
}

func (m *Model) Render() string {

	var view strings.Builder
	view.WriteString("Actions: ")

	for index, action := range m.Actions {
		styledAction := lipgloss.StyleRanges(
			action,
			lipgloss.NewRange(0, 1, appstyle.StyleActionFirstLetter),
			lipgloss.NewRange(1, len(action), appstyle.NewAppStyle()),
		)

		view.WriteString(styledAction)

		isLast := index == len(m.Actions)-1
		if !isLast {
			view.WriteString(" • ")
		}
	}

	switch m.footerType {
	case FooterStandalone:
		return appstyle.StyleActionFooter.Render(view.String())
	case FooterTab:
		return appstyle.StyleActionFooterTab.Render(view.String())
	case FooterNoStyle:
		fallthrough
	default:
		return view.String()
	}
}
