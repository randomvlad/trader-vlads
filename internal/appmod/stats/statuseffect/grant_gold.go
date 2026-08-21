package statuseffect

import (
	"github.com/oklog/ulid/v2"
	"github.com/randomvlad/trader-vlads/internal/util"
)

type GrantGoldEffectDef struct {
	TurnsMin int
	TurnsMax int
	MoneyMin int // can be combined into a Range(2, 4)
	MoneyMax int
}

// TODO: can create struct CreateParam: grantById, grantByType, rand
func (d *GrantGoldEffectDef) Create(r *util.RandomGenerator, grantById ulid.ULID, grandByType string) StatusEffect {

	return &BeginnersLuckEffect{
		id:          ulid.Make(),
		name:        "Beginner's Luck 🍀",
		grantById:   grantById,
		grantByType: grandByType,
		turnsLeft:   r.RollInt(d.TurnsMin, d.TurnsMax),
		grantMoney:  r.RollInt(d.MoneyMin, d.MoneyMax),
	}
}

type BeginnersLuckEffect struct {
	id          ulid.ULID
	name        string
	grantById   ulid.ULID
	grantByType string
	turnsLeft   int
	grantMoney  int
}

func (e *BeginnersLuckEffect) Id() ulid.ULID {
	return e.id
}

func (e *BeginnersLuckEffect) Name() string {
	return e.name
}

func (e *BeginnersLuckEffect) GrantedById() ulid.ULID {
	return e.grantById
}

func (e *BeginnersLuckEffect) GrantedByType() string {
	return e.grantByType
}

func (e *BeginnersLuckEffect) GetAppliedMessage() string {
	return "You are beginning to feel naively optimistic and unconcerned."
}

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

func (e *BeginnersLuckEffect) Apply(player PlayerEffectService) {
	player.AddMoney(e.grantMoney)
	e.turnsLeft--
}
