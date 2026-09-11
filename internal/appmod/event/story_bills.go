package event

import "github.com/randomvlad/trader-vlads/internal/util"

type StoryBills struct {
	*BaseStory
}

func NewStoryBills(player PlayerTurnService) *StoryBills {
	story := &StoryBills{
		BaseStory: NewBaseStory("UnexpectedBills", "Unexpected Bills"),
	}
	story.initScenes(player)
	return story
}

func (s *StoryBills) initScenes(player PlayerTurnService) {

	const storyName = "UnexpectedBills"

	scene := NewStorySceneBuilder(storyName).
		ViewStatic("An unexpected expense has come up and must be taken care of. "+
			"You begrudgingly pay "+util.FormatCurrency(25)+".").
		AddAction("Continue", func() {
			player.AddMoney(-25)
			s.MarkComplete()
		}).
		Build()
	s.addScenes(scene)
}
