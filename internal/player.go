package internal

import (
	"slices"

	eq "github.com/randomvlad/trader-vlads/internal/appmod/equipment"
	"github.com/randomvlad/trader-vlads/internal/appmod/market"
	"github.com/randomvlad/trader-vlads/internal/appmod/stats"
	"github.com/randomvlad/trader-vlads/internal/util"
)

type Player struct {
	money     int
	warehouse *Warehouse
	equipment map[eq.BodyPart]*eq.EqObject
	inventory []*eq.EqObject
	effects   []stats.StatusEffect
	random    *util.RandomGenerator
}

func NewPlayer(m *market.Market, r *util.RandomGenerator) *Player {

	warehouse := &Warehouse{
		capacity:  100,
		resources: make(map[string]int),
	}

	player := Player{
		money:     1000,
		warehouse: warehouse,
		equipment: make(map[eq.BodyPart]*eq.EqObject),
		inventory: []*eq.EqObject{},
		random:    r,
	}

	for _, item := range m.Resources {
		player.warehouse.resources[item.Name] = 0
	}

	// initialize every eq body part
	for bodyPart := range eq.BodyPartsMax {
		player.equipment[eq.BodyPart(bodyPart)] = nil
	}

	for _, eqObject := range eq.GetEqStarterSet(r) {
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

func (p *Player) GetInventoryObject(index int) *eq.EqObject {
	if index >= 0 || index < len(p.inventory) {
		return p.inventory[index]
	} else {
		return nil
	}
}

func (p *Player) GetEquipped() map[eq.BodyPart]*eq.EqObject {
	return p.equipment
}

func (p *Player) GetEquippedObject(bodyPart eq.BodyPart) *eq.EqObject {
	return p.equipment[bodyPart]
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
		if p.HasEquipped(eq.BodyPartFingerLeft) {
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
		if p.HasEquipped(eq.BodyPartHoldLeft) {
			p.wearBodyPart(eq.BodyPartHoldRight, object)
		} else {
			p.wearBodyPart(eq.BodyPartHoldLeft, object)
		}
	}
}

func (p *Player) WearInventory(invIndex int) bool {
	eqObject := p.GetInventoryObject(invIndex)
	if eqObject == nil || !eqObject.IsWearable() {
		return false
	}

	p.Wear(eqObject)
	p.inventory = slices.Delete(p.inventory, invIndex, invIndex+1)
	return true
}

func (p *Player) wearBodyPart(bodyPart eq.BodyPart, object *eq.EqObject) {
	if object != nil {
		p.Remove(bodyPart)
		p.equipment[bodyPart] = object
	}
}

func (p *Player) HasEquipped(bodyPart eq.BodyPart) bool {
	return p.equipment[bodyPart] != nil
}

func (p *Player) Remove(bodyPart eq.BodyPart) int {
	if p.HasEquipped(bodyPart) {
		indexAdded := p.AddInventory(p.equipment[bodyPart])
		p.equipment[bodyPart] = nil
		return indexAdded
	} else {
		return -1
	}
}

func (p *Player) AddInventory(object *eq.EqObject) int {
	p.inventory = append(p.inventory, object)
	return len(p.inventory) - 1
}

func (p *Player) AddEffects(effects ...stats.StatusEffect) {
	p.effects = append(p.effects, effects...)
}

func (p *Player) GetEffects() []stats.StatusEffect {
	return p.effects
}

func (p *Player) Use(invIndex int) bool {
	eqObject := p.GetInventoryObject(invIndex)
	if eqObject == nil || !eqObject.IsUsable() {
		return false
	}

	p.AddEffects(eqObject.Effects...)
	p.destroy(invIndex)
	return true
}

func (p *Player) NextWeek() []stats.StatusEffect {
	var expiredEffects []stats.StatusEffect

	if p.effects != nil {
		var nextWeekEffects []stats.StatusEffect

		for _, effect := range p.effects {
			effect.Apply(p)
			// TODO: think through turn countdown more. Does effect expire at 0 or -1?
			// when using a potion should its benefits apply immediately? or next turn?

			if effect.HasExpired() {
				expiredEffects = append(expiredEffects, effect)
			} else {
				nextWeekEffects = append(nextWeekEffects, effect)
			}
		}

		p.effects = nextWeekEffects
	}

	return expiredEffects
}

func (p *Player) destroy(invIndex int) {
	p.inventory = slices.Delete(p.inventory, invIndex, invIndex+1)
}
