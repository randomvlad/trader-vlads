package event

import (
	eff "github.com/randomvlad/trader-vlads/internal/appmod/stats/statuseffect"
	"github.com/randomvlad/trader-vlads/internal/util"
)

type EventTracker struct {
	randomGenerator *util.RandomGenerator
}

type Event struct {
	Name        string
	Description string
	Money       int
	Effects     []*eff.StatusEffect
}

func NewEventTracker(r *util.RandomGenerator) *EventTracker {
	return &EventTracker{
		randomGenerator: r,
	}
}

func (t *EventTracker) GetEvents() []*Event {
	return []*Event{
		{
			Name:        "Modest Inheritance (+50)",
			Description: "A distant relative has passed away and left you a modest sum of money.",
			Money:       50,
		},
		{
			Name:        "Unexpected Bills (-25)",
			Description: "An unexpected expense has come up and must be taken care of.",
			Money:       -25,
		},
		{
			Name:        "Winds of Fortune (+100)",
			Description: "You're in luck! An anonymous benefactor has donated to your cause.",
			Money:       100,
		},
		{
			Name:        "Bounty (-100)",
			Description: "Tough times dictate tough measures. You place a bounty with The House of Ancients to eliminate a hostile rival.",
			Money:       -100,
		},
	}
}

func (t *EventTracker) GetRandomEvent() *Event {
	if t.randomGenerator.RollChance(0.25) {
		events := t.GetEvents()
		index := t.randomGenerator.IntN(len(events))
		return events[index]
	} else {
		return nil
	}
}
