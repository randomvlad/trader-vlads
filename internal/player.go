package internal

type Player struct {
	money     int
	inventory map[string]int
	actions   []string
}

func NewPlayer(m *Market) *Player {

	player := Player{
		money: 1000,
		actions: []string{
			"Buy",
			"Sell",
			"Next Week",
			"Unlock Item",
			"Quit",
		},
	}

	inventory := make(map[string]int)
	for _, item := range m.items {
		inventory[item.name] = 0
	}
	player.inventory = inventory

	return &player
}

func (p *Player) isInventoryEmpty() bool {
	for _, count := range p.inventory {
		if count > 0 {
			return false
		}
	}
	return true
}
