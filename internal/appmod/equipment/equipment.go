package equipment

import (
	"strings"

	"github.com/randomvlad/trader-vlads/internal/appmod/stats"
)

type EqObject struct {
	Name    string
	Slot    EqSlot
	Usable  bool
	Effects []stats.StatusEffect
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

type BodyPart int

const ( // Homage: order reflects equipment display order in Darkmists
	BodyPartFingerLeft BodyPart = iota
	BodyPartFingerRight
	BodyPartNeck
	BodyPartTorso
	BodyPartHead
	BodyPartLegs
	BodyPartFeet
	BodyPartHands
	BodyPartWaist
	BodyPartHoldLeft
	BodyPartHoldRight
)

const BodyPartsMax = int(BodyPartHoldRight) + 1

func GetEqStarterSet() []*EqObject {

	return []*EqObject{
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
			Name:   "a potion of Beginner's Luck 🍀",
			Slot:   EqSlotInventory,
			Usable: true,
			Effects: []stats.StatusEffect{
				stats.NewEffectBeginnersLuck(6, 10),
			},
		},
		{
			Name:   "a jar of spicy pickles",
			Slot:   EqSlotInventory,
			Usable: true,
		},
		{
			Name: "quill made from a talon of the Blue Dragon",
			Slot: EqSlotHold,
		},
	}
}

func (o *EqObject) IsWearable() bool {
	return o.Slot != EqSlotInventory
}

func (o *EqObject) IsUsable() bool {
	return o.Usable
}

func (o *EqObject) ViewStats() string {
	var view strings.Builder
	view.WriteString("Object: " + o.Name + "\n")
	view.WriteString("Type: " + getEqSlotName(o.Slot) + "\n")

	if len(o.Effects) > 0 {
		for _, effect := range o.Effects {
			view.WriteString("Effect: " + effect.View() + "\n")
		}
	}

	return view.String()
}

func getEqSlotName(eqSlot EqSlot) string {
	var name string
	switch eqSlot {
	case EqSlotHead:
		name = "Head"
	case EqSlotNeck:
		name = "Neck"
	case EqSlotTorso:
		name = "Torso"
	case EqSlotHands:
		name = "Hands"
	case EqSlotFinger:
		name = "Finger"
	case EqSlotWaist:
		name = "Waist"
	case EqSlotLegs:
		name = "Legs"
	case EqSlotFeet:
		name = "Feet"
	case EqSlotWield, EqSlotHold:
		name = "Hold"
	case EqSlotInventory:
		name = "Inventory"
	}
	return name
}
