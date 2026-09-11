package event

import (
	"strconv"

	eq "github.com/randomvlad/trader-vlads/internal/appmod/equipment"
	"github.com/randomvlad/trader-vlads/internal/util"
	"github.com/randomvlad/trader-vlads/internal/util/stringutil"
)

type StoryNewBeginnings struct {
	*BaseStory
}

func NewStoryNewBeginnings(player PlayerTurnService, r *util.RandomGenerator) *StoryNewBeginnings {
	story := &StoryNewBeginnings{
		BaseStory: NewBaseStory("NewBeginnings", "To New Beginnings"),
	}
	story.initScenes(player, r)
	return story
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

	const (
		SceneNameStart      = "Start"
		SceneNameChestOpen  = "ChestOpen"
		SceneNameItemsAdded = "ItemsAdded"
	)

	sceneStart := NewStorySceneBuilder(SceneNameStart).
		ViewStatic("The Guild of Merchants has sent a standard edition wooden chest to get you started.").
		AddActionNextScene("Open", s, SceneNameChestOpen).
		Build()

	sceneChestOpen := NewStorySceneBuilder(SceneNameChestOpen).
		View(func() string {
			var render stringutil.Builder
			render.WriteLn("You open the chest and look inside:")

			for resource, count := range resources {
				render.Writef("    +%s %s\n", strconv.Itoa(count), resource)
			}

			for _, itemName := range items {
				render.Tab().WriteLn(itemName)
			}

			return render.String()
		}).
		AddAction("Take", func() {
			for resource, count := range resources {
				player.AddResourceQuantity(resource, count)
			}

			eqObjects := eq.Forge.Make(r, items...)
			for _, object := range eqObjects {
				player.AddInventory(object)
			}

			s.SetCurrentScene(SceneNameItemsAdded)
		}).
		Build()

	sceneItemsAdded := NewStorySceneBuilder(SceneNameItemsAdded).
		ViewStatic("You place the items in your inventory. Might be a good idea to try them on next.").
		AddActionComplete("Continue", s).
		Build()

	s.addScenes(sceneStart, sceneChestOpen, sceneItemsAdded)
}
