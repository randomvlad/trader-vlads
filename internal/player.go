package internal

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/randomvlad/trader-vlads/internal/appmod/market"
)

type Player struct {
	money     int
	warehouse *Warehouse
	equipment map[BodyPart]*EqObject
	inventory []*EqObject
}

func NewPlayer(m *market.Market) *Player {

	warehouse := &Warehouse{
		capacity:  100,
		resources: make(map[string]int),
	}

	player := Player{
		money:     1000,
		warehouse: warehouse,
		equipment: make(map[BodyPart]*EqObject),
		inventory: []*EqObject{},
	}

	for _, item := range m.Resources {
		player.warehouse.resources[item.Name] = 0
	}

	for _, eqObject := range getEqStarterSet() {
		if eqObject.IsWearable() {
			player.Wear(eqObject)
		} else {
			player.AddInventory(eqObject)
		}
	}

	return &player
}

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

type Warehouse struct {
	capacity  int
	resources map[string]int // Resource name → quantity
}

func (p *Player) GetMoney() int {
	return p.money
}

func (p *Player) AddMoney(amount int) {
	p.money += amount
}

func (p *Player) AddResourceQuantity(name string, quantity int) {
	p.warehouse.resources[name] += quantity
}

func (p *Player) IsWarehouseEmpty() bool {
	for _, count := range p.warehouse.resources {
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
	if object != nil {
		p.Remove(bodyPart)
		p.equipment[bodyPart] = object
	}
}

func (p *Player) HasEquipment(bodyPart BodyPart) bool {
	return p.equipment[bodyPart] != nil
}

func (p *Player) Remove(bodyPart BodyPart) {
	if eqObject, ok := p.equipment[bodyPart]; ok {
		p.equipment[bodyPart] = nil
		p.AddInventory(eqObject)
	}
}

func (p *Player) AddInventory(object *EqObject) {
	p.inventory = append(p.inventory, object)
}

func (p *Player) ViewEquipment() string {

	var view strings.Builder

	view.WriteString("You are using:\n")

	naked := true
	for bodyPart := range BodyPartHoldRight {
		if !p.HasEquipment(bodyPart) {
			continue
		}

		naked = false

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

		view.WriteString(fmt.Sprintf("     %-22s %s\n", wornOn, p.equipment[bodyPart].Name))
	}

	if naked {
		view.WriteString("     Nothing\n")
	}

	view.WriteString("\n")

	count := len(p.inventory)
	view.WriteString("You are carrying (" + strconv.Itoa(count) + "):\n")

	if count > 0 {
		for _, object := range p.inventory {
			view.WriteString("     " + object.Name + "\n")
		}
	} else {
		view.WriteString("     Nothing\n")
	}

	return view.String()
}
