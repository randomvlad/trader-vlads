package event

import (
	"strings"

	"github.com/oklog/ulid/v2"
	eff "github.com/randomvlad/trader-vlads/internal/appmod/stats/statuseffect"
	"github.com/randomvlad/trader-vlads/internal/util"
)

type TurnKeeper struct {
	turn         int
	player       PlayerService
	market       MarketService
	eventTracker *EventTracker
	toast        ToastMessenger
}

type PlayerService interface {
	AddMoney(amount int)
	AddResourceQuantity(name string, quantity int)
	AddEffects(effects ...eff.StatusEffect)
	GetEffects() []eff.StatusEffect
	RemoveEffects(ids ...ulid.ULID)
}

type MarketService interface {
	MovePrices()
}

type ToastMessenger interface {
	Message(text string, a ...any)
	Clear()
}

func NewTurnKeeper(player PlayerService, market MarketService, random *util.RandomGenerator, toast ToastMessenger) *TurnKeeper {
	return &TurnKeeper{
		turn:         1,
		player:       player,
		market:       market,
		eventTracker: NewEventTracker(random),
		toast:        toast,
	}
}

func (t *TurnKeeper) GetTurn() int {
	return t.turn
}

func (t *TurnKeeper) Next() {
	t.turn++

	t.toast.Clear()

	t.market.MovePrices()

	expiredEffects := t.applyEffects()

	events := t.eventTracker.GetRandomEvents()

	var toastMessage strings.Builder
	for _, effect := range expiredEffects {
		toastMessage.WriteString(effect.GetMessageEnd() + "\n")
	}

	for _, event := range events {
		toastMessage.WriteString("\n" + event.Name + "\n\n" + event.Description + "\n")
		t.player.AddMoney(event.Money)
		t.player.AddEffects(event.Effects...)
	}

	if toastMessage.Len() > 0 {
		t.toast.Message(toastMessage.String())
	}
}

func (t *TurnKeeper) applyEffects() []eff.StatusEffect {
	var expiredIds []ulid.ULID
	var expiredEffects []eff.StatusEffect

	for _, effect := range t.player.GetEffects() {
		effect.Apply(t.player) // TODO: move turn counting out of effect object

		if effect.HasEnded() {
			expiredIds = append(expiredIds, effect.Id())
			expiredEffects = append(expiredEffects, effect)
		}
	}

	t.player.RemoveEffects(expiredIds...)

	return expiredEffects
}
