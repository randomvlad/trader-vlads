package internal

import (
	"slices"

	"github.com/randomvlad/trader-vlads/internal/util"
)

type Market struct {
	week            int
	resources       map[string]*ResourceItem
	unlockCost      int
	lockedResources []*ResourceItem
	random          *util.RandomGenerator
}

type ResourceItem struct {
	name              string
	priceCurrent      int
	priceHistory      []int
	priceRangeMin     int
	priceRangeMax     int
	gainNegativeCap   int
	gainPositiveCap   int
	availableQuantity int
}

func NewMarket(r *util.RandomGenerator) *Market {
	return &Market{
		week: 1,
		resources: map[string]*ResourceItem{
			"Wood":    NewResourceItem("Wood", 10, 50),
			"Stone":   NewResourceItem("Stone", 10, 40),
			"Coal":    NewResourceItem("Coal", 10, 50),
			"Grain":   NewResourceItem("Grain", 15, 60),
			"Cloth":   NewResourceItem("Cloth", 15, 60),
			"Leather": NewResourceItem("Leather", 15, 80),
			"Iron":    NewResourceItem("Iron", 20, 100),
			"Glass":   NewResourceItem("Glass", 30, 200),
		},
		unlockCost: 100,
		lockedResources: []*ResourceItem{
			NewResourceItem("Emerald", 40, 250),
			NewResourceItem("Diamond", 100, 200),
			NewResourceItem("Purple Elixir", 30, 400),
			NewResourceItem("Mithril", 50, 500),
			NewResourceItem("Dragon Scales", 500, 1000), // future idea: rare resources have limited supply. Some turns available quantity will be low or zero.
			NewResourceItem("Dark Elixir", 100, 5000),
		},
		random: r,
	}
}

func NewResourceItem(name string, priceRangeMin, priceRangeMax int) *ResourceItem {
	item := &ResourceItem{
		name:              name,
		priceRangeMin:     priceRangeMin,
		priceRangeMax:     priceRangeMax,
		gainNegativeCap:   25,
		gainPositiveCap:   32, // experiment: throw in 7% boost for upward mobility
		availableQuantity: -1,
	}

	item.priceCurrent = ((item.priceRangeMin + item.priceRangeMax) / 2) - 5 // TODO: simple starting point for now
	item.priceHistory = append(item.priceHistory, item.priceCurrent)

	return item
}

func (m *Market) NextWeek() {
	m.week = m.week + 1

	for _, r := range m.resources {
		newPrice := m.random.RandomGain(r.priceCurrent, r.gainNegativeCap, r.gainPositiveCap, r.priceRangeMin, r.priceRangeMax)
		r.priceCurrent = newPrice
		r.priceHistory = append(r.priceHistory, newPrice)
	}
}

func (m *Market) GetPricesCurrent() map[string]int {
	currentPrices := make(map[string]int, len(m.resources))

	for _, item := range m.resources {
		currentPrices[item.name] = item.priceCurrent
	}

	return currentPrices
}

func (m *Market) GetPricePrevious(resourceName string) int {
	resource := m.resources[resourceName]
	if len(resource.priceHistory) == 1 {
		return resource.priceHistory[0]
	}

	return resource.priceHistory[len(resource.priceHistory)-2]
}

func (m *Market) UnlockResource() *ResourceItem {
	if len(m.lockedResources) == 0 {
		return nil
	}

	unlockedItem := m.lockedResources[0]
	m.lockedResources = slices.Delete(m.lockedResources, 0, 1)

	m.resources[unlockedItem.name] = unlockedItem

	return unlockedItem
}
