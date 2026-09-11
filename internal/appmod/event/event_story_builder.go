package event

import (
	"github.com/randomvlad/trader-vlads/internal/appmod/keybind/appaction"
)

type StorySceneBuilder struct {
	name     string
	viewFunc func() string
	actions  []appaction.AppAction
}

func NewStorySceneBuilder(name string) *StorySceneBuilder {
	return &StorySceneBuilder{
		name: name,
	}
}

func (b *StorySceneBuilder) View(viewFunc func() string) *StorySceneBuilder {
	b.viewFunc = viewFunc
	return b
}

func (b *StorySceneBuilder) ViewStatic(staticContent string) *StorySceneBuilder {
	b.viewFunc = func() string {
		return staticContent
	}
	return b
}

func (b *StorySceneBuilder) AddAction(name string, executeFunc func()) *StorySceneBuilder {
	b.actions = append(b.actions, appaction.NewAppAction(name, executeFunc))
	return b
}

func (b *StorySceneBuilder) AddActionNextScene(name string, story Story, nextScene string) *StorySceneBuilder {
	return b.AddAction(name, func() {
		story.SetCurrentScene(nextScene)
	})
}

func (b *StorySceneBuilder) AddActionComplete(name string, story Story) *StorySceneBuilder {
	return b.AddAction(name, func() {
		story.MarkComplete()
	})
}

func (b *StorySceneBuilder) Build() StoryScene {
	return StoryScene{
		Name: b.name,
		View: b.viewFunc,
		GetActions: func() []appaction.AppAction {
			return b.actions
		},
	}
}
