package event

import (
	"strconv"

	tea "charm.land/bubbletea/v2"
	eq "github.com/randomvlad/trader-vlads/internal/appmod/equipment"
	"github.com/randomvlad/trader-vlads/internal/util"
)

type StarterSetStory struct {
	phaseIndex int
	Name       string
	items      []string
	resources  map[string]int
	Complete   bool
	player     PlayerTurnService
	random     *util.RandomGenerator
}

func NewStarterSetStory(name string, player PlayerTurnService, r *util.RandomGenerator) *StarterSetStory {
	return &StarterSetStory{
		Name:   name,
		player: player,
		random: r,
		items: []string{
			"copper ring of a novice",
			"gray cotton tunic",
			"worn trousers",
			"brown leather sandals",
			"a potion of Beginner's Luck 🍀",
			"a jar of spicy pickles",
		},
		resources: map[string]int{
			"Wood":  3,
			"Stone": 3,
		},
	}
}

func (s *StarterSetStory) Render() string {
	var render util.StringBuilder
	switch s.phaseIndex {
	case 0:
		render.WriteLn("The Guild of Merchants has sent a standard edition wooden chest to get you started.")
	case 1:
		render.WriteLn("You open the chest and look inside:")

		for resource, count := range s.resources {
			render.Writef("    +%s %s\n", strconv.Itoa(count), resource)
		}

		for _, itemName := range s.items {
			render.Tab().WriteLn(itemName)
		}
	case 2:
		render.Write("You place the items in your inventory. Might be a good idea to try them on next.")
	}
	return render.String()
}

func (s *StarterSetStory) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "o", "O":
			s.phaseIndex = 1
			break
		case "t", "T":
			s.phaseIndex = 2

			for resource, count := range s.resources {
				s.player.AddResourceQuantity(resource, count)
			}

			eqObjects := eq.Forge.Make(s.random, s.items...)
			for _, object := range eqObjects {
				s.player.AddInventory(object)
			}

		case "c", "C":
			s.Complete = true
		}
	}

	return nil, nil
}

func (s *StarterSetStory) GetAvailableActions() []string {
	switch s.phaseIndex {
	case 0:
		return []string{"Open"}
	case 1:
		return []string{"Take"}
	default:
		return []string{"Continue"}
	}
}
