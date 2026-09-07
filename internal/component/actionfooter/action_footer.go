package actionfooter

import (
	"charm.land/lipgloss/v2"
	"github.com/randomvlad/trader-vlads/internal/appstyle"
	"github.com/randomvlad/trader-vlads/internal/util/stringutil"
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

	view := stringutil.NewBuilder().Write("Actions: ")

	if len(m.Actions) > 0 {
		for index, action := range m.Actions {
			styledAction := lipgloss.StyleRanges(
				action,
				lipgloss.NewRange(0, 1, appstyle.StyleActionFirstLetter),
				lipgloss.NewRange(1, len(action), appstyle.NewAppStyle()),
			)

			view.Write(styledAction)

			isLast := index == len(m.Actions)-1
			if !isLast {
				view.Write(" • ")
			}
		}
	} else {
		view.Write("None")
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
