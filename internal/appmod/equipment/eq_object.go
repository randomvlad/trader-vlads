package equipment

import (
	"strings"

	"github.com/oklog/ulid/v2"
	"github.com/randomvlad/trader-vlads/internal/appmod/stats/statuseffect"
)

type EqObject struct {
	Id      ulid.ULID
	Name    string
	Slot    EqSlot
	Usable  bool
	Effects []statuseffect.StatusEffect
}

type BodyPart int

const BodyPartsMax = int(BodyPartHoldRight) + 1

const ( // Homage: consts order reflects equipment display order in Darkmists
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

	if len(o.Effects) == 1 {
		view.WriteString("Effect: " + o.Effects[0].View() + "\n")
	} else if len(o.Effects) > 1 {
		view.WriteString("Effects:\n")
		for _, effect := range o.Effects {
			view.WriteString("  ‣ " + effect.View() + "\n")
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
