package statuseffect

import (
	"strconv"

	"github.com/oklog/ulid/v2"
	"github.com/randomvlad/trader-vlads/internal/util"
)

type GrantResourceEffectDef struct {
	Name      string
	Resource  string
	Permanent bool
	Chance    *util.RangeInt
	Amount    *util.RangeInt
	Turns     *util.RangeInt
}

func (def *GrantResourceEffectDef) Create(r *util.RandomGenerator, grantById ulid.ULID, grandByType string) StatusEffect {
	effect := &GrantResourceEffect{
		id:          ulid.Make(),
		name:        def.Name,
		grantById:   grantById,
		grantByType: grandByType,
		resource:    def.Resource,
		amount:      def.Amount.Generate(r),
	}

	if def.Chance != nil {
		effect.chance = def.Chance.Generate(r)
		effect.random = r
	}

	if def.Permanent {
		effect.permanent = true
	} else {
		effect.turnsLeft = def.Turns.Generate(r)
	}

	return effect
}

type GrantResourceEffect struct {
	id          ulid.ULID
	name        string
	grantById   ulid.ULID
	grantByType string
	permanent   bool
	turnsLeft   int
	resource    string
	chance      int
	amount      int
	random      *util.RandomGenerator
}

func (e *GrantResourceEffect) Id() ulid.ULID {
	return e.id
}

func (e *GrantResourceEffect) Name() string {
	return e.name
}

func (e *GrantResourceEffect) GrantedById() ulid.ULID {
	return e.grantById
}

func (e *GrantResourceEffect) GrantedByType() string {
	return e.grantByType
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

func (e *GrantResourceEffect) GetMessageStart() string {
	return "May the power of " + e.resource + " be with you. Always!"
}

func (e *GrantResourceEffect) GetMessageEnd() string {
	return "The power of " + e.resource + " has left you ..."
}

func (e *GrantResourceEffect) HasEnded() bool {
	if e.permanent {
		return false
	} else {
		return e.turnsLeft < 1
	}
}

func (e *GrantResourceEffect) GetTurns() int {
	return e.turnsLeft
}

func (e *GrantResourceEffect) Apply(player PlayerEffectService) {

	var grantResource bool
	if e.random != nil && e.chance > 0 {
		grantResource = e.random.RollChance(e.chance)
	} else {
		grantResource = true
	}

	if grantResource {
		player.AddResourceQuantity(e.resource, e.amount)
	}

	if !e.permanent {
		e.turnsLeft--
	}
}
