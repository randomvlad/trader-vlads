package equipment

import (
	eff "github.com/randomvlad/trader-vlads/internal/appmod/stats/statuseffect"
	"github.com/randomvlad/trader-vlads/internal/util"
)

type eqObjectDef struct {
	Name       string
	Slot       EqSlot
	Usable     bool
	EffectDefs []eff.EffectInstanceCreator
}

type EqSlot int

const (
	EqSlotHead EqSlot = iota
	EqSlotNeck
	EqSlotTorso
	EqSlotHands
	EqSlotFinger
	EqSlotWaist
	EqSlotLegs
	EqSlotFeet
	EqSlotWield
	EqSlotHold
	EqSlotInventory
)

type EqDefRegistry struct {
	definitions map[string]*eqObjectDef
}

func NewEqDefRegistry() *EqDefRegistry {
	registry := &EqDefRegistry{
		definitions: make(map[string]*eqObjectDef),
	}

	for _, def := range getDefinitions() {
		registry.definitions[def.Name] = def
	}

	return registry
}

func getDefinitions() []*eqObjectDef {
	return []*eqObjectDef{
		{
			Name: "copper ring of a novice",
			Slot: EqSlotFinger,
		},
		{
			Name: "gray cotton tunic",
			Slot: EqSlotTorso,
		},
		{
			Name: "worn trousers",
			Slot: EqSlotLegs,
		},
		{
			Name: "brown leather sandals",
			Slot: EqSlotFeet,
		},
		{
			Name: "quill made from a talon of the Blue Dragon",
			Slot: EqSlotHold,
		},
		{
			Name:   "a potion of Beginner's Luck 🍀",
			Slot:   EqSlotInventory,
			Usable: true,
			EffectDefs: []eff.EffectInstanceCreator{
				&eff.GrantGoldEffectDef{
					BaseEffectDef: &eff.BaseEffectDef{
						Name:         "Beginner's Luck 🍀",
						Duration:     eff.NewDuration().Turns(4, 6),
						MessageStart: "You are beginning to feel naively optimistic and unconcerned.",
						MessageEnd:   "Your carefree perspective and feeling of unabashed optimism fades.",
					},
					Gold: util.NewRangeInt(10, 15),
				},
				&eff.GrantResourceEffectDef{
					BaseEffectDef: &eff.BaseEffectDef{
						Name:     "Beaver's Bounty",
						Duration: eff.NewDuration().Turns(4, 6),
					},
					Resource: "Wood",
					Amount:   util.NewRangeInt(1, 1),
				},
			},
		},
		{
			Name: "goose with feathers of pure gold",
			Slot: EqSlotHold,
			EffectDefs: []eff.EffectInstanceCreator{
				&eff.GrantGoldEffectDef{
					BaseEffectDef: &eff.BaseEffectDef{
						Name:         "Blessings of the Honk 🪿",
						Duration:     eff.NewDuration().Permanent(),
						MessageStart: "The goose gently honks and starts following you.",
						MessageEnd:   "The goose waddles away from you.",
					},
					Gold: util.NewRangeInt(25, 35),
				},
			},
		},
		{
			Name:   "a jar of spicy pickles",
			Slot:   EqSlotInventory,
			Usable: true,
		},
		{
			Name: "Inexhaustible Cart of Lumber", // HoMM3 Homage https://homm.fandom.com/wiki/Inexhaustible_Cart_of_Lumber
			Slot: EqSlotHold,
			EffectDefs: []eff.EffectInstanceCreator{
				&eff.GrantResourceEffectDef{
					BaseEffectDef: &eff.BaseEffectDef{
						Name:     "Inexhaustible Cart of Lumber",
						Duration: eff.NewDuration().Permanent().Chance(45, 65),
					},
					Resource: "Wood",
					Amount:   util.NewRangeInt(2, 3),
				},
			},
		},
	}
}
