package toast

import (
	"fmt"

	"charm.land/lipgloss/v2"
	"github.com/randomvlad/trader-vlads/internal/util"
)

type Model struct {
	Message string
	Show    bool
}

// TODO: is there value in creating a New() constructor method?

var (
	blendColors = util.ToColors("#8BF578", "#6CCB5B", "#55A147", "#407C35", "#55A147", "#6CCB5B", "#8BF578")
	styleToast  = lipgloss.NewStyle().
			Width(50).
			Padding(2, 3).
			Align(lipgloss.Center).
			Foreground(lipgloss.White).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForegroundBlend(blendColors...)
)

func (m *Model) SetMessage(text string, a ...any) {
	m.Message = fmt.Sprintf(text, a...)
	m.Show = true
}

func (m *Model) Clear() {
	m.Message = ""
	m.Show = false
}

func (m *Model) Render() string {
	return styleToast.Render(m.Message)
}
