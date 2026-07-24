package internal

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/randomvlad/trader-vlads/internal/component/status"
	"github.com/randomvlad/trader-vlads/internal/component/toast"
	"github.com/randomvlad/trader-vlads/internal/screen"
	"github.com/randomvlad/trader-vlads/internal/util"
)

func NewGame() *GameData {
	market := NewMarket()
	player := NewPlayer(market)

	return &GameData{
		player:     player,
		market:     market,
		showScreen: ScreenMain,
		toast:      toast.Model{},
		status:     status.New(),
	}
}

func (gd *GameData) Init() tea.Cmd {
	return nil
}

func (gd *GameData) View() tea.View {
	var sbContent strings.Builder

	sbContent.WriteString(gd.status.Render(gd.market.week, gd.player.money))

	stylePanel := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#04B575")).
		Height(15).
		Padding(0, 2).
		Width(50).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#04B575"))

	viewInventory := stylePanel.Render(getViewInventory(gd.player))
	viewMarket := stylePanel.Render(getViewMarket(gd.market))

	sbContent.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, viewInventory, viewMarket) + "\n")

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

	return tea.NewView(compositor.Render())
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
	buyScreen := screen.NewBuyScreen(g.market.GetCurrentPrices(), g.player.money)

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
	sellScreen := screen.NewSellScreen(g.player.inventory, g.market.GetCurrentPrices())

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
	view := "Inventory:\n"
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
	view := "Local Market:\n"

	prices := m.GetCurrentPrices()
	items := slices.Sorted(maps.Keys(m.items))

	for _, item := range items {
		view += fmt.Sprintf(" %v: %v \n", item, util.FormatMoney(prices[item]))
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
	unlockedItem := gd.market.lockedItems[0]
	// TODO: needs to be implemented in market.go
	// player.lockedItems = slices.Delete(gd.market.lockedItems, 0, 1)
	player.inventory[unlockedItem.name] = 0

	gd.toast.SetMessage("New guild permit secured: %v", unlockedItem)
}

func getActionsBar(player *Player) string {

	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575"))
	styleFirstLetter := style.Bold(true).Underline(true)

	var sbContent strings.Builder
	sbContent.WriteString(style.Render("Actions: "))

	for index, action := range player.actions {
		styledAction := lipgloss.StyleRanges(action,
			lipgloss.NewRange(0, 1, styleFirstLetter),
			lipgloss.NewRange(1, len(action), style),
		)

		sbContent.WriteString(styledAction)

		if index < len(player.actions)-1 {
			sbContent.WriteString(style.Render(" • "))
		}
	}

	styleActionsBar := lipgloss.NewStyle().
		Width(100).
		Padding(0, 2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#04B575"))

	return styleActionsBar.Render(sbContent.String())
}
