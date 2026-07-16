package status

import (
	"strconv"

	"charm.land/lipgloss/v2"
	"github.com/randomvlad/trader-vlads/internal/util"
)

var (
	styleBadge = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#04B575")).
		Padding(0, 2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#04B575"))
)

type Model struct {
	badgeTurn  badge
	badgeMoney badge
}

type badge struct {
	name string
}

func New() Model {
	return Model{
		badgeTurn: badge{
			name: "Week #",
		},
		badgeMoney: badge{
			name: "Money",
		},
	}
}

func (m *Model) Render(week, money int) string {

	badgeWeeks := m.badgeTurn.render(strconv.Itoa(week))
	badgeMoney := m.badgeMoney.render(util.FormatMoney(money))

	return lipgloss.JoinHorizontal(lipgloss.Top, badgeWeeks, badgeMoney) + "\n"
}

func (b *badge) render(value string) string {
	return styleBadge.Render(b.name + ": " + value)
}
