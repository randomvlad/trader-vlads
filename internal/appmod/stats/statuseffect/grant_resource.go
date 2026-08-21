package statuseffect

import (
	"strconv"

	"github.com/oklog/ulid/v2"
	"github.com/randomvlad/trader-vlads/internal/util"
)

type GrantResourceEffectDef struct {
	Name      string
	Resource  string
	AmountMin int
	AmountMax int
	TurnsMin  int
	TurnsMax  int
}

func (d *GrantResourceEffectDef) Create(r *util.RandomGenerator, grantById ulid.ULID, grandByType string) StatusEffect {
	return &GrantResourceEffect{
		id:          ulid.Make(),
		name:        d.Name,
		grantById:   grantById,
		grantByType: grandByType,
		turnsLeft:   r.RollInt(d.TurnsMin, d.TurnsMax),
		resource:    d.Resource,
		amount:      r.RollInt(d.AmountMin, d.AmountMax),
	}
}

type GrantResourceEffect struct {
	id          ulid.ULID
	name        string
	grantById   ulid.ULID
	grantByType string
	turnsLeft   int
	resource    string
	amount      int
}

func (g *GrantResourceEffect) Id() ulid.ULID {
	return g.id
}

func (g *GrantResourceEffect) Name() string {
	return g.name
}

func (g *GrantResourceEffect) GrantedById() ulid.ULID {
	return g.grantById
}

func (g *GrantResourceEffect) GrantedByType() string {
	return g.grantByType
}

func (g *GrantResourceEffect) View() string {
	return "Grants +" + strconv.Itoa(g.amount) + " " + g.resource + " " + g.viewDurationLeft()
}

func (g *GrantResourceEffect) GetAppliedMessage() string {
	return "May the Wood be with you. Always!"
}

func (g *GrantResourceEffect) GetExpiredMessage() string {
	// TODO: does not apply if permanent; also messages are tied closely to item/event causing the effect
	return "The power of Wood has left you ..."
}

func (g *GrantResourceEffect) HasExpired() bool {
	if g.isPermanent() {
		return false
	} else {
		return g.turnsLeft <= 0
	}
}

func (g *GrantResourceEffect) Apply(player PlayerEffectService) {
	player.AddResourceQuantity(g.resource, g.amount)
	if !g.isPermanent() {
		g.turnsLeft--
	}
}

func (g *GrantResourceEffect) viewDurationLeft() string {
	if g.isPermanent() {
		return "each week"
	} else {
		var expiresSoonIcon string
		if g.turnsLeft == 1 {
			expiresSoonIcon = " ⌛"
		}

		return "for " + util.FormatCountPluralized(g.turnsLeft, "week") + expiresSoonIcon
	}
}

func (g *GrantResourceEffect) isPermanent() bool {
	return g.turnsLeft == -1
}
