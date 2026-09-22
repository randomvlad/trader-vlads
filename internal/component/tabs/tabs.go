package tabs

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/randomvlad/trader-vlads/internal/appmod/keybind/press"
	"github.com/randomvlad/trader-vlads/internal/appstyle"
	apppanel "github.com/randomvlad/trader-vlads/internal/component/panel"
)

type Model struct {
	Tabs      []string
	ActiveTab int
	visual    *modelVisual
}

type modelVisual struct {
	width       int
	inactiveTab lipgloss.Style
	activeTab   lipgloss.Style
}

func NewModel(tabNames ...string) *Model {
	return &Model{
		Tabs:      tabNames,
		ActiveTab: 0,
		visual:    newModelVisual(appstyle.AppWidth),
	}
}

func newModelVisual(width int) *modelVisual {

	inactiveTabBorder := tabBorderWithBottom("┴", "─", "┴")
	styleInactiveTab := lipgloss.NewStyle().
		Border(inactiveTabBorder).
		Foreground(lipgloss.Color("#696969")).
		BorderForeground(appstyle.AppBorderColor).
		Padding(0, 1)

	activeTabBorder := tabBorderWithBottom("┘", " ", "└")
	styleActiveTab := lipgloss.NewStyle().
		Border(activeTabBorder).
		Foreground(appstyle.AppTextColor).
		BorderForeground(appstyle.AppBorderColor).
		Bold(true).
		Padding(0, 1)

	return &modelVisual{
		width:       width,
		inactiveTab: styleInactiveTab,
		activeTab:   styleActiveTab,
	}
}

func NewTabPanel(actions ...press.KeyAction) *apppanel.Model {

	styleTabBody := appstyle.NewAppStyle().
		Padding(0, 2).
		Border(lipgloss.RoundedBorder()).
		BorderTop(false). // top border is drawn by tabs view
		BorderForeground(appstyle.AppBorderColor)

	model := apppanel.NewModel().
		WithStyle(styleTabBody).
		WithWidth(appstyle.AppWidth).
		WithHeight(appstyle.TabHeight)

	if len(actions) > 0 {
		model.WithFooter(actions...)
	}

	return model
}

// TODO: []TabModel struct with fields: id, display, shortcut key?

func (m *Model) SelectRight() {
	m.ActiveTab = min(m.ActiveTab+1, len(m.Tabs)-1)
}

func (m *Model) SelectLeft() {
	m.ActiveTab = max(m.ActiveTab-1, 0)
}

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

func (m *Model) View() string {
	renderedTabs := []string{getSpacerLeft()}

	for tabIndex, tabName := range m.Tabs {
		borderStyle := getTabStyle(tabIndex, m)
		renderedTabs = append(renderedTabs, borderStyle.Render(tabName))
	}

	viewContent := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)

	widthSoFar := lipgloss.Width(viewContent)
	if widthSoFar < m.visual.width {
		borderRightWidth := 1 // account for 1-char-width right border that is added by spacer style
		neededSpaceWidth := m.visual.width - widthSoFar - borderRightWidth

		// Fill remaining space with a bordered spacer so the bottom line runs all the way to max tabs width
		spaceFiller := getSpacerRight(neededSpaceWidth)

		viewContent = lipgloss.JoinHorizontal(lipgloss.Bottom, viewContent, spaceFiller)
	}

	return viewContent
}

func getSpacerLeft() string {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = "├"
	border.BottomRight = ""

	style := lipgloss.NewStyle().
		BorderForeground(appstyle.AppBorderColor).
		Border(border, false, false, true, true)

	return style.Render("   \n ")
}

func getTabStyle(tabIndex int, m *Model) lipgloss.Style {
	isFirst := tabIndex == 0
	isLast := tabIndex == len(m.Tabs)-1
	isActive := tabIndex == m.ActiveTab

	var style lipgloss.Style
	if isActive {
		style = m.visual.activeTab
	} else {
		style = m.visual.inactiveTab
	}
	border, _, _, _, _ := style.GetBorder()

	if isFirst {
		var left string
		if isActive {
			left = "┘"
		} else {
			left = "┴"
		}
		border.BottomLeft = left
	} else if isLast {
		var right string
		if isActive {
			right = "└"
		} else {
			right = "┴" // spacer with "─" border follows next to fill the gab from last tab to total tabs width
		}
		border.BottomRight = right
	}

	return style.Border(border)
}

func getSpacerRight(width int) string {
	spaceContent := strings.Repeat(" ", width) + "\n " // trailing newline & space to fill vertical space

	borderSpacer := lipgloss.RoundedBorder()
	borderSpacer.Right = "│"
	borderSpacer.BottomLeft = ""
	borderSpacer.Bottom = "─"
	borderSpacer.BottomRight = "┤" // to connect with "│" border of tab's content view

	style := lipgloss.NewStyle().
		Border(borderSpacer, false, true, true, false).
		BorderForeground(appstyle.AppBorderColor)

	// Fill remaining space with a bordered spacer so the bottom line runs all the way to max tabs width
	return style.Render(spaceContent)
}
