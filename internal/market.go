package internal

import (
	"slices"

	"github.com/randomvlad/trader-vlads/internal/util"
)

type Market struct {
	week        int
	items       map[string]*ResourceItem
	unlockCost  int
	lockedItems []*ResourceItem
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

func NewMarket() *Market {
	return &Market{
		week: 1,
		items: map[string]*ResourceItem{
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
		lockedItems: []*ResourceItem{
			NewResourceItem("Emerald", 40, 250),
			NewResourceItem("Diamond", 100, 200),
			NewResourceItem("Purple Elixir", 30, 400),
			NewResourceItem("Mithril", 50, 500),
			NewResourceItem("Dragon Scales", 500, 1000), // future idea: rare items have limited supply. Some turns available quantity will be low or zero.
			NewResourceItem("Dark Elixir", 100, 5000),
		},
	}
}

func NewResourceItem(name string, priceRangeMin, priceRangeMax int) *ResourceItem {
	item := &ResourceItem{
		name:              name,
		priceRangeMin:     priceRangeMin,
		priceRangeMax:     priceRangeMax,
		gainNegativeCap:   25,
		gainPositiveCap:   25,
		availableQuantity: -1,
	}

	item.priceCurrent = ((item.priceRangeMin + item.priceRangeMax) / 2) - 5 // TODO: simple starting point for now
	item.priceHistory = append(item.priceHistory, item.priceCurrent)

	return item
}

func (m *Market) NextWeek() {
	m.week = m.week + 1

	for _, item := range m.items {
		newPrice := util.RandomGain(item.priceCurrent, item.gainNegativeCap, item.gainPositiveCap, item.priceRangeMin, item.priceRangeMax)
		item.priceCurrent = newPrice
		item.priceHistory = append(item.priceHistory, newPrice)
	}
}

func (m *Market) GetPricesCurrent() map[string]int {
	currentPrices := make(map[string]int, len(m.items))

	for _, item := range m.items {
		currentPrices[item.name] = item.priceCurrent
	}

	return currentPrices
}

func (m *Market) GetPricePrevious(itemName string) int {
	item := m.items[itemName]
	if len(item.priceHistory) == 1 {
		return item.priceHistory[0]
	}

	return item.priceHistory[len(item.priceHistory)-2]
}

func (m *Market) UnlockItem() *ResourceItem {
	if len(m.lockedItems) == 0 {
		return nil
	}

	unlockedItem := m.lockedItems[0]
	m.lockedItems = slices.Delete(m.lockedItems, 0, 1)

	m.items[unlockedItem.name] = unlockedItem

	return unlockedItem
}
