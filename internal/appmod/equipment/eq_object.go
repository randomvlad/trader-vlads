package equipment

import (
	"fmt"

	"github.com/oklog/ulid/v2"
	"github.com/randomvlad/trader-vlads/internal/appmod/stats/statuseffect"
	"github.com/randomvlad/trader-vlads/internal/util/stringutil"
)

type EqObject struct {
	Id         ulid.ULID
	Name       string
	Slot       EqSlot
	Usable     bool
	Effects    []statuseffect.StatusEffect
	Attributes []StatAttribute
}

type AttributeType int

const (
	AttributeDefense AttributeType = iota
	AttributeSpeed
	AttributeIntuition
	AttributeWillpower
	AttributeOpulence
)

type StatAttribute struct {
	AttributeType AttributeType
	Value         int
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
	view := stringutil.NewBuilder().
		WriteLn("Object: " + o.Name).
		WriteLn("Type: " + getEqSlotName(o.Slot))

	if len(o.Attributes) == 1 {
		view.WriteLn("Stat: " + formatAttribute(o.Attributes[0]))
	} else if len(o.Attributes) > 1 {
		view.WriteLn("Stats:")
		for _, attribute := range o.Attributes {
			view.Write("    ").WriteLn(formatAttribute(attribute))
		}
	}

	if len(o.Effects) == 1 {
		view.WriteLn("Effect: " + o.Effects[0].View())
	} else if len(o.Effects) > 1 {
		view.WriteLn("Effects:")
		for _, effect := range o.Effects {
			view.WriteLn("  ‣ " + effect.View())
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

func formatAttribute(attribute StatAttribute) string {
	var displayType string

	switch attribute.AttributeType {
	case AttributeDefense:
		displayType = "Defense"
	case AttributeSpeed:
		displayType = "Speed"
	case AttributeIntuition:
		displayType = "Intuition"
	case AttributeWillpower:
		displayType = "Willpower"
	case AttributeOpulence:
		displayType = "Opulence"
	}

	return fmt.Sprintf("%+d %s", attribute.Value, displayType)
}
