package event

import "github.com/randomvlad/trader-vlads/internal/util"

type StoryFortune struct {
	*BaseStory
}

func NewStoryFortune(player PlayerTurnService) *StoryFortune {
	story := &StoryFortune{
		BaseStory: NewBaseStory("Fortune", "Winds of Fortune"),
	}
	story.initScenes(player)
	return story
}

func (s *StoryFortune) initScenes(player PlayerTurnService) {

	const storyName = "Fortune"

	scene := NewStorySceneBuilder(storyName).
		ViewStatic("You're in luck! An anonymous benefactor donates "+
			util.FormatCurrency(100)+" to your cause.").
		AddAction("Continue", func() {
			player.AddMoney(100)
			s.MarkComplete()
		}).
		Build()
	s.addScenes(scene)
}
