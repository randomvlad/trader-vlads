package internal

type Market struct {
	week        int
	items       map[string]ResourceItem
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

	//lockedItems: []string{
	//	"Emerald",
	//	"Diamond",
	//	"Purple Elixir",
	//	"Dark Elixir",
	//	"Mithril",
	//	"Adamantine",
	//	"Elven Silk",
	//	"Dragon Scales", // future idea: rare items have limited supply. Some turns available quantity will be low or zero.
	//},

	return &Market{
		week: 1,
		items: map[string]ResourceItem{
			"Wood":    NewResourceItem("Wood", 10, 50),
			"Stone":   NewResourceItem("Stone", 10, 40),
			"Coal":    NewResourceItem("Coal", 10, 50),
			"Grain":   NewResourceItem("Grain", 15, 60),
			"Cloth":   NewResourceItem("Cloth", 15, 60),
			"Leather": NewResourceItem("Leather", 15, 80),
			"Iron":    NewResourceItem("Iron", 20, 100),
			"Glass":   NewResourceItem("Glass", 30, 200),
		},
		unlockCost:  100,
		lockedItems: []*ResourceItem{}, // TODO: tackle later
	}
}

func NewResourceItem(name string, priceRangeMin, priceRangeMax int) ResourceItem {
	item := ResourceItem{
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
	// TODO: update prices
}

func (m *Market) GetCurrentPrices() map[string]int {
	currentPrices := make(map[string]int, len(m.items))

	for _, item := range m.items {
		currentPrices[item.name] = item.priceCurrent
	}

	return currentPrices
}

//func getTurnData(player *Player) *TurnData {
//
//	randomItemIndices := rand.Perm(len(player.activeItems))[:5]
//
//	marketPrices := make(map[string]int)
//	for _, randomIndex := range randomItemIndices {
//		itemName := player.activeItems[randomIndex]
//		itemPrice := rand.IntN(41) + 10
//
//		marketPrices[itemName] = itemPrice
//	}
//
//	marketItems := slices.Sorted(maps.Keys(marketPrices))
//
//	return &TurnData{
//
//	}
//}
