package player

import (
	"slices"

	"github.com/oklog/ulid/v2"
	eq "github.com/randomvlad/trader-vlads/internal/appmod/equipment"
	"github.com/randomvlad/trader-vlads/internal/appmod/market"
	eff "github.com/randomvlad/trader-vlads/internal/appmod/stats/statuseffect"
	"github.com/randomvlad/trader-vlads/internal/util"
)

type Player struct {
	money     int
	Warehouse *Warehouse
	equipped  map[eq.BodyPart]*eq.EqObject
	inventory []*eq.EqObject
	effects   []eff.StatusEffect
	random    *util.RandomGenerator
}

func NewPlayer(m *market.Market, r *util.RandomGenerator) *Player {

	warehouse := &Warehouse{
		Capacity:  100,
		Resources: make(map[string]int),
	}

	player := Player{
		money:     1000,
		Warehouse: warehouse,
		equipped:  make(map[eq.BodyPart]*eq.EqObject),
		inventory: []*eq.EqObject{},
		random:    r,
	}

	for _, item := range m.Resources {
		player.Warehouse.Resources[item.Name] = 0
	}

	// initialize every eq body part
	for bodyPart := range eq.BodyPartsMax {
		player.equipped[eq.BodyPart(bodyPart)] = nil
	}

	//equipped := eq.Forge.Make(
	//	r,
	//	"quill made from a talon of the Blue Dragon",
	//)
	//for _, eqObject := range equipped {
	//	player.Wear(eqObject)
	//}
	//
	//inv := eq.Forge.Make(
	//	r,
	//	"Inexhaustible Cart of Lumber",
	//	"goose with feathers of pure gold",
	//)
	//for _, eqObject := range inv {
	//	player.AddInventory(eqObject)
	//}

	return &player
}

type Warehouse struct {
	Capacity  int
	Resources map[string]int // Resource name → quantity
}

func (p *Player) GetMoney() int {
	return p.money
}

func (p *Player) AddMoney(amount int) {
	p.money += amount
}

func (p *Player) AddResourceQuantity(name string, quantity int) {
	p.Warehouse.Resources[name] += quantity
}

func (p *Player) IsWarehouseEmpty() bool {
	for _, count := range p.Warehouse.Resources {
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
	return p.equipped
}

func (p *Player) GetEquippedObject(bodyPart eq.BodyPart) *eq.EqObject {
	return p.equipped[bodyPart]
}

func (p *Player) Wear(object *eq.EqObject) {
	var bodyPart eq.BodyPart

	switch object.Slot {
	case eq.EqSlotHead:
		bodyPart = eq.BodyPartHead
	case eq.EqSlotNeck:
		bodyPart = eq.BodyPartNeck
	case eq.EqSlotTorso:
		bodyPart = eq.BodyPartTorso
	case eq.EqSlotHands:
		bodyPart = eq.BodyPartHands
	case eq.EqSlotFinger:
		if p.HasEquipped(eq.BodyPartFingerLeft) {
			bodyPart = eq.BodyPartFingerRight
		} else {
			bodyPart = eq.BodyPartFingerLeft
		}
	case eq.EqSlotWaist:
		bodyPart = eq.BodyPartWaist
	case eq.EqSlotLegs:
		bodyPart = eq.BodyPartLegs
	case eq.EqSlotFeet:
		bodyPart = eq.BodyPartFeet
	case eq.EqSlotWield, eq.EqSlotHold:
		if p.HasEquipped(eq.BodyPartHoldLeft) {
			bodyPart = eq.BodyPartHoldRight
		} else {
			bodyPart = eq.BodyPartHoldLeft
		}
	default:
		return
	}

	p.wearOnBody(bodyPart, object)
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

func (p *Player) wearOnBody(part eq.BodyPart, object *eq.EqObject) {
	if object != nil {
		p.Remove(part)
		p.equipped[part] = object
	}
}

func (p *Player) HasEquipped(bodyPart eq.BodyPart) bool {
	return p.equipped[bodyPart] != nil
}

func (p *Player) Remove(bodyPart eq.BodyPart) int {
	if p.HasEquipped(bodyPart) {
		indexAdded := p.AddInventory(p.equipped[bodyPart])
		p.equipped[bodyPart] = nil
		return indexAdded
	} else {
		return -1
	}
}

func (p *Player) AddInventory(object *eq.EqObject) int {
	p.inventory = append(p.inventory, object)
	return len(p.inventory) - 1
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

func (p *Player) AddEffects(effects ...eff.StatusEffect) {
	for _, effect := range effects {
		if p.HasEffect(effect.Name()) {
			continue
		}

		p.effects = append(p.effects, effect)
	}
}

func (p *Player) HasEffect(name string) bool {
	// Note: revisit effect matching later on. Most effects have a unique name even if they provide the same type of grant
	return slices.ContainsFunc(p.GetEffects(), func(effect eff.StatusEffect) bool {
		return effect.Name() == name
	})
}

func (p *Player) GetEffects() []eff.StatusEffect {
	var eqEffects []eff.StatusEffect

	// equipment granted effects
	for _, eqObject := range p.equipped {
		if eqObject != nil {
			eqEffects = append(eqEffects, eqObject.Effects...)
		}
	}

	// other effects (example: granted by events)
	return append(eqEffects, p.effects...)
}

func (p *Player) RemoveEffects(ids ...ulid.ULID) {
	p.effects = slices.DeleteFunc(p.effects, func(effect eff.StatusEffect) bool {
		return slices.Contains(ids, effect.Id())
	})
}

func (p *Player) destroy(invIndex int) {
	p.inventory = slices.Delete(p.inventory, invIndex, invIndex+1)
}
