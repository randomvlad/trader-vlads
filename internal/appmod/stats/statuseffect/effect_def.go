package statuseffect

import (
	"github.com/oklog/ulid/v2"
	"github.com/randomvlad/trader-vlads/internal/util"
)

type EffectInstanceCreator interface {
	Create(r *util.RandomGenerator, grantById ulid.ULID, grandByType string) StatusEffect
}

type BaseEffectDef struct {
	EffectInstanceCreator
	Name         string
	Duration     *EffectDuration
	MessageStart string
	MessageEnd   string
}

type EffectDuration struct {
	PermanentDuration bool
	ChanceRange       *util.RangeInt
	TurnsRange        *util.RangeInt
}

func NewDuration() *EffectDuration {
	return &EffectDuration{}
}

func (e *EffectDuration) Turns(valueMin, valueMax int) *EffectDuration {
	e.TurnsRange = util.NewRangeInt(valueMin, valueMax)
	return e
}

func (e *EffectDuration) Permanent() *EffectDuration {
	e.PermanentDuration = true
	return e
}

func (e *EffectDuration) Chance(valueMin, valueMax int) *EffectDuration {
	e.ChanceRange = util.NewRangeInt(valueMin, valueMax)
	return e
}
