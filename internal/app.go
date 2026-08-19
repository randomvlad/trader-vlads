package internal

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	appeq "github.com/randomvlad/trader-vlads/internal/appmod/equipment"
	appmarket "github.com/randomvlad/trader-vlads/internal/appmod/market"
	appstats "github.com/randomvlad/trader-vlads/internal/appmod/stats"
	"github.com/randomvlad/trader-vlads/internal/appstyle"
	"github.com/randomvlad/trader-vlads/internal/component/actionfooter"
	"github.com/randomvlad/trader-vlads/internal/component/status"
	"github.com/randomvlad/trader-vlads/internal/component/tabs"
	toastcmp "github.com/randomvlad/trader-vlads/internal/component/toast"
	"github.com/randomvlad/trader-vlads/internal/util"
)

type GameData struct {
	player       *Player
	marketModel  *appmarket.Model
	eqModel      *appeq.Model
	statsModel   *appstats.Model
	eventTrack   *EventTracker
	tabs         *tabs.Model
	actionFooter *actionfooter.Model
	toast        *toastcmp.Toast
	status       status.Model
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
	player := NewPlayer(market, random)
	toast := &toastcmp.Toast{}

	tabNames := []string{"📜 Events", "🏦 Market", "💠 Equipment", "🔍 Stats"}

	return &GameData{
		player:       player,
		eqModel:      appeq.NewTuiModel(player, toast),
		marketModel:  appmarket.NewTuiModel(market, player, toast),
		statsModel:   appstats.NewTuiModel(player),
		eventTrack:   NewEventTracker(random),
		tabs:         tabs.NewModel(tabNames, appstyle.AppWidth),
		actionFooter: actionfooter.NewModel(actionfooter.FooterStandalone, "Next Week", "Quit"),
		toast:        toast,
		status:       status.New(),
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
	var view strings.Builder

	view.WriteString(gd.status.Render(gd.marketModel.Market.Week, gd.player.money))

	view.WriteString(gd.tabs.View().Content + "\n")

	activeTab := TabId(gd.tabs.ActiveTab)
	switch activeTab {
	case TabEvents:
		activeTabContent := "Events History"
		view.WriteString(appstyle.StyleTabView.Render(activeTabContent) + "\n")
	case TabMarket:
		gd.marketModel.Resources = gd.player.warehouse.resources
		view.WriteString(gd.marketModel.View().Content)
		view.WriteString("\n")
	case TabEquipment:
		view.WriteString(gd.eqModel.View().Content)
		view.WriteString("\n")
	case TabStats:
		view.WriteString(gd.statsModel.View().Content)
		view.WriteString("\n")
	}

	view.WriteString(gd.actionFooter.Render())

	layerMain := lipgloss.NewLayer(view.String())

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

			var toastMessage string
			event := gd.eventTrack.getRandomEvent()
			if event != nil {
				toastMessage = event.name + "\n\n" + event.description + "\n"
				gd.player.money += event.money
			}

			expiredEffects := gd.player.NextWeek()
			for _, effect := range expiredEffects {
				toastMessage += effect.GetExpiredMessage() + "\n"
			}

			if toastMessage != "" {
				gd.toast.Message(toastMessage)
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
		switch TabId(gd.tabs.ActiveTab) {
		case TabMarket:
			_, cmd := gd.marketModel.Update(msg)
			cmds = append(cmds, cmd)
		case TabEquipment:
			_, cmd := gd.eqModel.Update(msg)
			cmds = append(cmds, cmd)
		case TabStats:
			_, cmd := gd.statsModel.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return gd, tea.Batch(cmds...)
}
