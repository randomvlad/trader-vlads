package market

import (
	"maps"
	"slices"
	"strconv"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/randomvlad/trader-vlads/internal/appmod/keybind"
	"github.com/randomvlad/trader-vlads/internal/appmod/keybind/press"
	"github.com/randomvlad/trader-vlads/internal/component/tabs"
	"github.com/randomvlad/trader-vlads/internal/util"
	"github.com/randomvlad/trader-vlads/internal/util/stringutil"
)

type Model struct {
	market         *Market
	playerService  PlayerService
	keyBinder      *keybind.KeyBinder
	toastMessenger ToastMessenger
	Resources      map[string]int
	state          marketModuleState
	screenBuy      *BuyScreen
	screenSell     *SellScreen
}

func NewTuiModel(market *Market, player PlayerService, keyBinder *keybind.KeyBinder, toast ToastMessenger) *Model {
	return &Model{
		market:         market,
		playerService:  player,
		keyBinder:      keyBinder,
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

	// TODO issue: app state/location needs to account for buy/sell forms.
	// Temp workaround: handle with if statement in each action func

	actionBuyStart := press.NewActionBuilder().
		Name("Buy").
		ActionCmd(func() tea.Cmd {
			if m.state != stateList {
				return nil
			}

			m.state = stateBuy
			return m.initBuyScreen()
		}).
		Permanent().
		SortOrder(1).
		Build()

	actionSellStart := press.NewActionBuilder().
		Name("Sell").
		ActionCmd(func() tea.Cmd {
			if m.state != stateList {
				return nil
			}

			if m.playerService.IsWarehouseEmpty() {
				m.toastMessenger.Message("Warehouse is empty. Nothing to sell.")
				return nil
			} else {
				m.state = stateSell
				return m.initSellScreen()
			}
		}).
		Permanent().
		SortOrder(2).
		Build()

	actionUnlock := press.NewActionBuilder().
		Name("Unlock Resource").
		Action(func() {
			if m.state != stateList {
				return
			}

			m.unlockItem()
		}).
		Permanent().
		SortOrder(3).
		Build()

	m.keyBinder.AddActions("tui-market", actionBuyStart, actionSellStart, actionUnlock)

	return nil
}

func (m *Model) View() tea.View {

	panel := tabs.NewTabPanel().
		WriteLn(viewMarket(m.market)).
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

	return panel.RenderTeaView()
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

	}

	return m, tea.Batch(cmds...)
}

func viewMarket(m *Market) string {
	view := stringutil.NewBuilder().WriteLn("Market:")

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

		view.Writef(
			" %v: %v (%v)\n",
			name,
			util.FormatCurrency(item.PriceCurrent),
			changeStyle.Render(changeDisplay),
		)
	}
	return view.String()
}

func viewWarehouse(resources map[string]int) string {
	view := stringutil.NewBuilder().WriteLn("Your Warehouse:")
	hasItems := false
	items := slices.Sorted(maps.Keys(resources))
	for _, item := range items {
		count := resources[item]
		if count > 0 {
			hasItems = true
			view.Writef(" - %v: %v\n", item, count)
		}
	}

	if !hasItems {
		view.WriteLn("   Empty")
	}

	return view.String()
}

func (m *Model) initBuyScreen() tea.Cmd {
	buyScreen := NewBuyScreen(m.market.GetPricesCurrent(), m.playerService.GetMoney())

	buyScreen.OnAborted = func() {
		m.state = stateList
	}

	buyScreen.OnComplete = func(purchaseOrder map[string]int) {
		m.state = stateList

		for item, quantity := range purchaseOrder {
			if quantity == 0 {
				continue
			}

			cost := m.market.Resources[item].PriceCurrent * quantity
			m.playerService.AddMoney(-cost)
			m.playerService.AddResourceQuantity(item, quantity)
		}
	}

	m.screenBuy = buyScreen

	return m.screenBuy.Init()
}

func (m *Model) initSellScreen() tea.Cmd {
	sellScreen := NewSellScreen(m.Resources, m.market.GetPricesCurrent())

	sellScreen.OnAborted = func() {
		m.state = stateList
	}

	sellScreen.OnComplete = func(sellOrder map[string]int) {
		m.state = stateList

		for item, quantity := range sellOrder {
			if quantity == 0 {
				continue
			}

			cost := m.market.Resources[item].PriceCurrent * quantity
			m.playerService.AddMoney(cost)
			m.playerService.AddResourceQuantity(item, -quantity)
		}
	}

	m.screenSell = sellScreen

	return m.screenSell.Init()
}

func (m *Model) unlockItem() {

	if len(m.market.LockedResources) == 0 {
		m.toastMessenger.Message("You have already unlocked all resources.")
		return
	}

	if m.playerService.GetMoney() < m.market.UnlockCost {
		m.toastMessenger.Message("You don't have enough to unlock a new item. Stop being poor!")
		return
	}

	m.playerService.AddMoney(-m.market.UnlockCost)
	unlockedItem := m.market.UnlockResource()
	if unlockedItem != nil {
		m.playerService.AddResourceQuantity(unlockedItem.Name, 0)
		m.toastMessenger.Message("New resource permit secured: %v", unlockedItem.Name)
	}
}
