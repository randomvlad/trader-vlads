package player

import (
	"slices"

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

	equipped := eq.Forge.Make(
		r,
		"copper ring of a novice",
		"gray cotton tunic",
		"worn trousers",
		"brown leather sandals",
		"quill made from a talon of the Blue Dragon",
	)
	for _, eqObject := range equipped {
		player.Wear(eqObject)
	}

	inv := eq.Forge.Make(
		r,
		"a potion of Beginner's Luck 🍀",
		"a jar of spicy pickles",
		"Inexhaustible Cart of Lumber",
		"goose with feathers of pure gold",
	)
	for _, eqObject := range inv {
		player.AddInventory(eqObject)
	}

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

func (p *Player) AddEffects(effects ...eff.StatusEffect) {
	p.effects = append(p.effects, effects...)
}

func (p *Player) GetEffects() []eff.StatusEffect {

	var combined []eff.StatusEffect

	// equipment granted effects
	for _, eqObject := range p.equipped {
		if eqObject != nil {
			combined = append(combined, eqObject.Effects...)
		}
	}

	// other effects (example: granted by events)
	combined = append(combined, p.effects...)

	return combined
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

func (p *Player) NextWeek() []eff.StatusEffect {
	var expiredEffects []eff.StatusEffect

	for _, effect := range p.GetEffects() {
		// TODO: think through turn countdown more. Does effect expire at 0 or -1?
		// when using a potion should its benefits apply immediately? or next turn?
		effect.Apply(p)
	}

	// filter out expired events (equipment granted effects are permanent)
	if p.effects != nil {
		var nextWeekEffects []eff.StatusEffect

		for _, effect := range p.effects {
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
