package event

import "github.com/randomvlad/trader-vlads/internal/util"

type StoryInheritance struct {
	*BaseStory
}

func NewStoryInheritance(player PlayerTurnService) *StoryInheritance {
	story := &StoryInheritance{
		BaseStory: NewBaseStory("MeagerInheritance", "Meager Inheritance"),
	}
	story.initScenes(player)
	return story
}

func (s *StoryInheritance) initScenes(player PlayerTurnService) {

	const storyName = "Inheritance"

	scene := NewStorySceneBuilder(storyName).
		ViewStatic("A distant relative has passed away and left you a meager sum of "+
			util.FormatCurrency(50)+".").
		AddAction("Continue", func() {
			player.AddMoney(50)
			s.MarkComplete()
		}).
		Build()
	s.addScenes(scene)
}
