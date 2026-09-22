package internal

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	eq "github.com/randomvlad/trader-vlads/internal/appmod/equipment"
	ev "github.com/randomvlad/trader-vlads/internal/appmod/event"
	"github.com/randomvlad/trader-vlads/internal/appmod/keybind"
	"github.com/randomvlad/trader-vlads/internal/appmod/keybind/press"
	appmarket "github.com/randomvlad/trader-vlads/internal/appmod/market"
	p "github.com/randomvlad/trader-vlads/internal/appmod/player"
	appstats "github.com/randomvlad/trader-vlads/internal/appmod/stats"
	"github.com/randomvlad/trader-vlads/internal/appstyle"
	"github.com/randomvlad/trader-vlads/internal/component/actionfooter"
	apppanel "github.com/randomvlad/trader-vlads/internal/component/panel"
	"github.com/randomvlad/trader-vlads/internal/component/status"
	"github.com/randomvlad/trader-vlads/internal/component/tabs"
	toastcmp "github.com/randomvlad/trader-vlads/internal/component/toast"
	"github.com/randomvlad/trader-vlads/internal/util"
	"github.com/randomvlad/trader-vlads/internal/util/stringutil"
)

type GameData struct {
	turnKeeper   *ev.TurnKeeper
	player       *p.Player
	marketModel  *appmarket.Model
	eqModel      *eq.Model
	statsModel   *appstats.Model
	eventTrack   *ev.EventTracker
	keyBinder    *keybind.KeyBinder
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
	keyBinder := keybind.NewKeyBinder()
	random := util.NewRandomGenerator(nil)

	market := appmarket.NewMarket(random)
	player := p.NewPlayer(market, random)
	toast := &toastcmp.Toast{}
	turnKeeper := ev.NewTurnKeeper(player, market, keyBinder, random, toast)

	return &GameData{
		player:      player,
		turnKeeper:  turnKeeper,
		keyBinder:   keyBinder,
		eqModel:     eq.NewTuiModel(player, keyBinder, toast),
		marketModel: appmarket.NewTuiModel(market, player, keyBinder, toast),
		statsModel:  appstats.NewTuiModel(player),
		tabs:        tabs.NewModel("📜 Events", "🏦 Market", "💠 Equipment", "🔍 Stats"),
		actionFooter: actionfooter.NewModel(actionfooter.FooterPanel).
			WithKeyBinder(keyBinder).
			WithDisplayRightActionContext("tui-control"),
		toast:  toast,
		status: status.New(),
	}
}

func (gd *GameData) Init() tea.Cmd {
	gd.bindActions()

	cmdEq := gd.eqModel.Init()
	cmdMarket := gd.marketModel.Init()

	return tea.Batch(cmdEq, cmdMarket)
}

func (gd *GameData) bindActions() {

	nextWeek := press.NewActionBuilder().
		Name("End Turn").
		Action(func() { gd.turnKeeper.Next() }).
		Permanent().
		SortOrder(1).
		Build()

	quit := press.NewActionBuilder().
		Name("Quit").
		ActionCmd(func() tea.Cmd {
			gd.toast.Message("Farewell and safe travels!")
			return tea.Quit
		}).
		Permanent().
		SortOrder(2).
		Build()

	tabLeft := press.NewActionBuilder().
		NameKeyPress("NavTabLeft", "left").
		Action(func() { gd.tabs.SelectLeft() }).
		Permanent().
		FooterHidden().
		Build()

	tabRight := press.NewActionBuilder().
		NameKeyPress("NavTabRight", "right").
		Action(func() { gd.tabs.SelectRight() }).
		Permanent().
		FooterHidden().
		Build()

	clearToast := press.NewActionBuilder().
		NameKeyPress("ClearToast", "esc").
		Action(func() { gd.toast.Clear() }).
		FooterHidden().
		Permanent().
		Build()

	gd.keyBinder.AddActions("tui-control", tabLeft, tabRight, clearToast, nextWeek, quit)
}

func (gd *GameData) View() tea.View {
	var view stringutil.Builder

	if activeStory := gd.turnKeeper.EventTracker.GetActiveStory(); activeStory != nil {
		panel := apppanel.NewModel().
			WithTitle(activeStory.GetName()).
			WithFooter(activeStory.GetCurrentActions()...).
			Write(activeStory.View())
		view.WriteStyle(panel.Render(), appstyle.StyleStoryContainer)

	} else {
		// status bar
		view.WriteLn(gd.status.Render(gd.turnKeeper.GetTurn(), gd.player.GetMoney()))

		// tabs and tab content
		view.WriteLn(gd.tabs.View())

		var actionsContext string
		activeTab := TabId(gd.tabs.ActiveTab)
		switch activeTab {
		case TabEvents:
			panel := tabs.NewTabPanel().
				WithBodyBorderFooterCompatible().
				WriteLn("Events History")
			view.WriteLn(panel.Render())
		case TabMarket:
			gd.marketModel.Resources = gd.player.Warehouse.Resources
			view.WriteLn(gd.marketModel.View().Content)
			actionsContext = "tui-market"
		case TabEquipment:
			view.WriteLn(gd.eqModel.View().Content)
			actionsContext = "tui-eq"
		case TabStats:
			view.WriteLn(gd.statsModel.View().Content)
			actionsContext = "tui-stats"
		}

		// footer
		view.Write(gd.actionFooter.Render(actionsContext))
	}

	layerMain := lipgloss.NewLayer(view.String())

	compositor := lipgloss.NewCompositor(layerMain)

	if gd.toast.Show {
		layerToast := lipgloss.NewLayer(gd.toast.Render()).X(25).Y(8).Z(2)
		compositor.AddLayers(layerToast)
	}

	return tea.NewView(appstyle.StyleAppContainer.Render(compositor.Render()))
}

func (gd *GameData) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmds []tea.Cmd

	if story := gd.turnKeeper.EventTracker.GetActiveStory(); story != nil {
		if gd.keyBinder.IsBound(story.GetStateId(), msg) {
			cmd := gd.keyBinder.Execute(story.GetStateId(), msg)
			cmds = append(cmds, cmd)
		}
	} else if gd.keyBinder.IsBound("tui-control", msg) {
		cmd := gd.keyBinder.Execute("tui-control", msg)
		cmds = append(cmds, cmd)
	} else {
		switch TabId(gd.tabs.ActiveTab) {
		case TabMarket:
			var cmd tea.Cmd
			if gd.keyBinder.IsBound("tui-market", msg) {
				cmd = gd.keyBinder.Execute("tui-market", msg)
			} else {
				_, cmd = gd.marketModel.Update(msg)
			}
			cmds = append(cmds, cmd)
		case TabEquipment:
			var cmd tea.Cmd
			if gd.keyBinder.IsBound("tui-eq", msg) {
				cmd = gd.keyBinder.Execute("tui-eq", msg)
			} else {
				_, cmd = gd.eqModel.Update(msg)
			}
			cmds = append(cmds, cmd)
		case TabStats:
			_, cmd := gd.statsModel.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return gd, tea.Batch(cmds...)
}
