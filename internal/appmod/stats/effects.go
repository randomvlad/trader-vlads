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

func (e *BeginnersLuckEffect) View() string {
	var expiresSoonIcon string
	if e.turnsLeft == 1 {
		expiresSoonIcon = " ⌛"
	}

	return "Grants " + util.FormatMoney(e.grantMoney) + " for " +
		util.FormatCountPluralized(e.turnsLeft, "week") + expiresSoonIcon
}

func (e *BeginnersLuckEffect) Apply(player PlayerStatsService) {
	player.AddMoney(e.grantMoney)
	e.turnsLeft--
}
