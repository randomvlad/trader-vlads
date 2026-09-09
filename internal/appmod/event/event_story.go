package event

type Story interface {
	GetName() string
	View() string
	IsComplete() bool
	GetActions() []AppAction
	GetStateId() string
}

type StoryScene struct {
	Name       string
	View       func() string
	GetActions func() []AppAction
}

type AppAction struct {
	Name        string
	KeyPress    string
	executeFunc func()
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

func (b *BaseStory) IsComplete() bool {
	return b.complete
}

func (b *BaseStory) GetActions() []AppAction {
	current := b.currentScene
	return b.scenes[current].GetActions()
}

func (b *BaseStory) GetStateId() string {
	return b.id + "_" + b.currentScene
}
