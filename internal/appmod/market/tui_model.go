package market

import (
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/randomvlad/trader-vlads/internal/component/tabs"
	"github.com/randomvlad/trader-vlads/internal/util"
)

type Model struct {
	Market         *Market
	playerService  PlayerService
	toastMessenger ToastMessenger
	Resources      map[string]int
	state          marketModuleState
	screenBuy      *BuyScreen
	screenSell     *SellScreen
}

func NewTuiModel(market *Market, player PlayerService, toast ToastMessenger) *Model {
	return &Model{
		Market:         market,
		playerService:  player,
		toastMessenger: toast,
	}
}

type marketModuleState int

const (
	stateList marketModuleState = iota
	stateBuy
	stateSell
)

type PlayerService interface {
	GetMoney() int
	AddMoney(amount int)
	AddResourceQuantity(name string, quantity int)
	IsWarehouseEmpty() bool
}

type ToastMessenger interface {
	Message(text string, a ...any)
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) View() tea.View {

	panel := tabs.NewTabPanel("Buy", "Sell", "Unlock Resource")

	panel.
		WriteLn(viewMarket(m.Market)).
		AddLn().
		WriteLn(viewWarehouse(m.Resources))

	if m.state == stateBuy {
		viewBuy := m.screenBuy.View()
		layerBuyFlow := lipgloss.NewLayer(viewBuy.Content).X(15).Y(3).Z(1)
		panel.AddLayer(layerBuyFlow)
	} else if m.state == stateSell {
		viewSell := m.screenSell.View()
		layerSellFlow := lipgloss.NewLayer(viewSell.Content).X(15).Y(3).Z(1)
		panel.AddLayer(layerSellFlow)
	}

	return tea.NewView(panel.Render())
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmds []tea.Cmd

	switch m.state {
	case stateBuy:
		_, cmd := m.screenBuy.Update(msg)
		cmds = append(cmds, cmd)
	case stateSell:
		_, cmd := m.screenSell.Update(msg)
		cmds = append(cmds, cmd)
	case stateList:
		switch msg := msg.(type) {
		case tea.KeyPressMsg:
			switch msg.String() {
			case "b", "B":
				m.state = stateBuy
				cmd := m.initBuyScreen()
				cmds = append(cmds, cmd)
			case "s", "S":
				if m.playerService.IsWarehouseEmpty() {
					m.toastMessenger.Message("Warehouse is empty. Nothing to sell.")
				} else {
					m.state = stateSell
					cmd := m.initSellScreen()
					cmds = append(cmds, cmd)
				}
			case "u", "U":
				m.unlockItem()
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func viewMarket(m *Market) string {
	var view strings.Builder
	view.WriteString("Market:\n")

	itemNames := slices.Sorted(maps.Keys(m.Resources))

	for _, name := range itemNames {
		item := m.Resources[name]

		priceChange := item.PriceCurrent - m.GetPricePrevious(name)

		var changeStyle lipgloss.Style
		var changeDisplay string
		if priceChange > 0 {
			changeDisplay = "↑ " + strconv.Itoa(priceChange)
			changeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#8FBC8B"))
		} else if priceChange < 0 {
			changeDisplay = "↓ " + strconv.Itoa(priceChange)
			changeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#CD5C5C"))
		} else {
			changeDisplay = "± 0"
			changeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#D3D3D3"))
		}

		view.WriteString(fmt.Sprintf(
			" %v: %v (%v)\n",
			name,
			util.FormatMoney(item.PriceCurrent),
			changeStyle.Render(changeDisplay),
		))
	}
	return view.String()
}

func viewWarehouse(resources map[string]int) string {
	var view strings.Builder
	view.WriteString("Your Warehouse:\n")
	hasItems := false
	items := slices.Sorted(maps.Keys(resources))
	for _, item := range items {
		count := resources[item]
		if count > 0 {
			hasItems = true
			view.WriteString(fmt.Sprintf(" - %v: %v\n", item, count))
		}
	}

	if !hasItems {
		view.WriteString("   Empty\n")
	}

	return view.String()
}

func (m *Model) initBuyScreen() tea.Cmd {
	buyScreen := NewBuyScreen(m.Market.GetPricesCurrent(), m.playerService.GetMoney())

	buyScreen.OnAborted = func() {
		m.state = stateList
	}

	buyScreen.OnComplete = func(purchaseOrder map[string]int) {
		m.state = stateList

		for item, quantity := range purchaseOrder {
			if quantity == 0 {
				continue
			}

			cost := m.Market.Resources[item].PriceCurrent * quantity
			m.playerService.AddMoney(-cost)
			m.playerService.AddResourceQuantity(item, quantity)
		}
	}

	m.screenBuy = buyScreen

	return m.screenBuy.Init()
}

func (m *Model) initSellScreen() tea.Cmd {
	sellScreen := NewSellScreen(m.Resources, m.Market.GetPricesCurrent())

	sellScreen.OnAborted = func() {
		m.state = stateList
	}

	sellScreen.OnComplete = func(sellOrder map[string]int) {
		m.state = stateList

		for item, quantity := range sellOrder {
			if quantity == 0 {
				continue
			}

			cost := m.Market.Resources[item].PriceCurrent * quantity
			m.playerService.AddMoney(cost)
			m.playerService.AddResourceQuantity(item, -quantity)
		}
	}

	m.screenSell = sellScreen

	return m.screenSell.Init()
}

func (m *Model) unlockItem() {

	if len(m.Market.LockedResources) == 0 {
		m.toastMessenger.Message("You have already unlocked all resources.")
		return
	}

	if m.playerService.GetMoney() < m.Market.UnlockCost {
		m.toastMessenger.Message("You don't have enough to unlock a new item. Stop being poor!")
		return
	}

	m.playerService.AddMoney(-m.Market.UnlockCost)
	unlockedItem := m.Market.UnlockResource()
	if unlockedItem != nil {
		m.playerService.AddResourceQuantity(unlockedItem.Name, 0)
		m.toastMessenger.Message("New resource permit secured: %v", unlockedItem.Name)
	}
}
