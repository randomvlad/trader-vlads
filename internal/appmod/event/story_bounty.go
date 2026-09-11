package event

import "github.com/randomvlad/trader-vlads/internal/util"

type StoryBounty struct {
	*BaseStory
}

func NewStoryBounty(player PlayerTurnService) *StoryBounty {
	story := &StoryBounty{
		BaseStory: NewBaseStory("UnrulyRival", "Unruly Rival"),
	}
	story.initScenes(player)
	return story
}

func (s *StoryBounty) initScenes(player PlayerTurnService) {

	const storyName = "UnrulyRival"

	scene := NewStorySceneBuilder(storyName).
		ViewStatic("Tough times dictate tough measures. You place a "+util.FormatCurrency(100)+
			" bounty with The House of Ancients to eliminate a hostile rival.").
		AddAction("Continue", func() {
			player.AddMoney(-100)
			s.MarkComplete()
		}).
		Build()
	s.addScenes(scene)
}
