package statuseffect

import (
	"github.com/oklog/ulid/v2"
	"github.com/randomvlad/trader-vlads/internal/util"
)

type StatusEffect interface {
	Id() ulid.ULID
	Name() string
	GrantedById() ulid.ULID
	GrantedByType() string
	View() string
	GetMessageStart() string
	GetMessageEnd() string
	HasEnded() bool
	GetTurns() int
	Apply(player PlayerEffectService)
}

type StatusEffectDef interface {
	Create(r *util.RandomGenerator, grantById ulid.ULID, grandByType string) StatusEffect
}

type BaseEffectDef struct {
	StatusEffectDef
	Name         string
	Permanent    bool
	Chance       *util.RangeInt
	Turns        *util.RangeInt
	MessageStart string
	MessageEnd   string
}

type PlayerEffectService interface {
	AddMoney(amount int)
	AddResourceQuantity(name string, quantity int)
}

type BaseStatusEffect struct {
	StatusEffect
	id            ulid.ULID
	name          string
	grantedById   ulid.ULID
	grantedByType string
	permanent     bool
	turnsLeft     int
	chance        int
	messageStart  string
	messageEnd    string
	applyFunc     func(player PlayerEffectService)
	random        *util.RandomGenerator
}

func NewBaseStatusEffect(baseDef *BaseEffectDef, r *util.RandomGenerator, grantById ulid.ULID, grandByType string) *BaseStatusEffect {

	base := &BaseStatusEffect{
		id:            ulid.Make(),
		name:          baseDef.Name,
		grantedById:   grantById,
		grantedByType: grandByType,
		random:        r,
		messageStart:  baseDef.MessageStart,
		messageEnd:    baseDef.MessageEnd,
	}

	if baseDef.Chance != nil {
		base.chance = baseDef.Chance.Generate(r)
	}

	if baseDef.Permanent {
		base.permanent = true
	} else {
		base.turnsLeft = baseDef.Turns.Generate(r)
	}

	return base
}

func (b *BaseStatusEffect) Id() ulid.ULID {
	return b.id
}

func (b *BaseStatusEffect) Name() string {
	return b.name
}

func (b *BaseStatusEffect) GrantedById() ulid.ULID {
	return b.grantedById
}

func (b *BaseStatusEffect) GrantedByType() string {
	return b.grantedByType
}

func (b *BaseStatusEffect) HasEnded() bool {
	if b.permanent {
		return false
	} else {
		return b.turnsLeft < 1
	}
}

func (b *BaseStatusEffect) GetTurns() int {
	return b.turnsLeft
}

func (b *BaseStatusEffect) Apply(player PlayerEffectService) {

	var shouldApply bool
	if b.random != nil && b.chance > 0 {
		shouldApply = b.random.RollChance(b.chance)
	} else {
		shouldApply = true
	}

	if shouldApply {
		b.applyFunc(player)
	}

	if !b.permanent {
		b.turnsLeft--
	}
}

func (b *BaseStatusEffect) GetMessageStart() string {
	return b.messageStart
}

func (b *BaseStatusEffect) GetMessageEnd() string {
	return b.messageEnd
}
