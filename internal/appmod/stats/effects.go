package stats

import (
	"github.com/randomvlad/trader-vlads/internal/util"
)

type StatusEffect interface {
	Name() string
	View() string
	GetAppliedMessage() string
	GetExpiredMessage() string
	HasExpired() bool
	Apply(player PlayerStatsService)
}

type BeginnersLuckEffect struct {
	name          string
	durationTurns int
	turnsLeft     int
	grantMoney    int
}

func NewEffectBeginnersLuck(turns int, money int) StatusEffect {
	// TODO: random grants/effects: +[5-10] gold for [3-6] turns
	return &BeginnersLuckEffect{
		name:          "Beginner's Luck 🍀",
		durationTurns: turns,
		turnsLeft:     turns,
		grantMoney:    money,
	}
}

func (e *BeginnersLuckEffect) Name() string {
	return e.name
}

func (e *BeginnersLuckEffect) GetAppliedMessage() string {
	return "You are beginning to feel naively optimistic and unconcerned."
}

// TODO: use it as a toast at the beginning of turn?
func (e *BeginnersLuckEffect) GetExpiredMessage() string {
	return "Your carefree perspective and feeling of unabashed optimism fades."
}

func (e *BeginnersLuckEffect) HasExpired() bool {
	return e.turnsLeft <= 0
}

// TODO: two separate views ... 1) definition and 2) in-progress
func (e *BeginnersLuckEffect) View() string {
	// TODO: if e.turnsLeft == 1 then append ⌛? Makes sense for in-progress view of stats,
	// but less so when showing a selected item that hasn't been used yet
	return "Grants " + util.FormatMoney(e.grantMoney) + " for " +
		util.FormatCountPluralized(e.durationTurns, "week")
}

func (e *BeginnersLuckEffect) Apply(player PlayerStatsService) {
	player.AddMoney(e.grantMoney)
	e.turnsLeft--
}
