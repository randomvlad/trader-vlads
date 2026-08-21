package equipment

import (
	"github.com/randomvlad/trader-vlads/internal/appmod/stats/statuseffect"
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
					TurnsMin: 4,
					TurnsMax: 6,
					MoneyMin: 10,
					MoneyMax: 15,
				},
			},
		},
		{
			Name: "goose with feathers of pure gold",
			Slot: EqSlotHold,
			EffectDefs: []statuseffect.StatusEffectDef{
				&statuseffect.GrantGoldEffectDef{ // TODO: support name for effect
					TurnsMin: 100, // TODO: support permanent
					TurnsMax: 100,
					MoneyMin: 25,
					MoneyMax: 25,
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
					AmountMin: 1,
					AmountMax: 1,
					TurnsMin:  100, // TODO: implement permanent
					TurnsMax:  200,
				},
			},
		},
	}
}
