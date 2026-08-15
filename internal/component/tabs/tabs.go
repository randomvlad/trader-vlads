package tabs

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/randomvlad/trader-vlads/internal/appstyle"
)

type Model struct {
	Tabs      []string
	ActiveTab int
	width     int
	styles    *tabStyles
}

func NewModel(tabNames []string, width int) *Model {
	return &Model{
		Tabs:      tabNames,
		ActiveTab: 0,
		width:     width,
		styles:    newStyles(),
	}
}

// TODO: TabModel struct with fields: id, display, shortcut key?

type tabStyles struct {
	tabsContainer lipgloss.Style
	inactiveTab   lipgloss.Style
	activeTab     lipgloss.Style
}

func newStyles() *tabStyles {

	inactiveTabBorder := tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder := tabBorderWithBottom("┘", " ", "└")
	inactiveColor := lipgloss.Color("#696969")

	s := new(tabStyles)
	s.tabsContainer = appstyle.NewAppStyle()
	s.inactiveTab = lipgloss.NewStyle().
		Border(inactiveTabBorder, true).
		Foreground(inactiveColor).
		BorderForeground(appstyle.AppBorderColor).
		Padding(0, 1)
	s.activeTab = lipgloss.NewStyle().
		Border(activeTabBorder, true).
		Padding(0, 1).
		Foreground(appstyle.AppTextColor).
		BorderForeground(appstyle.AppBorderColor).
		Bold(true)

	return s
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch keypress := msg.String(); keypress {
		case "right", "tab":
			m.ActiveTab = min(m.ActiveTab+1, len(m.Tabs)-1)
			return m, nil
		case "left", "shift+tab":
			m.ActiveTab = max(m.ActiveTab-1, 0)
			return m, nil
		}
	}

	return m, nil
}

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

func (m *Model) View() tea.View {
	var renderedTabs []string
	for tabIndex, tabName := range m.Tabs {
		borderStyle := getTabStyle(tabIndex, m)
		renderedTabs = append(renderedTabs, borderStyle.Render(tabName))
	}

	viewContent := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)

	widthSoFar := lipgloss.Width(viewContent)
	if widthSoFar < m.width {
		borderRightWidth := 1 // account for 1-char-width right border that is added by spacer style
		neededSpaceWidth := m.width - widthSoFar - borderRightWidth

		// Fill remaining space with a bordered spacer so the bottom line runs all the way to max tabs width
		spaceFiller := getSpacerStyle().Render(strings.Repeat(" ", neededSpaceWidth))
		viewContent = lipgloss.JoinHorizontal(lipgloss.Bottom, viewContent, spaceFiller)
	}

	return tea.NewView(m.styles.tabsContainer.Render(viewContent))
}

func getTabStyle(tabIndex int, m *Model) lipgloss.Style {
	isFirst := tabIndex == 0
	isLast := tabIndex == len(m.Tabs)-1
	isActive := tabIndex == m.ActiveTab

	var style lipgloss.Style
	if isActive {
		style = m.styles.activeTab
	} else {
		style = m.styles.inactiveTab
	}
	border, _, _, _, _ := style.GetBorder()

	if isFirst {
		var left string
		if isActive {
			left = "│"
		} else {
			left = "├"
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

func getSpacerStyle() lipgloss.Style {
	borderSpacer := lipgloss.RoundedBorder()
	borderSpacer.Right = ""
	borderSpacer.BottomLeft = ""
	borderSpacer.Bottom = "─"
	borderSpacer.BottomRight = "╮" // rounded down to connect with "│" border of tab's content view

	return lipgloss.NewStyle().
		Border(borderSpacer, false, true, true, false).
		BorderForeground(appstyle.AppBorderColor)
}
