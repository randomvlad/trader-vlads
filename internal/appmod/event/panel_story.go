package event

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/randomvlad/trader-vlads/internal/appmod/keybind"
	"github.com/randomvlad/trader-vlads/internal/appstyle"
	apppanel "github.com/randomvlad/trader-vlads/internal/component/panel"
)

type StoryPanel struct {
	keyBinder *keybind.KeyBinder
	style     lipgloss.Style
}

func NewStoryPanel(keyBinder *keybind.KeyBinder) StoryPanel {
	stylePanelContainer := appstyle.NewAppStyle().
		Width(appstyle.AppWidth).
		Height(appstyle.AppHeight).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForegroundBlend(appstyle.GreenBlendColors...).
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center)

	return StoryPanel{keyBinder: keyBinder, style: stylePanelContainer}
}

func (p *StoryPanel) Render(story Story) string {
	panel := apppanel.NewPanelBuilder().
		Size(80, 25).
		Footer(p.keyBinder, story.GetStateId(), "").
		Build()

	panel.Write(story.View())

	return panel.RenderStyle(p.style)
}

func (p *StoryPanel) Update(storyStateId string, msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	if p.keyBinder.IsBound(storyStateId, msg) {
		cmd := p.keyBinder.Execute(storyStateId, msg)
		cmds = append(cmds, cmd)
	}
	return tea.Batch(cmds...)
}
