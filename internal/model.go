package internal

import (
	"github.com/randomvlad/trader-vlads/internal/component/status"
	"github.com/randomvlad/trader-vlads/internal/component/toast"
	"github.com/randomvlad/trader-vlads/internal/screen"
)

type GameData struct {
	player     *Player
	turn       *TurnData
	toast      toast.Model
	status     status.Model
	showScreen Screen
	screenBuy  *screen.BuyScreen
	screenSell *screen.SellScreen
}

type TurnData struct {
	MarketItems  []string
	MarketPrices map[string]int
}

type Player struct {
	money       int
	week        int
	inventory   map[string]int
	activeItems []string
	unlockCost  int
	lockedItems []string
	actions     []string
	actionIndex int
}

type Screen int

const (
	ScreenMain Screen = iota
	ScreenBuy
	ScreenSell
)
