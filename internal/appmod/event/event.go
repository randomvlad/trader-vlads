package event

import (
	"github.com/oklog/ulid/v2"
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
	EffectDefs  []eff.StatusEffectDef
	Effects     []eff.StatusEffect
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
		{
			Name:        "Blessings of Evergreen",
			Description: "The forest nymphs of Evergreen have bestowed their blessings upon you.",
			EffectDefs: []eff.StatusEffectDef{
				&eff.GrantResourceEffectDef{
					Name:     "Blessings of Evergreen",
					Resource: "Wood",
					Amount:   util.NewRangeInt(1, 1),
					Turns:    util.NewRangeInt(2, 4),
				},
			},
		},
	}
}

func (t *EventTracker) GetRandomEvent() *Event {
	if t.randomGenerator.RollChance(25) {
		event := t.randomGenerator.Pick(t.GetEvents())

		var effects []eff.StatusEffect
		for _, def := range event.EffectDefs {
			effects = append(effects, def.Create(t.randomGenerator, ulid.Make(), "event"))
		}
		event.Effects = effects
		return event
	} else {
		return nil
	}
}
