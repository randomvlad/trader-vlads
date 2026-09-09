package event

import (
	"strconv"

	tea "charm.land/bubbletea/v2"
	eq "github.com/randomvlad/trader-vlads/internal/appmod/equipment"
	"github.com/randomvlad/trader-vlads/internal/appmod/keybind"
	"github.com/randomvlad/trader-vlads/internal/util"
	"github.com/randomvlad/trader-vlads/internal/util/stringutil"
)

type StoryNewBeginnings struct {
	*BaseStory
}

func NewStoryNewBeginnings(
	name string,
	player PlayerTurnService,
	keyBinder *keybind.KeyBinder,
	r *util.RandomGenerator,
) *StoryNewBeginnings {

	story := &StoryNewBeginnings{
		BaseStory: NewBaseStory("NewBeginnings", name),
	}

	story.initScenes(player, r)

	story.bindKeys(keyBinder)

	return story
}

// TODO: move into keybinder
func (s *StoryNewBeginnings) bindKeys(keyBinder *keybind.KeyBinder) {
	for _, scene := range s.scenes {
		for _, action := range scene.GetActions() {
			keyBinder.Add(s.id+"_"+scene.Name, action.KeyPress, func() tea.Cmd {
				action.executeFunc()
				return nil
			})
		}
	}
}

func (s *StoryNewBeginnings) initScenes(player PlayerTurnService, r *util.RandomGenerator) {
	items := []string{
		"copper ring of a novice",
		"gray cotton tunic",
		"worn trousers",
		"brown leather sandals",
		"a potion of Beginner's Luck 🍀",
		"a jar of spicy pickles",
	}

	resources := map[string]int{
		"Wood":  3,
		"Stone": 3,
	}

	scenes := []StoryScene{
		{
			Name: "SceneStart",
			View: func() string {
				return "The Guild of Merchants has sent a standard edition wooden chest to get you started."
			},
			GetActions: func() []AppAction {
				return []AppAction{
					{
						Name:     "Open",
						KeyPress: "O",
						executeFunc: func() {
							s.currentScene = "SceneChestOpen"
						},
					},
				}
			},
		},
		{
			Name: "SceneChestOpen",
			View: func() string {
				var render stringutil.Builder
				render.WriteLn("You open the chest and look inside:")

				for resource, count := range resources {
					render.Writef("    +%s %s\n", strconv.Itoa(count), resource)
				}

				for _, itemName := range items {
					render.Tab().WriteLn(itemName)
				}

				return render.String()
			},
			GetActions: func() []AppAction {
				return []AppAction{
					{
						Name:     "Take",
						KeyPress: "T",
						executeFunc: func() {
							for resource, count := range resources {
								player.AddResourceQuantity(resource, count)
							}

							eqObjects := eq.Forge.Make(r, items...)
							for _, object := range eqObjects {
								player.AddInventory(object)
							}

							s.currentScene = "SceneItemsAdded"
						},
					},
				}
			},
		},
		{
			Name: "SceneItemsAdded",
			View: func() string {
				return "You place the items in your inventory. Might be a good idea to try them on next."
			},
			GetActions: func() []AppAction {
				return []AppAction{
					{
						Name:     "Continue",
						KeyPress: "C",
						executeFunc: func() {
							s.complete = true
						},
					},
				}
			},
		},
	}

	for _, scene := range scenes {
		s.scenes[scene.Name] = scene
	}
	s.currentScene = scenes[0].Name
}
