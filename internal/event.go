package internal

import (
	"github.com/randomvlad/trader-vlads/internal/appmod/stats"
	"github.com/randomvlad/trader-vlads/internal/util"
)

type EventTracker struct {
	randomGenerator *util.RandomGenerator
}

type Event struct {
	name        string
	description string
	money       int
	affects     []*stats.StatusEffect
}

func NewEventTracker(r *util.RandomGenerator) *EventTracker {
	return &EventTracker{
		randomGenerator: r,
	}
}

func (t *EventTracker) getEvents() []*Event {
	return []*Event{
		{
			name:        "Modest Inheritance (+50)",
			description: "A distant relative has passed away and left you a modest sum of money.",
			money:       50,
		},
		{
			name:        "Unexpected Bills (-25)",
			description: "An unexpected expense has come up and must be taken care of.",
			money:       -25,
		},
		{
			name:        "Winds of Fortune (+100)",
			description: "You're in luck! An anonymous benefactor has donated to your cause.",
			money:       100,
		},
		{
			name:        "Bounty (-100)",
			description: "Tough times dictate tough measures. You place a bounty with The House of Ancients to eliminate a hostile rival.",
			money:       -100,
		},
	}
}

func (t *EventTracker) getRandomEvent() *Event {
	if t.randomGenerator.RollChance(0.25) {
		events := t.getEvents()
		index := t.randomGenerator.IntN(len(events))
		return events[index]
	} else {
		return nil
	}
}
