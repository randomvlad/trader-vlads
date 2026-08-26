package statuseffect

import (
	"strconv"

	"github.com/oklog/ulid/v2"
	"github.com/randomvlad/trader-vlads/internal/util"
)

type GrantResourceEffectDef struct {
	*BaseEffectDef
	Resource string
	Amount   *util.RangeInt
}

func (def *GrantResourceEffectDef) Create(r *util.RandomGenerator, grantById ulid.ULID, grandByType string) StatusEffect {
	effect := &GrantResourceEffect{
		BaseStatusEffect: NewBaseStatusEffect(def.BaseEffectDef, r, grantById, grandByType),
		resource:         def.Resource,
		amount:           def.Amount.Generate(r),
	}

	effect.applyFunc = func(player PlayerEffectService) {
		player.AddResourceQuantity(effect.resource, effect.amount)
	}

	effect.messageStart = "May the power of " + effect.resource + " be with you. Always!"
	effect.messageEnd = "The power of " + effect.resource + " has left you ..."

	return effect
}

type GrantResourceEffect struct {
	*BaseStatusEffect
	resource string
	amount   int
}

func (e *GrantResourceEffect) View() string {

	var grantsDisplay string
	if e.chance > 0 {
		grantsDisplay = strconv.Itoa(e.chance) + "% to grant"
	} else {
		grantsDisplay = "Grants"
	}

	return grantsDisplay + " +" + strconv.Itoa(e.amount) + " " + e.resource + " " + util.ViewTurnsLeft(e.permanent, e.turnsLeft)
}
