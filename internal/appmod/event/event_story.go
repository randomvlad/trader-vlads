package event

import (
	"github.com/randomvlad/trader-vlads/internal/appmod/keybind/press"
)

type Story interface {
	GetName() string
	View() string
	MarkComplete()
	IsComplete() bool
	GetCurrentActions() []press.KeyAction
	GetActions() map[string][]press.KeyAction
	GetStateId() string
	SetCurrentScene(name string)
}

type StoryScene struct {
	Name       string
	View       func() string
	GetActions func() []press.KeyAction
}

type BaseStory struct {
	Story
	id           string
	name         string
	scenes       map[string]StoryScene
	currentScene string
	complete     bool
}

func NewBaseStory(id, name string) *BaseStory {

	baseStory := &BaseStory{
		id:     id,
		name:   name,
		scenes: make(map[string]StoryScene),
	}

	return baseStory
}

func (b *BaseStory) GetName() string {
	return b.name
}

func (b *BaseStory) View() string {
	return b.scenes[b.currentScene].View()
}

func (b *BaseStory) MarkComplete() {
	b.complete = true
}

func (b *BaseStory) IsComplete() bool {
	return b.complete
}

func (b *BaseStory) GetCurrentActions() []press.KeyAction {
	current := b.currentScene
	return b.scenes[current].GetActions()
}

func (b *BaseStory) GetActions() map[string][]press.KeyAction {
	mapActions := make(map[string][]press.KeyAction)
	for _, scene := range b.scenes {
		contextId := b.id + "_" + scene.Name
		mapActions[contextId] = scene.GetActions() // TODO: revisit idea of making this a func rather than array of elements
	}
	return mapActions
}

func (b *BaseStory) GetStateId() string {
	return b.id + "_" + b.currentScene
}

func (b *BaseStory) SetCurrentScene(name string) {
	b.currentScene = name
}

func (b *BaseStory) addScenes(scenes ...StoryScene) {
	for _, scene := range scenes {
		b.scenes[scene.Name] = scene
	}

	if b.currentScene == "" {
		b.currentScene = scenes[0].Name
	}
}
