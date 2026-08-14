package equipment

import (
	"fmt"
	"strconv"
	"strings"

	tea "charm.land/bubbletea/v2"
)

type Model struct {
	player PlayerService
}

func NewTuiModel(player PlayerService) *Model {
	return &Model{
		player: player,
	}
}

type PlayerService interface {
	GetInventory() []*EqObject
	HasEquipment(BodyPart) bool
	GetEquipment() map[BodyPart]*EqObject
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) View() tea.View {
	var content strings.Builder

	content.WriteString(m.viewEqInv())

	return tea.NewView(content.String())
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	return m, tea.Batch(cmds...)
}

func (m *Model) viewEqInv() string {

	var view strings.Builder

	view.WriteString("You are using:\n")

	eq := m.player.GetEquipment()
	if len(eq) > 0 {
		for bodyPart := range BodyPartHoldRight {
			eqObject, hasEqOnBodyPart := eq[bodyPart]
			if !hasEqOnBodyPart {
				continue
			}

			var wornOn string
			switch bodyPart {
			case BodyPartFingerLeft, BodyPartFingerRight:
				wornOn = "worn on finger"
			case BodyPartNeck:
				wornOn = "worn around neck"
			case BodyPartTorso:
				wornOn = "worn on torso"
			case BodyPartHead:
				wornOn = "worn on head"
			case BodyPartLegs:
				wornOn = "worn on legs"
			case BodyPartFeet:
				wornOn = "worn on feet"
			case BodyPartHands:
				wornOn = "worn on hands"
			case BodyPartWaist:
				wornOn = "worn on waist"
			case BodyPartHoldLeft, BodyPartHoldRight:
				wornOn = "held"
			}
			wornOn = "<" + wornOn + ">"

			view.WriteString(fmt.Sprintf("     %-22s %s\n", wornOn, eqObject.Name))
		}
	} else {
		view.WriteString("     Nothing\n")
	}

	view.WriteString("\n")

	inventory := m.player.GetInventory()
	count := len(inventory)
	view.WriteString("You are carrying (" + strconv.Itoa(count) + "):\n")

	if count > 0 {
		for _, object := range inventory {
			view.WriteString("     " + object.Name + "\n")
		}
	} else {
		view.WriteString("     Nothing\n")
	}

	return view.String()
}
