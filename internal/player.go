package internal

import (
	"fmt"
	"strings"
)

type Player struct {
	money     int
	inventory map[string]int // TODO: rename to warehouse. separate struct
	equipment map[BodyPart]*EqObject
	actions   []string
}

func NewPlayer(m *Market) *Player {

	player := Player{
		money:     1000,
		inventory: make(map[string]int),
		equipment: make(map[BodyPart]*EqObject),
		actions: []string{ // TODO: move out of player to action bar component
			"Buy",
			"Sell",
			"Next Week",
			"Unlock Item",
			"Quit",
		},
	}

	for _, item := range m.items {
		player.inventory[item.name] = 0
	}

	for _, eqObject := range getEqStarterSet() {
		player.Wear(eqObject)
	}

	return &player
}

type BodyPart int

const ( // Homage: order reflects equipment display order in Darkmists(
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

func (p *Player) isInventoryEmpty() bool {
	for _, count := range p.inventory {
		if count > 0 {
			return false
		}
	}
	return true
}

func (p *Player) Wear(object *EqObject) {
	switch object.Slot {
	case EqSlotHead:
		p.wearBodyPart(BodyPartHead, object)
	case EqSlotNeck:
		p.wearBodyPart(BodyPartNeck, object)
	case EqSlotTorso:
		p.wearBodyPart(BodyPartTorso, object)
	case EqSlotHands:
		p.wearBodyPart(BodyPartHands, object)
	case EqSlotFinger:
		if p.HasEquipment(BodyPartFingerLeft) {
			p.wearBodyPart(BodyPartFingerRight, object)
		} else {
			p.wearBodyPart(BodyPartFingerLeft, object)
		}
	case EqSlotWaist:
		p.wearBodyPart(BodyPartWaist, object)
	case EqSlotLegs:
		p.wearBodyPart(BodyPartLegs, object)
	case EqSlotFeet:
		p.wearBodyPart(BodyPartFeet, object)
	case EqSlotWield, EqSlotHold:
		if p.HasEquipment(BodyPartHoldLeft) {
			p.wearBodyPart(BodyPartHoldRight, object)
		} else {
			p.wearBodyPart(BodyPartHoldLeft, object)
		}
	}
}

func (p *Player) wearBodyPart(bodyPart BodyPart, object *EqObject) {
	if object == nil {
		return
	}

	// TODO: object eq type must match body part

	// TODO: handle case of left & right finger. Also left & right hold.

	p.Remove(bodyPart)
	p.equipment[bodyPart] = object
}

func (p *Player) HasEquipment(bodyPart BodyPart) bool {
	return p.equipment[bodyPart] != nil
}

func (p *Player) Remove(bodyPart BodyPart) {
	if _ /* eqObject */, ok := p.equipment[bodyPart]; ok {
		p.equipment[bodyPart] = nil
		// TODO: move removed eqObject to inventory?
		// need to separate storage of resources vs eq? warehouse vs inventory?
	}
}

func (p *Player) ViewEquipment() string {

	var view strings.Builder

	view.WriteString("You are using:\n")

	for bodyPart := range BodyPartHoldRight {
		if !p.HasEquipment(bodyPart) {
			continue
		}

		var wornOn string
		switch bodyPart {
		case BodyPartFingerLeft, BodyPartFingerRight:
			wornOn = "worn on finger"
		case BodyPartNeck:
			wornOn = "worn around neck"
		case BodyPartTorso:
			wornOn = "worn on torso"
		case BodyPartHead:
			wornOn = "worn on head"
		case BodyPartLegs:
			wornOn = "worn on legs"
		case BodyPartFeet:
			wornOn = "worn on feet"
		case BodyPartHands:
			wornOn = "worn on hands"
		case BodyPartWaist:
			wornOn = "worn on waist"
		case BodyPartHoldLeft, BodyPartHoldRight:
			wornOn = "held"
		}
		wornOn = "<" + wornOn + ">"

		view.WriteString(fmt.Sprintf("  %-22s %s\n", wornOn, p.equipment[bodyPart].Name))

	}

	view.WriteString("\n")
	view.WriteString("You are carrying (0):\n")
	view.WriteString("     Nothing\n") // TODO: implement inventory

	return view.String()
}
