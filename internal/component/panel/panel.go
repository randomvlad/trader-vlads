package panel

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/randomvlad/trader-vlads/internal/appstyle"
	"github.com/randomvlad/trader-vlads/internal/component/actionfooter"
	"github.com/randomvlad/trader-vlads/internal/util"
)

type Model struct {
	title      string
	body       util.StringBuilder
	bodyLayers []*lipgloss.Layer
	footer     *actionfooter.Model
	visual     *modelVisual
}

type modelVisual struct {
	width      int
	height     int
	stylePanel lipgloss.Style
}

func NewModel() *Model {
	m := &Model{
		visual: &modelVisual{
			stylePanel: appstyle.NewAppStyle().
				Padding(1, 2).
				Border(lipgloss.RoundedBorder(), false, true, true, true).
				BorderForeground(appstyle.AppBorderColor),
		},
	}
	return m.WithWidth(80).WithHeight(20)
}

func (p *Model) WithTitle(title string) *Model {
	p.title = title
	return p
}

func (p *Model) WithWidth(value int) *Model {
	p.visual.width = value
	p.visual.stylePanel = p.visual.stylePanel.Width(value)
	return p
}

func (p *Model) WithHeight(value int) *Model {
	p.visual.height = value
	p.visual.stylePanel = p.visual.stylePanel.Height(value)
	return p
}

func (p *Model) WithFooter(actions ...string) *Model {
	p.footer = actionfooter.NewModel(actionfooter.FooterNoStyle, actions...)
	return p
}

func (p *Model) Write(bodyContent string) *Model {
	p.body.Write(bodyContent)
	return p
}

func (p *Model) WriteLn(bodyContent string) *Model {
	return p.Write(bodyContent).AddLn()
}

func (p *Model) AddLn() *Model {
	p.body.Ln()
	return p
}

func (p *Model) AddLayer(bodyLayer *lipgloss.Layer) {
	p.bodyLayers = append(p.bodyLayers, bodyLayer)
}

func (p *Model) Render() string {
	compositor := lipgloss.NewCompositor()

	compositor.AddLayers(lipgloss.NewLayer(p.renderBodyWithFooter())) // main layer

	// additional layers that are typically small popups
	for _, layer := range p.bodyLayers {
		compositor.AddLayers(layer)
	}

	return compositor.Render()
}

func (p *Model) RenderTeaView() tea.View {
	return tea.NewView(p.Render())
}

func (p *Model) renderBodyWithFooter() string {
	var render util.StringBuilder

	panelStyle := p.visual.stylePanel

	if p.title != "" {
		render.Write(p.renderTopBorderWithTitle())
	} else {
		// given no title, let lipgloss draw a round border on all sides
		panelStyle = panelStyle.Border(lipgloss.RoundedBorder())
	}

	// append footer at the end of body
	p.body.Ln().WriteLn(p.footer.Render())

	return render.WriteStylized(p.body.String(), panelStyle).String()
}

func (p *Model) renderTopBorderWithTitle() string {

	styleBorderTop := appstyle.NewAppStyle().Foreground(appstyle.AppBorderColor)
	styleTitle := appstyle.NewAppStyle()
	stylePointedStart := appstyle.NewAppStyle().Foreground(lipgloss.Color("#F54927"))

	// 2 characters for border corners, 4 for left and right padding each
	borderToFill := p.visual.width - 2 - 4 - 4 - len(p.title)
	sideLengthLeft := borderToFill / 2
	sideLengthRight := borderToFill - sideLengthLeft

	borderDef := lipgloss.RoundedBorder()

	return new(util.StringBuilder).
		WithStyle(styleBorderTop).
		Write(borderDef.TopLeft).
		WriteRepeat(borderDef.Top, sideLengthLeft).
		Write("[").
		WriteStylized(" ✧ ", stylePointedStart).
		WriteStylized(util.Truncate(p.title, 60), styleTitle).
		WriteStylized(" ✧ ", stylePointedStart).
		Write("]").
		WriteRepeat(borderDef.Top, sideLengthRight).
		WriteLn(borderDef.TopRight).
		String()
}
