package stats

import (
	"slices"
	"strings"

	tea "charm.land/bubbletea/v2"
	eff "github.com/randomvlad/trader-vlads/internal/appmod/stats/statuseffect"
	"github.com/randomvlad/trader-vlads/internal/component/tabs"
)

type Model struct {
	player PlayerStatsService
}

func NewTuiModel(player PlayerStatsService) *Model {
	return &Model{
		player: player,
	}
}

type PlayerStatsService interface {
	GetEffects() []eff.StatusEffect
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) View() tea.View {

	panel := tabs.NewTabPanel()
	panel.WriteLn("You are affected by:")

	effects := m.player.GetEffects()
	if len(effects) > 0 {

		// sort by turns left, then name
		slices.SortFunc(effects, func(first, second eff.StatusEffect) int {
			diff := first.GetTurns() - second.GetTurns()
			if diff == 0 {
				return strings.Compare(first.Name(), second.Name())
			} else {
				return diff
			}
		})

		for _, effect := range effects {
			panel.WriteLn("    " + effect.Name() + " : " + effect.View())
		}
	} else {
		panel.WriteLn("    Nothing")
	}

	return panel.RenderTeaView()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	//switch msg := msg.(type) {
	//case tea.KeyMsg:
	//
	//}

	return m, tea.Batch(cmds...)
}
