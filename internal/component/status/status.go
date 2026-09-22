package status

import (
	"strconv"

	"charm.land/lipgloss/v2"
	"github.com/randomvlad/trader-vlads/internal/appstyle"
	"github.com/randomvlad/trader-vlads/internal/util"
	"github.com/randomvlad/trader-vlads/internal/util/stringutil"
)

type Model struct {
	badgeTurn  badge
	badgeGold  badge
	styleBadge lipgloss.Style
}

type badge struct {
	name string
}

func New() Model {
	return Model{
		badgeTurn: badge{
			name: "Week",
		},
		badgeGold: badge{
			name: "Gold",
		},
		styleBadge: appstyle.NewAppStyle().Padding(0, 1).MarginRight(2),
	}
}

func (m *Model) Render(week, gold int) string {
	view := stringutil.NewBuilder().
		Write(m.badgeTurn.render(strconv.Itoa(week), m.styleBadge)).
		Write(m.badgeGold.render(util.FormatCurrency(gold), m.styleBadge))

	styleStatusBar := appstyle.NewAppStyle().
		Border(lipgloss.RoundedBorder(), true, true, false, true).
		BorderForeground(appstyle.AppBorderColor).
		Width(appstyle.AppWidth).
		PaddingBottom(1)

	return view.StringStyle(styleStatusBar)
}

func (b *badge) render(value string, style lipgloss.Style) string {
	return style.Render(b.name + ": " + value)
}
