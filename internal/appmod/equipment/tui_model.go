package equipment

import (
	"fmt"
	"strconv"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	eff "github.com/randomvlad/trader-vlads/internal/appmod/stats/statuseffect"
	"github.com/randomvlad/trader-vlads/internal/appstyle"
	"github.com/randomvlad/trader-vlads/internal/component/tabs"
)

type Model struct {
	selectionIndex int
	player         PlayerService
	toast          ToastMessenger
}

func NewTuiModel(player PlayerService, toast ToastMessenger) *Model {
	return &Model{
		player: player,
		toast:  toast,
	}
}

type PlayerService interface {
	HasEquipped(BodyPart) bool
	GetEquipped() map[BodyPart]*EqObject
	GetEquippedObject(bodyPart BodyPart) *EqObject
	Remove(bodyPart BodyPart) int
	GetInventory() []*EqObject
	GetInventoryObject(invIndex int) *EqObject
	WearInventory(invIndex int) bool
	AddEffects(effects ...eff.StatusEffect)
	Use(invIndex int) bool
}

type ToastMessenger interface {
	Message(text string, a ...any)
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) View() tea.View {

	panel := tabs.NewTabPanel(m.getActions()...)

	panel.
		WriteLn(m.renderEq()).
		AddLn().
		WriteLn(m.renderInv())

	eqObject := m.getSelectedObject()
	if eqObject != nil {
		stats := eqObject.ViewStats()
		panel.AddLayer(lipgloss.NewLayer(appstyle.StyleEqStats.Render(stats)).X(68).Y(0).Z(1))
	}

	return panel.RenderTeaView()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "down":
			m.moveCursorPosition(true)
		case "up":
			m.moveCursorPosition(false)
		case "r", "R":
			invIndex := m.player.Remove(BodyPart(m.selectionIndex))
			if invIndex >= 0 {
				// move selection to the removed item in inventory
				m.selectionIndex = invIndex + BodyPartsMax
			}
		case "w", "W":
			invIndex := m.selectionIndex - BodyPartsMax
			eqObject := m.player.GetInventoryObject(invIndex)
			ok := m.player.WearInventory(invIndex)
			if ok {
				if len(eqObject.Effects) > 0 {
					// TODO: this is a temp hack. need event pub/sub system. Player publishes event messages.
					m.toast.Message(eqObject.Effects[0].GetAppliedMessage())
				}

				if m.selectionIndex >= BodyPartsMax+len(m.player.GetInventory()) {
					m.selectionIndex -= 1 // inv has shrunk so move to previous item
				}
			}
		case "u", "U":
			invIndex := m.selectionIndex - BodyPartsMax
			eqObject := m.player.GetInventoryObject(invIndex)
			if eqObject != nil && eqObject.IsUsable() {
				used := m.player.Use(invIndex)
				if used {
					if m.selectionIndex >= BodyPartsMax+len(m.player.GetInventory()) {
						m.selectionIndex -= 1 // inv has shrunk so move to previous item
					}

					if len(eqObject.Effects) > 0 {
						m.toast.Message(eqObject.Effects[0].GetAppliedMessage())
					}
				}
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) moveCursorPosition(nextOrPrev bool) {

	totalLength := BodyPartsMax + len(m.player.GetInventory())

	if nextOrPrev { // move selection to next
		m.selectionIndex++
		if m.selectionIndex >= totalLength {
			m.selectionIndex = 0 // wrap to start
		}
	} else { // move selection to previous
		m.selectionIndex--
		if m.selectionIndex < 0 {
			m.selectionIndex = totalLength - 1 // wrap to end
		}
	}
}

func (m *Model) renderEq() string {

	var view strings.Builder

	view.WriteString("You are using:\n")

	eq := m.player.GetEquipped()

	for bodyPartIndex := range BodyPartsMax {
		bodyPart := BodyPart(bodyPartIndex)
		eqObject := eq[bodyPart]

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
		case BodyPartHoldLeft:
			wornOn = "held in left"
		case BodyPartHoldRight:
			wornOn = "held in right"
		}
		wornOn = "<" + wornOn + ">"

		var eqObjectDisplay string
		if eqObject != nil {
			eqObjectDisplay = eqObject.Name
		} else {
			eqObjectDisplay = " "
		}

		if m.selectionIndex == bodyPartIndex {
			view.WriteString(appstyle.SelectionPointer)
		} else {
			view.WriteString(" ")
		}

		view.WriteString(fmt.Sprintf("   %-22s %s\n", wornOn, eqObjectDisplay))
	}

	return view.String()
}

func (m *Model) renderInv() string {
	var view strings.Builder

	inventory := m.player.GetInventory()
	count := len(inventory)
	view.WriteString("You are carrying (" + strconv.Itoa(count) + "):\n")

	positionEqOffset := len(m.player.GetEquipped())

	if count > 0 {
		for index, object := range inventory {
			if (index + positionEqOffset) == m.selectionIndex {
				view.WriteString(appstyle.SelectionPointer)
			} else {
				view.WriteString(" ")
			}
			view.WriteString("   " + object.Name + "\n")
		}
	} else {
		view.WriteString("    Nothing\n")
	}

	return view.String()
}

func (m *Model) getActions() []string {
	var actions []string

	if m.selectionIndex < BodyPartsMax {
		if m.player.HasEquipped(BodyPart(m.selectionIndex)) {
			actions = append(actions, "Remove")
		}
	} else {
		invObject := m.getSelectedObject()
		if invObject.IsWearable() {
			actions = append(actions, "Wear")
		} else if invObject.IsUsable() {
			actions = append(actions, "Use")
		}
	}

	return actions
}

func (m *Model) getSelectedObject() *EqObject {
	if m.selectionIndex < BodyPartsMax {
		return m.player.GetEquippedObject(BodyPart(m.selectionIndex))
	} else {
		invIndex := m.selectionIndex - BodyPartsMax
		return m.player.GetInventoryObject(invIndex)
	}
}
