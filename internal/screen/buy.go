package screen

import (
	"errors"
	"fmt"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
	"github.com/randomvlad/trader-vlads/internal/component/gimme"
	"github.com/randomvlad/trader-vlads/internal/util"
)

type BuyScreen struct {
	form         *gimme.Form
	marketItems  []string
	marketPrices map[string]int
	money        int
	OnAborted    func()
	OnComplete   func(map[string]int)
}

func NewBuyScreen(marketItems []string, marketPrices map[string]int, money int) *BuyScreen {
	return &BuyScreen{
		marketItems:  marketItems,
		marketPrices: marketPrices,
		money:        money,
	}
}

func (b *BuyScreen) View() tea.View {
	return tea.NewView(style.Render(b.form.View().Content))
}

func (b *BuyScreen) Init() tea.Cmd {

	keyMap := huh.NewDefaultKeyMap()
	keyMap.Quit = key.NewBinding(key.WithKeys("esc"))
	keyMap.Input.Prev = key.NewBinding(key.WithKeys("up", "shift+tab"), key.WithHelp("up / shift+tab", "back"))
	keyMap.Input.Next = key.NewBinding(key.WithKeys("down", "enter", "tab"), key.WithHelp("down / enter / tab", "next"))

	// Future note: "enter" key submits the form when on last field. Want a more convenient and faster way to submit form. Allow to hit enter at any point

	inputFields := util.CreateFormInputFields(
		b.marketItems,
		func(name string) string {
			return fmt.Sprintf("Item: %v; Price: %v", name, util.FormatMoney(b.marketPrices[name]))
		},
		func(value string, input *gimme.Input) error {
			maxRange := b.money / b.marketPrices[input.GetKey()]
			return gimme.ValidateNumberStringInRange(value, 0, maxRange)
		})

	group := gimme.NewGroup(inputFields...).
		Title("Purchase Order:").
		WithValidation(func(group *gimme.Group) error {
			totalCost := 0
			for item, quantity := range group.GetValuesInt() {
				totalCost += b.marketPrices[item] * quantity
			}
			if totalCost > b.money {
				return errors.New(fmt.Sprintf("Requested purchase costs %v and exceeds available funds", util.FormatMoney(totalCost)))
			}

			return nil
		})

	b.form = gimme.NewForm(group).
		WithKeyMap(keyMap).
		WithShowHelp(true).
		WithShowErrors(true)

	return b.form.Init()
}

func (b *BuyScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	_, cmd := b.form.Update(msg)
	cmds = append(cmds, cmd)

	if b.form.State == gimme.StateAborted {
		b.OnAborted()
	} else if b.form.State == gimme.StateCompleted {
		inputValues := b.form.GetValuesInt()
		b.OnComplete(inputValues)
	}

	return b, tea.Batch(cmds...)
}
