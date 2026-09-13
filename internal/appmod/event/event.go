package event

import (
	"github.com/randomvlad/trader-vlads/internal/appmod/keybind"
	"github.com/randomvlad/trader-vlads/internal/util"
)

type EventTracker struct {
	keyBinder       *keybind.KeyBinder
	player          PlayerTurnService
	randomGenerator *util.RandomGenerator
	activeStory     Story
}

func NewEventTracker(player PlayerTurnService, keyBinder *keybind.KeyBinder, r *util.RandomGenerator) *EventTracker {

	tracker := &EventTracker{
		keyBinder:       keyBinder,
		player:          player,
		randomGenerator: r,
	}

	// TODO: corresponds to start of the game (turn 0). introduce a dedicated "new turn" phase
	tracker.setActive(NewStoryNewBeginnings(player, r))
	return tracker
}

func (t *EventTracker) GetActiveStory() Story {
	if t.activeStory != nil && !t.activeStory.IsComplete() {
		return t.activeStory
	} else {
		return nil
	}
}

func (t *EventTracker) GenerateActiveStory(turn int) {
	// for now limiting to 1 event per method call/turn

	if turn == 2 { // TODO: introduce a schedule later on.
		t.setActive(NewStoryEvergreenBlessings(t.player, t.randomGenerator))
	} else if turn == 3 {
		t.setActive(NewStoryFortune(t.player))
	} else {
		if t.randomGenerator.RollChance(25) {
			event := t.randomGenerator.Pick(t.getStories())
			t.setActive(event)
		}
	}
}

func (t *EventTracker) setActive(story Story) {
	t.activeStory = story
	t.keyBinder.AddActionsMap(story.GetActions())
}

func (t *EventTracker) getStories() []Story {
	// TODO: introduce a story registry. Create stories more efficiently. Pick from a list of strings. Create/generate only if story is picked.
	return []Story{
		NewStoryBills(t.player),
		NewStoryInheritance(t.player),
		NewStoryBounty(t.player),
	}
}
