package screen

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

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
	var sb strings.Builder
	sb.WriteString("Purchase Order:\n")
	sb.WriteString(b.form.View().Content)
	return tea.NewView(style.Render(sb.String()))
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
		})

	// Add min/max quantity validation on each input field
	for _, field := range inputFields {
		field.(*gimme.Input).Validate(func(inputValue string) error {
			if strings.TrimSpace(inputValue) == "" {
				return nil
			}

			minRange := 0
			maxRange := b.money / b.marketPrices[field.GetKey()]

			quantity, err := strconv.Atoi(inputValue)
			if err != nil || quantity < minRange || quantity > maxRange {
				return errors.New(fmt.Sprintf("Invalid value. Enter a number between %v - %v or leave blank", minRange, maxRange))
			}

			// TODO: rework this. introduce group level validation
			totalCost := b.getTotalCost()
			if totalCost > b.money {
				return errors.New(fmt.Sprintf("Total cost of %v exceeds available funds", util.FormatMoney(totalCost)))
			}

			return err
		})
	}

	b.form = gimme.NewForm(gimme.NewGroup(inputFields...)).
		WithShowHelp(true).
		WithKeyMap(keyMap).
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

func (b *BuyScreen) getTotalCost() int {
	inputValues := b.form.GetValuesInt()

	totalCost := 0
	for item, quantity := range inputValues {
		totalCost += b.marketPrices[item] * quantity
	}

	return totalCost
}
