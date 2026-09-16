package equipment

import (
	"github.com/oklog/ulid/v2"
	"github.com/randomvlad/trader-vlads/internal/util"
)

var Forge = newEqForge()

type EqForge struct {
	defRegistry *EqDefRegistry
}

func newEqForge() *EqForge {
	return &EqForge{
		defRegistry: NewEqDefRegistry(),
	}
}

func (f *EqForge) Make(random *util.RandomGenerator, eqDefs ...string) []*EqObject {

	var objects []*EqObject

	for _, def := range eqDefs {
		eqDef := f.defRegistry.definitions[def]
		if eqDef == nil {
			continue
		}

		eqObject := &EqObject{
			Id:     ulid.Make(),
			Name:   eqDef.Name,
			Slot:   eqDef.Slot,
			Usable: eqDef.Usable,
		}

		if eqDef.StatDefense != nil {
			eqObject.Attributes = append(eqObject.Attributes, createStat(AttributeDefense, eqDef.StatDefense, random))
		}

		if eqDef.StatSpeed != nil {
			eqObject.Attributes = append(eqObject.Attributes, createStat(AttributeSpeed, eqDef.StatSpeed, random))
		}

		if eqDef.StatIntuition != nil {
			eqObject.Attributes = append(eqObject.Attributes, createStat(AttributeIntuition, eqDef.StatIntuition, random))
		}

		if eqDef.StatWillpower != nil {
			eqObject.Attributes = append(eqObject.Attributes, createStat(AttributeWillpower, eqDef.StatWillpower, random))
		}

		if eqDef.StatOpulence != nil {
			eqObject.Attributes = append(eqObject.Attributes, createStat(AttributeOpulence, eqDef.StatOpulence, random))
		}

		for _, effectDef := range eqDef.EffectDefs {
			effect := effectDef.Create(random, eqObject.Id, "equipment")
			eqObject.Effects = append(eqObject.Effects, effect)
		}

		objects = append(objects, eqObject)
	}

	return objects
}

func createStat(attributeType AttributeType, valueRange *util.RangeInt, random *util.RandomGenerator) StatAttribute {
	return StatAttribute{
		AttributeType: attributeType,
		Value:         valueRange.Generate(random),
	}
}
