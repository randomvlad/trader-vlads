package statuseffect

import (
	"github.com/oklog/ulid/v2"
	"github.com/randomvlad/trader-vlads/internal/util"
)

type GrantGoldEffectDef struct {
	Name         string
	Permanent    bool
	Turns        *util.RangeInt
	Gold         *util.RangeInt
	MessageStart string
	MessageEnd   string
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

	effect.messageStart = util.GetOrDefault(def.MessageStart, "The tides of prosperity seem to favor you.")
	effect.messageEnd = util.GetOrDefault(def.MessageEnd, "The tides of prosperity no longer favor you.")

	return effect
}

type BeginnersLuckEffect struct {
	id           ulid.ULID
	name         string
	grantById    ulid.ULID
	grantByType  string
	permanent    bool
	turnsLeft    int
	goldGain     int
	messageStart string
	messageEnd   string
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

func (e *BeginnersLuckEffect) GetMessageStart() string {
	return e.messageStart
}

func (e *BeginnersLuckEffect) GetMessageEnd() string {
	return e.messageEnd
}

func (e *BeginnersLuckEffect) HasEnded() bool {
	if e.permanent {
		return false
	} else {
		return e.turnsLeft < 1
	}
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
