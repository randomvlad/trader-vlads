package event

import (
	"github.com/oklog/ulid/v2"
	eq "github.com/randomvlad/trader-vlads/internal/appmod/equipment"
	"github.com/randomvlad/trader-vlads/internal/appmod/keybind"
	eff "github.com/randomvlad/trader-vlads/internal/appmod/stats/statuseffect"
	"github.com/randomvlad/trader-vlads/internal/util"
	"github.com/randomvlad/trader-vlads/internal/util/stringutil"
)

type TurnKeeper struct {
	turn         int
	player       PlayerTurnService
	market       MarketService
	EventTracker *EventTracker
	toast        ToastMessenger
}

type PlayerTurnService interface {
	AddMoney(amount int)
	AddResourceQuantity(name string, quantity int)
	AddInventory(object *eq.EqObject) int
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

func NewTurnKeeper(player PlayerTurnService, market MarketService, keyBinder *keybind.KeyBinder, random *util.RandomGenerator, toast ToastMessenger) *TurnKeeper {
	eventTracker := NewEventTracker(player, keyBinder, random)
	return &TurnKeeper{
		turn:         1,
		player:       player,
		market:       market,
		EventTracker: eventTracker,
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

	events := t.EventTracker.GetRandomEvents()

	var toastMessage stringutil.Builder
	for _, effect := range expiredEffects {
		toastMessage.WriteLn(effect.GetMessageEnd())
	}

	for _, event := range events {
		toastMessage.
			Ln().
			WriteLn(event.Name).
			Ln().
			WriteLn(event.Description)
		t.player.AddMoney(event.Money)
		t.player.AddEffects(event.Effects...)
	}

	if toastMessage.IsNotBlank() {
		t.toast.Message(toastMessage.String())
	}
}

func (t *TurnKeeper) applyEffects() []eff.StatusEffect {
	var expiredIds []ulid.ULID
	var expiredEffects []eff.StatusEffect

	for _, effect := range t.player.GetEffects() {
		effect.Apply(t.player)

		if effect.HasEnded() {
			expiredIds = append(expiredIds, effect.Id())
			expiredEffects = append(expiredEffects, effect)
		}
	}

	t.player.RemoveEffects(expiredIds...)

	return expiredEffects
}
