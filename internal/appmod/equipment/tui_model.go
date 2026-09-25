package equipment

import (
	"strconv"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/randomvlad/trader-vlads/internal/appmod/keybind"
	"github.com/randomvlad/trader-vlads/internal/appmod/keybind/press"
	eff "github.com/randomvlad/trader-vlads/internal/appmod/stats/statuseffect"
	"github.com/randomvlad/trader-vlads/internal/appstyle"
	"github.com/randomvlad/trader-vlads/internal/component/tabs"
	"github.com/randomvlad/trader-vlads/internal/util/stringutil"
)

type Model struct {
	selectionIndex int
	player         PlayerService
	keyBinder      *keybind.KeyBinder
	toast          ToastMessenger
}

func NewTuiModel(player PlayerService, keyBinder *keybind.KeyBinder, toast ToastMessenger) *Model {
	return &Model{
		player:    player,
		keyBinder: keyBinder,
		toast:     toast,
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

	eqNext := press.NewActionBuilder().
		NameKeyPress("EqNext", "down").
		Action(func() { m.moveCursorPosition(true) }).
		Permanent().
		FooterHidden().
		Build()

	eqPrev := press.NewActionBuilder().
		NameKeyPress("EqPrev", "up").
		Action(func() { m.moveCursorPosition(false) }).
		FooterHidden().
		Permanent().
		Build()

	actionRemove := press.NewActionBuilder().
		Name("Remove").
		Action(func() { m.removeEq() }).
		FooterVisible(func() bool {
			return m.isSelectedBodyPart() && m.player.HasEquipped(BodyPart(m.selectionIndex))
		}).
		Permanent().
		Build()

	actionWear := press.NewActionBuilder().
		Name("Wear").
		Action(func() { m.wearEq() }).
		FooterVisible(func() bool {
			return !m.isSelectedBodyPart() && m.getSelectedObject().IsWearable()
		}).
		Permanent().
		Build()

	actionUse := press.NewActionBuilder().
		Name("Use").
		Action(func() { m.useItem() }).
		FooterVisible(func() bool {
			return !m.isSelectedBodyPart() && m.getSelectedObject().IsUsable()
		}).
		Permanent().
		Build()

	m.keyBinder.AddActions("tui-eq", eqNext, eqPrev, actionRemove, actionWear, actionUse)

	return nil
}

func (m *Model) View() tea.View {

	panel := tabs.NewTabPanel()

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
	return m, nil
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

func (m *Model) removeEq() {
	invIndex := m.player.Remove(BodyPart(m.selectionIndex))
	if invIndex >= 0 {
		// move selection to the removed item in inventory
		m.selectionIndex = invIndex + BodyPartsMax
	}
}

func (m *Model) wearEq() {
	invIndex := m.selectionIndex - BodyPartsMax
	eqObject := m.player.GetInventoryObject(invIndex)
	ok := m.player.WearInventory(invIndex)
	if ok {
		if len(eqObject.Effects) > 0 {
			// TODO: this is a temp hack. need event pub/sub system. Player publishes event messages.
			m.toast.Message(eqObject.Effects[0].GetMessageStart())
		}

		if m.selectionIndex >= BodyPartsMax+len(m.player.GetInventory()) {
			m.selectionIndex -= 1 // inv has shrunk so move to previous item
		}
	}
}

func (m *Model) useItem() {
	invIndex := m.selectionIndex - BodyPartsMax
	eqObject := m.player.GetInventoryObject(invIndex)
	if eqObject != nil && eqObject.IsUsable() {
		used := m.player.Use(invIndex)
		if used {
			if m.selectionIndex >= BodyPartsMax+len(m.player.GetInventory()) {
				m.selectionIndex -= 1 // inv has shrunk so move to previous item
			}

			if len(eqObject.Effects) > 0 {
				m.toast.Message(eqObject.Effects[0].GetMessageStart())
			}
		}
	}
}

func (m *Model) renderEq() string {

	view := stringutil.NewBuilder().WriteLn("You are using:")

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
			view.Write(appstyle.SelectionPointer)
		} else {
			view.Write(" ")
		}

		view.Writef("   %-22s %s\n", wornOn, eqObjectDisplay)
	}

	return view.String()
}

func (m *Model) renderInv() string {

	inventory := m.player.GetInventory()
	count := len(inventory)

	view := stringutil.NewBuilder().
		WriteLn("You are carrying (" + strconv.Itoa(count) + "):")

	positionEqOffset := len(m.player.GetEquipped())

	if count > 0 {
		for index, object := range inventory {
			if (index + positionEqOffset) == m.selectionIndex {
				view.Write(appstyle.SelectionPointer)
			} else {
				view.Write(" ")
			}
			view.WriteLn("   " + object.Name)
		}
	} else {
		view.WriteLn("    Nothing")
	}

	return view.String()
}

func (m *Model) getSelectedObject() *EqObject {
	if m.isSelectedBodyPart() {
		return m.player.GetEquippedObject(BodyPart(m.selectionIndex))
	} else {
		invIndex := m.selectionIndex - BodyPartsMax
		return m.player.GetInventoryObject(invIndex)
	}
}

func (m *Model) isSelectedBodyPart() bool {
	return m.selectionIndex < BodyPartsMax
}
