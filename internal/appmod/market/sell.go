package market

import (
	"fmt"
	"slices"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
	"github.com/randomvlad/trader-vlads/internal/appstyle"
	"github.com/randomvlad/trader-vlads/internal/component/gimme"
	"github.com/randomvlad/trader-vlads/internal/util"
)

type SellScreen struct {
	form               *gimme.Form
	warehouseResources map[string]int
	marketPrices       map[string]int
	OnAborted          func()
	OnComplete         func(map[string]int)
}

func NewSellScreen(resources map[string]int, marketPrices map[string]int) *SellScreen {
	return &SellScreen{
		warehouseResources: resources,
		marketPrices:       marketPrices,
	}
}

func (s *SellScreen) View() tea.View {
	return tea.NewView(appstyle.StylePopup.Render(s.form.View().Content))
}

func (s *SellScreen) Init() tea.Cmd {
	keyMap := huh.NewDefaultKeyMap()
	keyMap.Quit = key.NewBinding(key.WithKeys("esc"))
	keyMap.Input.Prev = key.NewBinding(key.WithKeys("up", "shift+tab"), key.WithHelp("up / shift+tab", "back"))
	keyMap.Input.Next = key.NewBinding(key.WithKeys("down", "enter", "tab"), key.WithHelp("down / enter / tab", "next"))

	// Future note: "enter" key submits the form when on last field. Want a more convenient and faster way to submit form. Allow to hit enter at any point

	var items []string
	for item, quantity := range s.warehouseResources {
		if quantity > 0 {
			items = append(items, item)
		}
	}

	if items == nil {
		items = []string{}
	}

	slices.Sort(items)

	inputFields := util.CreateFormInputFields(
		items,
		func(name string) string {
			return fmt.Sprintf("Item: %v; Price: %v; # Available: %v",
				name, util.FormatCurrency(s.marketPrices[name]), s.warehouseResources[name])
		},
		func(value string, input *gimme.Input) error {
			maxRange := s.warehouseResources[input.GetKey()]
			return gimme.ValidateNumberStringInRange(value, 0, maxRange)
		})

	group := gimme.NewGroup(inputFields...).Title("Sell Order:")

	s.form = gimme.NewForm(group).
		WithKeyMap(keyMap).
		WithShowHelp(true).
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
		inputValues := s.form.GetValuesInt()
		s.OnComplete(inputValues)
	}

	return s, tea.Batch(cmds...)
}
