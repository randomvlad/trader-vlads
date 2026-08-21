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
	amount      int
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
	return "Grants +" + strconv.Itoa(e.amount) + " " + e.resource + " " + util.ViewTurnsLeft(e.permanent, e.turnsLeft)
}

func (e *GrantResourceEffect) GetAppliedMessage() string {
	return "May the power of " + e.resource + " be with you. Always!"
}

func (e *GrantResourceEffect) GetExpiredMessage() string {
	if e.permanent {
		return ""
	} else {
		return "The power of " + e.resource + " has left you ..."
	}
}

func (e *GrantResourceEffect) HasExpired() bool {
	if e.permanent {
		return false
	} else {
		return e.turnsLeft < 1
	}
}

func (e *GrantResourceEffect) Apply(player PlayerEffectService) {
	player.AddResourceQuantity(e.resource, e.amount)
	if !e.permanent {
		e.turnsLeft--
	}
}
