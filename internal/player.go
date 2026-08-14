package internal

import (
	eq "github.com/randomvlad/trader-vlads/internal/appmod/equipment"
	"github.com/randomvlad/trader-vlads/internal/appmod/market"
)

type Player struct {
	money     int
	warehouse *Warehouse
	equipment map[eq.BodyPart]*eq.EqObject
	inventory []*eq.EqObject
}

func NewPlayer(m *market.Market) *Player {

	warehouse := &Warehouse{
		capacity:  100,
		resources: make(map[string]int),
	}

	player := Player{
		money:     1000,
		warehouse: warehouse,
		equipment: make(map[eq.BodyPart]*eq.EqObject),
		inventory: []*eq.EqObject{},
	}

	for _, item := range m.Resources {
		player.warehouse.resources[item.Name] = 0
	}

	for _, eqObject := range eq.GetEqStarterSet() {
		if eqObject.IsWearable() {
			player.Wear(eqObject)
		} else {
			player.AddInventory(eqObject)
		}
	}

	return &player
}

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

func (p *Player) GetInventory() []*eq.EqObject {
	return p.inventory
}

func (p *Player) GetEquipment() map[eq.BodyPart]*eq.EqObject {
	return p.equipment
}

func (p *Player) Wear(object *eq.EqObject) {
	switch object.Slot {
	case eq.EqSlotHead:
		p.wearBodyPart(eq.BodyPartHead, object)
	case eq.EqSlotNeck:
		p.wearBodyPart(eq.BodyPartNeck, object)
	case eq.EqSlotTorso:
		p.wearBodyPart(eq.BodyPartTorso, object)
	case eq.EqSlotHands:
		p.wearBodyPart(eq.BodyPartHands, object)
	case eq.EqSlotFinger:
		if p.HasEquipment(eq.BodyPartFingerLeft) {
			p.wearBodyPart(eq.BodyPartFingerRight, object)
		} else {
			p.wearBodyPart(eq.BodyPartFingerLeft, object)
		}
	case eq.EqSlotWaist:
		p.wearBodyPart(eq.BodyPartWaist, object)
	case eq.EqSlotLegs:
		p.wearBodyPart(eq.BodyPartLegs, object)
	case eq.EqSlotFeet:
		p.wearBodyPart(eq.BodyPartFeet, object)
	case eq.EqSlotWield, eq.EqSlotHold:
		if p.HasEquipment(eq.BodyPartHoldLeft) {
			p.wearBodyPart(eq.BodyPartHoldRight, object)
		} else {
			p.wearBodyPart(eq.BodyPartHoldLeft, object)
		}
	}
}

func (p *Player) wearBodyPart(bodyPart eq.BodyPart, object *eq.EqObject) {
	if object != nil {
		p.Remove(bodyPart)
		p.equipment[bodyPart] = object
	}
}

func (p *Player) HasEquipment(bodyPart eq.BodyPart) bool {
	_, has := p.equipment[bodyPart]
	return has
}

func (p *Player) Remove(bodyPart eq.BodyPart) {
	if eqObject, ok := p.equipment[bodyPart]; ok {
		delete(p.equipment, bodyPart)
		p.AddInventory(eqObject)
	}
}

func (p *Player) AddInventory(object *eq.EqObject) {
	p.inventory = append(p.inventory, object)
}
