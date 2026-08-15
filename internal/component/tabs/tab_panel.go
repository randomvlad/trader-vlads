package tabs

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/randomvlad/trader-vlads/internal/appstyle"
	"github.com/randomvlad/trader-vlads/internal/component/actionfooter"
)

type TabPanel struct {
	body       strings.Builder
	bodyLayers []*lipgloss.Layer
	footer     *actionfooter.Model // might be possible to simplify
}

func NewTabPanel(actions ...string) *TabPanel {
	return &TabPanel{
		footer: actionfooter.NewModel(actionfooter.FooterTab, actions...),
	}
}

func (p *TabPanel) Write(bodyContent string) *TabPanel {
	p.body.WriteString(bodyContent)
	return p
}

func (p *TabPanel) WriteLn(bodyContent string) *TabPanel {
	return p.Write(bodyContent).AddLn()
}

func (p *TabPanel) AddLn() *TabPanel {
	p.body.WriteString("\n")
	return p
}

func (p *TabPanel) AddLayer(bodyLayer *lipgloss.Layer) {
	p.bodyLayers = append(p.bodyLayers, bodyLayer)
}

func (p *TabPanel) Render() string {
	compositor := lipgloss.NewCompositor()

	compositor.AddLayers(lipgloss.NewLayer(p.renderBodyWithFooter())) // main layer

	// additional layers that are typically small popups
	for _, layer := range p.bodyLayers {
		compositor.AddLayers(layer)
	}

	return compositor.Render()
}

func (p *TabPanel) RenderTeaView() tea.View {
	return tea.NewView(p.Render())
}

func (p *TabPanel) renderBodyWithFooter() string {
	var view strings.Builder
	view.WriteString(appstyle.StyleTabBodyView.Render(p.body.String()))
	view.WriteString("\n")
	view.WriteString(p.footer.Render())
	return view.String()
}
