package statuseffect

import (
	"github.com/oklog/ulid/v2"
	"github.com/randomvlad/trader-vlads/internal/util"
)

type GrantGoldEffectDef struct {
	Name      string
	Permanent bool
	Turns     *util.RangeInt
	Gold      *util.RangeInt
}

// TODO: consider creating struct EffectCreateParams to centralize method arguments
func (def *GrantGoldEffectDef) Create(r *util.RandomGenerator, grantById ulid.ULID, grandByType string) StatusEffect {

	effect := &BeginnersLuckEffect{
		id:          ulid.Make(),
		name:        def.Name,
		grantById:   grantById,
		grantByType: grandByType,
		goldGain:    def.Gold.Generate(r),
	}

	if def.Permanent {
		effect.permanent = true
	} else {
		effect.turnsLeft = def.Turns.Generate(r)
	}

	return effect
}

type BeginnersLuckEffect struct {
	id          ulid.ULID
	name        string
	grantById   ulid.ULID
	grantByType string
	permanent   bool
	turnsLeft   int
	goldGain    int
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
	return "You are beginning to feel naively optimistic and unconcerned." // TODO: make more generic or allow to override
}

func (e *BeginnersLuckEffect) GetExpiredMessage() string {
	return "Your carefree perspective and feeling of unabashed optimism fades." // TODO: make more generic or allow to override
}

func (e *BeginnersLuckEffect) HasExpired() bool {
	return e.turnsLeft <= 0
}

func (e *BeginnersLuckEffect) View() string {
	return "Grants +" + util.FormatMoney(e.goldGain) + " " + util.ViewTurnsLeft(e.permanent, e.turnsLeft)
}

func (e *BeginnersLuckEffect) Apply(player PlayerEffectService) {
	player.AddMoney(e.goldGain)

	if !e.permanent {
		e.turnsLeft--
	}
}
