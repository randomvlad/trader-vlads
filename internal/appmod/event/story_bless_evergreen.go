package event

import (
	"github.com/oklog/ulid/v2"
	eff "github.com/randomvlad/trader-vlads/internal/appmod/stats/statuseffect"
	"github.com/randomvlad/trader-vlads/internal/util"
	"github.com/randomvlad/trader-vlads/internal/util/stringutil"
)

type StoryEvergreenBlessings struct {
	*BaseStory
	effect eff.StatusEffect
}

func NewStoryEvergreenBlessings(player PlayerTurnService, r *util.RandomGenerator) *StoryEvergreenBlessings {
	story := &StoryEvergreenBlessings{
		BaseStory: NewBaseStory("EvergreenBlessings", "Blessings of Evergreen"),
	}
	story.initScenes(player, r)
	return story
}

func (s *StoryEvergreenBlessings) initScenes(player PlayerTurnService, r *util.RandomGenerator) {

	const (
		SceneNameStart    = "Start"
		SceneNameAccepted = "Accepted"
		SceneNameDeclined = "Declined"
	)

	sceneStart := NewStorySceneBuilder(SceneNameStart).
		ViewStatic("The forest nymphs of Evergreen approach you and offer their blessings.").
		AddAction("Accept", func() {
			effectDef := &eff.GrantResourceEffectDef{
				BaseEffectDef: &eff.BaseEffectDef{
					Name:     "Blessings of Evergreen",
					Duration: eff.NewDuration().Turns(2, 4),
				},
				Resource: "Wood",
				Amount:   util.NewRangeInt(1, 1),
			}

			s.effect = effectDef.Create(r, ulid.Make(), "event")
			player.AddEffects(s.effect)
			s.SetCurrentScene(SceneNameAccepted)

		}).AddActionNextScene("Decline", s, SceneNameDeclined).
		Build()

	sceneAccepted := NewStorySceneBuilder(SceneNameAccepted).
		View(func() string {
			return stringutil.NewBuilder().
				WriteLn("The forest nymphs of Evergreen bestow their blessings upon you.").
				Ln().
				WriteLn("A faint voice whispers: \"" + s.effect.GetMessageStart() + "\" 🍃").
				String()
		}).
		AddActionComplete("Continue", s).
		Build()

	sceneDeclined := NewStorySceneBuilder(SceneNameDeclined).
		ViewStatic("The forest nymphs of Evergreen nod solemnly and disappear. An earthy scent of bark and moss lingers.").
		AddActionComplete("Continue", s).
		Build()

	s.addScenes(sceneStart, sceneAccepted, sceneDeclined)
}
