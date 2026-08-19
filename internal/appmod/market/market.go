package market

import (
	"slices"

	"github.com/randomvlad/trader-vlads/internal/util"
)

type Market struct {
	// TODO: revisit public vs private scope
	Week            int
	Resources       map[string]*ResourceItem
	UnlockCost      int
	LockedResources []*ResourceItem
	random          *util.RandomGenerator
}

type ResourceItem struct {
	Name              string
	PriceCurrent      int
	PriceHistory      []int
	PriceRangeMin     int
	PriceRangeMax     int
	GainNegativeCap   int
	GainPositiveCap   int
	AvailableQuantity int
}

func NewMarket(r *util.RandomGenerator) *Market {
	return &Market{
		Week: 1,
		Resources: map[string]*ResourceItem{
			"Wood":    NewResourceItem("Wood", 10, 50),
			"Stone":   NewResourceItem("Stone", 10, 40),
			"Coal":    NewResourceItem("Coal", 10, 50),
			"Grain":   NewResourceItem("Grain", 15, 60),
			"Cloth":   NewResourceItem("Cloth", 15, 60),
			"Leather": NewResourceItem("Leather", 15, 80),
			"Iron":    NewResourceItem("Iron", 20, 100),
			"Glass":   NewResourceItem("Glass", 30, 200),
		},
		UnlockCost: 100,
		LockedResources: []*ResourceItem{
			NewResourceItem("Emerald", 40, 250),
			NewResourceItem("Diamond", 100, 200),
			NewResourceItem("Purple Elixir", 30, 400),
			NewResourceItem("Mithril", 50, 500),
			NewResourceItem("Dragon Scales", 500, 1000), // future idea: rare Resources have limited supply. Some turns available quantity will be low or zero.
			NewResourceItem("Dark Elixir", 100, 5000),
		},
		random: r,
	}
}

func NewResourceItem(name string, priceRangeMin, priceRangeMax int) *ResourceItem {
	item := &ResourceItem{
		Name:              name,
		PriceRangeMin:     priceRangeMin,
		PriceRangeMax:     priceRangeMax,
		GainNegativeCap:   25,
		GainPositiveCap:   32, // experiment: throw in 7% boost for upward mobility
		AvailableQuantity: -1,
	}

	item.PriceCurrent = ((item.PriceRangeMin + item.PriceRangeMax) / 2) - 5 // TODO: simple starting point for now
	item.PriceHistory = append(item.PriceHistory, item.PriceCurrent)

	return item
}

func (m *Market) NextWeek() {
	m.Week = m.Week + 1

	for _, r := range m.Resources {
		newPrice := m.random.GainPercent(r.PriceCurrent, r.GainNegativeCap, r.GainPositiveCap, r.PriceRangeMin, r.PriceRangeMax)
		r.PriceCurrent = newPrice
		r.PriceHistory = append(r.PriceHistory, newPrice)
	}
}

func (m *Market) GetPricesCurrent() map[string]int {
	currentPrices := make(map[string]int, len(m.Resources))

	for _, item := range m.Resources {
		currentPrices[item.Name] = item.PriceCurrent
	}

	return currentPrices
}

func (m *Market) GetPricePrevious(resourceName string) int {
	resource := m.Resources[resourceName]
	if len(resource.PriceHistory) == 1 {
		return resource.PriceHistory[0]
	}

	return resource.PriceHistory[len(resource.PriceHistory)-2]
}

func (m *Market) UnlockResource() *ResourceItem {
	if len(m.LockedResources) == 0 {
		return nil
	}

	unlockedItem := m.LockedResources[0]
	m.LockedResources = slices.Delete(m.LockedResources, 0, 1)

	m.Resources[unlockedItem.Name] = unlockedItem

	return unlockedItem
}
