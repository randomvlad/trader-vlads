package screen

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
	"github.com/randomvlad/trader-vlads/internal/component/gimme"
	"github.com/randomvlad/trader-vlads/internal/util"
)

type SellScreen struct {
	form         *gimme.Form
	inventory    map[string]int
	marketPrices map[string]int
	OnAborted    func()
	OnComplete   func(map[string]int)
}

func NewSellScreen(inventory map[string]int, marketPrices map[string]int) *SellScreen {
	return &SellScreen{
		inventory:    inventory,
		marketPrices: marketPrices,
	}
}

func (s *SellScreen) View() tea.View {
	var sb strings.Builder
	sb.WriteString("Sell Order:\n")
	sb.WriteString(s.form.View().Content)
	return tea.NewView(style.Render(sb.String()))
}

func (s *SellScreen) Init() tea.Cmd {
	keyMap := huh.NewDefaultKeyMap()
	keyMap.Quit = key.NewBinding(key.WithKeys("esc"))
	keyMap.Input.Prev = key.NewBinding(key.WithKeys("up", "shift+tab"), key.WithHelp("up / shift+tab", "back"))
	keyMap.Input.Next = key.NewBinding(key.WithKeys("down", "enter", "tab"), key.WithHelp("down / enter / tab", "next"))

	// Future note: "enter" key submits the form when on last field. Want a more convenient and faster way to submit form. Allow to hit enter at any point

	var items []string
	for item, quantity := range s.inventory {
		if quantity > 0 {
			items = append(items, item)
		}
	}

	// TODO: handle case where inventory is empty. Currently results in a panic crash

	slices.Sort(items)

	inputFields := util.CreateFormInputFields(
		items,
		func(name string) string {
			return fmt.Sprintf("Item: %v; Price: %v; # Available: %v",
				name, util.FormatMoney(s.marketPrices[name]), s.inventory[name])
		})

	s.form = gimme.NewForm(gimme.NewGroup(inputFields...)).
		WithShowHelp(true).
		WithKeyMap(keyMap).
		WithShowErrors(true)

	return s.form.Init()
}

func (s *SellScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	_, cmd := s.form.Update(msg)
	cmds = append(cmds, cmd)

	if s.form.State == gimme.StateAborted {
		s.OnAborted()
	} else if s.form.State == gimme.StateCompleted {
		// TODO: need validation to enforce selling up to max of inventory
		inputFields := slices.Collect(maps.Keys(s.inventory))
		result := s.form.GetResultInts(inputFields)
		s.OnComplete(result)
	}

	return s, tea.Batch(cmds...)
}
