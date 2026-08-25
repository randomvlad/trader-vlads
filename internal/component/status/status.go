package status

import (
	"strconv"

	"charm.land/lipgloss/v2"
	"github.com/randomvlad/trader-vlads/internal/appstyle"
	"github.com/randomvlad/trader-vlads/internal/util"
)

type Model struct {
	badgeTurn badge
	badgeGold badge
}

type badge struct {
	name string
}

func New() Model {
	return Model{
		badgeTurn: badge{
			name: "Week #",
		},
		badgeGold: badge{
			name: "Gold",
		},
	}
}

func (m *Model) Render(week, gold int) string {

	badgeWeeks := m.badgeTurn.render(strconv.Itoa(week))
	badgeGold := m.badgeGold.render(util.FormatCurrency(gold))

	return lipgloss.JoinHorizontal(lipgloss.Top, badgeWeeks, badgeGold) + "\n"
}

func (b *badge) render(value string) string {
	return appstyle.StyleBadge.Render(b.name + ": " + value)
}
