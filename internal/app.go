package internal

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	appeq "github.com/randomvlad/trader-vlads/internal/appmod/equipment"
	appmarket "github.com/randomvlad/trader-vlads/internal/appmod/market"
	"github.com/randomvlad/trader-vlads/internal/appstyle"
	"github.com/randomvlad/trader-vlads/internal/component/status"
	"github.com/randomvlad/trader-vlads/internal/component/tabs"
	toastcmp "github.com/randomvlad/trader-vlads/internal/component/toast"
	"github.com/randomvlad/trader-vlads/internal/util"
)

type GameData struct {
	player      *Player
	marketModel *appmarket.Model
	eqModel     *appeq.Model
	eventTrack  *EventTracker
	tabs        *tabs.TabsModel
	toast       *toastcmp.Toast
	status      status.Model
}

type TabId int

const (
	TabEvents TabId = iota
	TabMarket
	TabEquipment
	TabStats
)

func NewGame() *GameData {
	random := util.NewRandomGenerator(nil)

	market := appmarket.NewMarket(random)
	player := NewPlayer(market)
	toast := &toastcmp.Toast{}

	tabNames := []string{"📜 Events", "🏦 Market", "💠 Equipment", "🔍 Stats"}

	return &GameData{
		player:      player,
		eqModel:     appeq.NewTuiModel(player),
		marketModel: appmarket.NewTuiModel(market, player, toast),
		eventTrack:  NewEventTracker(random),
		tabs:        tabs.NewTabsModel(tabNames, appstyle.AppWidth),
		toast:       toast,
		status:      status.New(),
	}
}

func (gd *GameData) Init() tea.Cmd {
	var cmds []tea.Cmd

	cmdMarket := gd.marketModel.Init()
	cmds = append(cmds, cmdMarket)

	cmdTabs := gd.tabs.Init()
	cmds = append(cmds, cmdTabs)

	return tea.Batch(cmds...)
}

func (gd *GameData) View() tea.View {
	var sbContent strings.Builder

	sbContent.WriteString(gd.status.Render(gd.marketModel.Market.Week, gd.player.money))

	sbContent.WriteString(gd.tabs.View().Content + "\n")

	var activeTabContent string
	switch TabId(gd.tabs.ActiveTab) {
	case TabEvents:
		activeTabContent = "Events History"
	case TabMarket:
		gd.marketModel.Resources = gd.player.warehouse.resources
		activeTabContent = gd.marketModel.View().Content
	case TabEquipment:
		activeTabContent = gd.eqModel.View().Content
	case TabStats:
		activeTabContent = "Stats"
	}

	sbContent.WriteString(appstyle.StyleTabView.Render(activeTabContent) + "\n")

	sbContent.WriteString(getActionsBar() + "\n")

	layerMain := lipgloss.NewLayer(sbContent.String())

	compositor := lipgloss.NewCompositor(layerMain)

	if gd.toast.Show {
		layerToast := lipgloss.NewLayer(gd.toast.Render()).X(25).Y(8).Z(1)
		compositor.AddLayers(layerToast)
	}

	return tea.NewView(appstyle.StyleAppContainer.Render(compositor.Render()))
}

func (gd *GameData) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmds []tea.Cmd

	globalKeyPress := false

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "n", "N":
			gd.marketModel.Market.NextWeek()
			gd.toast.Clear()

			event := gd.eventTrack.getRandomEvent()
			if event != nil {
				gd.toast.Message(event.name + "\n\n" + event.description)
				gd.player.money += event.money
			}
			globalKeyPress = true
		case "left", "right":
			_, cmd := gd.tabs.Update(msg)
			cmds = append(cmds, cmd)
			globalKeyPress = true
		case "enter", "esc":
			gd.toast.Clear()
			// TODO: Key Press & App State → delegate to corresponding Model Update()
			// need exact app state to give key presses the right context. For example: if user is on market AND buy/sell popup, then enter has different behavior
		case "q", "ctrl+c":
			gd.toast.Message("Farewell and safe travels!")
			return gd, tea.Quit
		}
	}

	if !globalKeyPress {
		if TabId(gd.tabs.ActiveTab) == TabMarket {
			_, cmd := gd.marketModel.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return gd, tea.Batch(cmds...)
}

func getActionsBar() string {

	var sbContent strings.Builder
	sbContent.WriteString("Actions: ")

	actions := []string{
		"Buy",
		"Sell",
		"Next Week",
		"Unlock Item",
		"Quit",
	}

	for index, action := range actions {
		styledAction := lipgloss.StyleRanges(
			action,
			lipgloss.NewRange(0, 1, appstyle.StyleTextFirstLetter),
			lipgloss.NewRange(1, len(action), appstyle.NewAppStyle()),
		)

		sbContent.WriteString(styledAction)

		if index < len(actions)-1 {
			sbContent.WriteString(" • ")
		}
	}

	return appstyle.StyleActionsBar.Render(sbContent.String())
}
