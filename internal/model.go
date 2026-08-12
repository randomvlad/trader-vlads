package internal

import (
	"github.com/randomvlad/trader-vlads/internal/component/status"
	"github.com/randomvlad/trader-vlads/internal/component/tabs"
	"github.com/randomvlad/trader-vlads/internal/component/toast"
	"github.com/randomvlad/trader-vlads/internal/screen"
)

type GameData struct {
	player     *Player
	market     *Market
	eventTrack *EventTracker
	turn       *TurnData
	tabs       *tabs.TabsModel
	toast      toast.Model
	status     status.Model
	showScreen Screen
	screenBuy  *screen.BuyScreen
	screenSell *screen.SellScreen
}

type TurnData struct {
	// Note: keeping around for now
}

type Screen int

const (
	ScreenMain Screen = iota
	ScreenBuy
	ScreenSell
)
