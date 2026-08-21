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

		for _, effectDef := range eqDef.EffectDefs {
			effect := effectDef.Create(random, eqObject.Id, "equipment")
			eqObject.Effects = append(eqObject.Effects, effect)
		}

		objects = append(objects, eqObject)
	}

	return objects
}
