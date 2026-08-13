package market

// TODO: rename file ... it's more than just a view

import (
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/randomvlad/trader-vlads/internal/appstyle"
	"github.com/randomvlad/trader-vlads/internal/util"
)

type MarketModule struct {
	Market         *Market
	PlayerService  PlayerService
	ToastMessenger ToastMessenger
	Resources      map[string]int
	Operation      int // create enum? rename to state maybe?
	screenBuy      *BuyScreen
	screenSell     *SellScreen
}

type PlayerService interface {
	GetMoney() int
	AddMoney(amount int)
	AddResourceQuantity(name string, quantity int)
	IsWarehouseEmpty() bool
}

type ToastMessenger interface {
	Message(text string, a ...any)
}

func (m *MarketModule) Init() tea.Cmd {
	return nil
}

func (m *MarketModule) View() tea.View {

	var sbContent strings.Builder
	sbContent.WriteString(viewMarket(m.Market))
	sbContent.WriteString("\n\n")
	sbContent.WriteString(viewWarehouse(m.Resources))

	layerMain := lipgloss.NewLayer(sbContent.String())

	compositor := lipgloss.NewCompositor(layerMain)

	if m.Operation == 1 /* ScreenBuy */ {
		viewBuy := m.screenBuy.View()
		layerBuyFlow := lipgloss.NewLayer(viewBuy.Content).X(15).Y(3).Z(1)
		compositor.AddLayers(layerBuyFlow)
	} else if m.Operation == 2 /* ScreenSell */ {
		viewSell := m.screenSell.View()
		layerSellFlow := lipgloss.NewLayer(viewSell.Content).X(15).Y(3).Z(1)
		compositor.AddLayers(layerSellFlow)
	}

	return tea.NewView(appstyle.StyleAppContainer.Render(compositor.Render()))
}

func (m *MarketModule) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmds []tea.Cmd

	if m.Operation == 1 /* ScreenBuy */ {
		_, cmd := m.screenBuy.Update(msg)
		cmds = append(cmds, cmd)
	} else if m.Operation == 2 /* ScreenSell */ {

		_, cmd := m.screenSell.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		switch msg := msg.(type) {
		case tea.KeyPressMsg:
			switch msg.String() {
			case "b", "B":
				m.Operation = 1
				cmd := m.initBuyScreen()
				cmds = append(cmds, cmd)
			case "s", "S":
				if m.PlayerService.IsWarehouseEmpty() {
					m.ToastMessenger.Message("Warehouse is empty. Nothing to sell.")
				} else {
					m.Operation = 2
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
	view := "Market:\n"

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

		view += fmt.Sprintf(
			" %v: %v (%v)\n",
			name,
			util.FormatMoney(item.PriceCurrent),
			changeStyle.Render(changeDisplay),
		)
	}
	return view
}

func viewWarehouse(resources map[string]int) string {
	view := "Your Warehouse:\n"
	hasItems := false
	items := slices.Sorted(maps.Keys(resources))
	for _, item := range items {
		count := resources[item]
		if count > 0 {
			hasItems = true
			view += fmt.Sprintf(" - %v: %v\n", item, count)
		}
	}

	if !hasItems {
		view += "   Empty\n"
	}

	return view
}

func (m *MarketModule) initBuyScreen() tea.Cmd {
	buyScreen := NewBuyScreen(m.Market.GetPricesCurrent(), m.PlayerService.GetMoney())

	buyScreen.OnAborted = func() {
		m.Operation = 0
	}

	buyScreen.OnComplete = func(purchaseOrder map[string]int) {
		m.Operation = 0

		for item, quantity := range purchaseOrder {
			if quantity == 0 {
				continue
			}

			cost := m.Market.Resources[item].PriceCurrent * quantity
			m.PlayerService.AddMoney(-cost)
			m.PlayerService.AddResourceQuantity(item, quantity)
		}
	}

	m.screenBuy = buyScreen

	return m.screenBuy.Init()
}

func (m *MarketModule) initSellScreen() tea.Cmd {
	sellScreen := NewSellScreen(m.Resources, m.Market.GetPricesCurrent())

	sellScreen.OnAborted = func() {
		m.Operation = 0
	}

	sellScreen.OnComplete = func(sellOrder map[string]int) {
		m.Operation = 0

		for item, quantity := range sellOrder {
			if quantity == 0 {
				continue
			}

			cost := m.Market.Resources[item].PriceCurrent * quantity
			m.PlayerService.AddMoney(cost)
			m.PlayerService.AddResourceQuantity(item, -quantity)
		}
	}

	m.screenSell = sellScreen

	return m.screenSell.Init()
}

func (m *MarketModule) unlockItem() {

	if len(m.Market.LockedResources) == 0 {
		m.ToastMessenger.Message("You have already unlocked all resources.")
		return
	}

	if m.PlayerService.GetMoney() < m.Market.UnlockCost {
		m.ToastMessenger.Message("You don't have enough to unlock a new item. Stop being poor!")
		return
	}

	m.PlayerService.AddMoney(-m.Market.UnlockCost)
	unlockedItem := m.Market.UnlockResource()
	if unlockedItem != nil {
		m.PlayerService.AddResourceQuantity(unlockedItem.Name, 0)
		m.ToastMessenger.Message("New resource permit secured: %v", unlockedItem.Name)
	}
}
