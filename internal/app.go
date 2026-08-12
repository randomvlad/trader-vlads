package internal

import (
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/randomvlad/trader-vlads/internal/component/status"
	"github.com/randomvlad/trader-vlads/internal/component/tabs"
	"github.com/randomvlad/trader-vlads/internal/component/toast"
	"github.com/randomvlad/trader-vlads/internal/screen"
	"github.com/randomvlad/trader-vlads/internal/util"
)

func NewGame() *GameData {
	random := util.NewRandomGenerator(nil)

	market := NewMarket(random)
	player := NewPlayer(market)

	return &GameData{
		player:     player,
		market:     market,
		eventTrack: NewEventTracker(random),
		showScreen: ScreenMain,
		tabs:       tabs.NewTabsModel(screen.AppWidth),
		toast:      toast.Model{},
		status:     status.New(),
	}
}

func (gd *GameData) Init() tea.Cmd {

	// TODO: think more about how to structure Init()'s
	gd.tabs.Init()

	return nil
}

func (gd *GameData) View() tea.View {
	var sbContent strings.Builder

	sbContent.WriteString(gd.status.Render(gd.market.week, gd.player.money))

	sbContent.WriteString(gd.tabs.View().Content + "\n")

	var activeTabContent string
	switch gd.tabs.ActiveTab {
	case 0:
		activeTabContent = "Events History"
	case 1:
		activeTabContent = getViewMarket(gd.market) + "\n\n" + getViewInventory(gd.player)
	case 2:
		activeTabContent = gd.player.ViewEquipment()
	case 3:
		activeTabContent = "Stats"
	}

	sbContent.WriteString(screen.StyleTabView.Render(activeTabContent) + "\n")

	sbContent.WriteString(getActionsBar(gd.player) + "\n")

	layerMain := lipgloss.NewLayer(sbContent.String())

	compositor := lipgloss.NewCompositor(layerMain)

	if gd.toast.Show {
		layerToast := lipgloss.NewLayer(gd.toast.Render()).X(25).Y(8).Z(1)
		compositor.AddLayers(layerToast)
	}

	if gd.showScreen == ScreenBuy {
		viewBuy := gd.screenBuy.View()
		layerBuyFlow := lipgloss.NewLayer(viewBuy.Content).X(15).Y(3).Z(1)
		compositor.AddLayers(layerBuyFlow)
	} else if gd.showScreen == ScreenSell {
		viewSell := gd.screenSell.View()
		layerSellFlow := lipgloss.NewLayer(viewSell.Content).X(15).Y(3).Z(1)
		compositor.AddLayers(layerSellFlow)
	}

	return tea.NewView(screen.StyleAppContainer.Render(compositor.Render()))
}

func (gd *GameData) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmds []tea.Cmd

	if gd.showScreen == ScreenBuy {
		_, cmd := gd.screenBuy.Update(msg)
		cmds = append(cmds, cmd)
	} else if gd.showScreen == ScreenSell {

		_, cmd := gd.screenSell.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		switch msg := msg.(type) {
		case tea.KeyPressMsg:
			switch msg.String() {
			case "n", "N":
				gd.showScreen = ScreenMain
				gd.market.NextWeek()
				gd.toast.Clear()

				event := gd.eventTrack.getRandomEvent()
				if event != nil {
					gd.toast.SetMessage(event.name + "\n\n" + event.description)
					gd.player.money += event.money
				}
			case "b", "B":
				gd.showScreen = ScreenBuy
				cmd := gd.initBuyScreen()
				cmds = append(cmds, cmd)
			case "s", "S":
				if gd.player.isInventoryEmpty() {
					gd.toast.SetMessage("Inventory is empty. Nothing to sell.")
				} else {
					gd.showScreen = ScreenSell
					cmd := gd.initSellScreen()
					cmds = append(cmds, cmd)
				}
			case "u", "U":
				gd.showScreen = ScreenMain
				unlockItem(gd)
			case "left", "right":
				_, cmd := gd.tabs.Update(msg)
				cmds = append(cmds, cmd)
			case "enter", "esc":
				gd.toast.Clear()
			case "q", "ctrl+c":
				gd.toast.SetMessage("Farewell and safe travels!")
				return gd, tea.Quit
			default:
				gd.toast.SetMessage("Unknown action: %v", msg.String())
			}
		}
	}

	return gd, tea.Batch(cmds...)
}

func (g *GameData) initBuyScreen() tea.Cmd {
	buyScreen := screen.NewBuyScreen(g.market.GetPricesCurrent(), g.player.money)

	buyScreen.OnAborted = func() {
		g.showScreen = ScreenMain
	}

	buyScreen.OnComplete = func(purchaseOrder map[string]int) {
		g.showScreen = ScreenMain

		for item, quantity := range purchaseOrder {
			if quantity == 0 {
				continue
			}

			cost := g.market.items[item].priceCurrent * quantity
			g.player.money -= cost
			g.player.inventory[item] += quantity
		}
	}

	g.screenBuy = buyScreen

	return g.screenBuy.Init()
}

func (g *GameData) initSellScreen() tea.Cmd {
	sellScreen := screen.NewSellScreen(g.player.inventory, g.market.GetPricesCurrent())

	sellScreen.OnAborted = func() {
		g.showScreen = ScreenMain
	}

	sellScreen.OnComplete = func(sellOrder map[string]int) {
		g.showScreen = ScreenMain

		for item, quantity := range sellOrder {
			if quantity == 0 {
				continue
			}

			cost := g.market.items[item].priceCurrent * quantity
			g.player.money += cost
			g.player.inventory[item] -= quantity
		}
	}

	g.screenSell = sellScreen

	return g.screenSell.Init()
}

func getViewInventory(player *Player) string {
	view := "Inventory:\n" // TODO: rename to ware house?
	hasItems := false
	items := slices.Sorted(maps.Keys(player.inventory))
	for _, item := range items {
		count := player.inventory[item]
		if count > 0 {
			hasItems = true
			view += fmt.Sprintf(" - %v: %v\n", item, count)
		}
	}

	if !hasItems {
		view += " (Empty)\n"
	}

	return view
}

func getViewMarket(m *Market) string {
	view := "Market:\n"

	itemNames := slices.Sorted(maps.Keys(m.items))

	for _, name := range itemNames {
		item := m.items[name]

		priceChange := item.priceCurrent - m.GetPricePrevious(name)

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
			util.FormatMoney(item.priceCurrent),
			changeStyle.Render(changeDisplay),
		)
	}
	return view
}

func unlockItem(gd *GameData) {
	player := gd.player
	if len(gd.market.lockedItems) == 0 {
		gd.toast.SetMessage("You have already unlocked all items.")
		return
	}

	if player.money < gd.market.unlockCost {
		gd.toast.SetMessage("You don't have enough to unlock a new item. Stop being poor!")
		return
	}

	player.money -= gd.market.unlockCost
	unlockedItem := gd.market.UnlockItem()
	if unlockedItem != nil {
		player.inventory[unlockedItem.name] = 0
		gd.toast.SetMessage("New guild permit secured: %v", unlockedItem.name)
	}
}

func getActionsBar(player *Player) string {

	var sbContent strings.Builder
	sbContent.WriteString("Actions: ")

	for index, action := range player.actions {
		styledAction := lipgloss.StyleRanges(
			action,
			lipgloss.NewRange(0, 1, screen.StyleTextFirstLetter),
			lipgloss.NewRange(1, len(action), screen.NewAppStyle()),
		)

		sbContent.WriteString(styledAction)

		if index < len(player.actions)-1 {
			sbContent.WriteString(" • ")
		}
	}

	return screen.StyleActionsBar.Render(sbContent.String())
}
