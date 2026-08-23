package equipment

import (
	"github.com/randomvlad/trader-vlads/internal/appmod/stats/statuseffect"
	"github.com/randomvlad/trader-vlads/internal/util"
)

type eqObjectDef struct {
	Name       string
	Slot       EqSlot
	Usable     bool
	EffectDefs []statuseffect.StatusEffectDef
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
			EffectDefs: []statuseffect.StatusEffectDef{
				&statuseffect.GrantGoldEffectDef{
					Name:         "Beginner's Luck 🍀",
					Turns:        util.NewRangeInt(4, 6),
					Gold:         util.NewRangeInt(10, 15),
					MessageStart: "You are beginning to feel naively optimistic and unconcerned.",
					MessageEnd:   "Your carefree perspective and feeling of unabashed optimism fades.",
				},
				&statuseffect.GrantResourceEffectDef{
					Name:     "Beaver's Bounty",
					Resource: "Wood",
					Turns:    util.NewRangeInt(4, 6),
					Amount:   util.NewRangeInt(1, 1),
				},
			},
		},
		{
			Name: "goose with feathers of pure gold",
			Slot: EqSlotHold,
			EffectDefs: []statuseffect.StatusEffectDef{
				&statuseffect.GrantGoldEffectDef{
					Name:         "Blessings of the Honk 🪿",
					Permanent:    true,
					Gold:         util.NewRangeInt(25, 35),
					MessageStart: "The goose gently honks and starts following you.",
					MessageEnd:   "The goose waddles away from you.",
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
			EffectDefs: []statuseffect.StatusEffectDef{
				&statuseffect.GrantResourceEffectDef{
					Name:      "Inexhaustible Cart of Lumber",
					Resource:  "Wood",
					Permanent: true,
					Chance:    util.NewRangeInt(45, 65),
					Amount:    util.NewRangeInt(2, 3),
				},
			},
		},
	}
}
