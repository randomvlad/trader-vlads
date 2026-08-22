package statuseffect

import (
	"github.com/oklog/ulid/v2"
	"github.com/randomvlad/trader-vlads/internal/util"
)

type StatusEffect interface {
	Id() ulid.ULID
	Name() string
	GrantedById() ulid.ULID
	GrantedByType() string
	View() string
	GetMessageStart() string
	GetMessageEnd() string
	HasEnded() bool
	Apply(player PlayerEffectService)
}

type StatusEffectDef interface {
	Create(r *util.RandomGenerator, grantById ulid.ULID, grandByType string) StatusEffect
}

type PlayerEffectService interface {
	AddMoney(amount int)
	AddResourceQuantity(name string, quantity int)
}
