package statuseffect

import (
	"github.com/oklog/ulid/v2"
	"github.com/randomvlad/trader-vlads/internal/util"
)

type GrantGoldEffectDef struct {
	*BaseEffectDef
	Gold *util.RangeInt
}

func (def *GrantGoldEffectDef) Create(r *util.RandomGenerator, grantById ulid.ULID, grandByType string) StatusEffect {

	effect := &GrantGoldEffect{
		BaseStatusEffect: NewBaseStatusEffect(def.BaseEffectDef, r, grantById, grandByType),
		goldGain:         def.Gold.Generate(r),
	}

	effect.messageStart = util.GetOrDefault(def.MessageStart, "The tides of prosperity seem to favor you.")
	effect.messageEnd = util.GetOrDefault(def.MessageEnd, "The tides of prosperity no longer favor you.")

	effect.applyFunc = func(player PlayerEffectService) {
		player.AddMoney(effect.goldGain)
	}

	return effect
}

type GrantGoldEffect struct {
	*BaseStatusEffect
	goldGain int
}

func (e *GrantGoldEffect) View() string {
	return "Grants +" + util.FormatCurrency(e.goldGain) + " " + util.ViewTurnsLeft(e.permanent, e.turnsLeft)
}
