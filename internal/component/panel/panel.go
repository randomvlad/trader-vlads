package panel

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/randomvlad/trader-vlads/internal/appmod/keybind/press"
	"github.com/randomvlad/trader-vlads/internal/appstyle"
	"github.com/randomvlad/trader-vlads/internal/component/actionfooter"
	"github.com/randomvlad/trader-vlads/internal/util/stringutil"
)

type Model struct {
	title      string
	body       stringutil.Builder
	bodyLayers []*lipgloss.Layer
	footer     *actionfooter.Model
	visual     *modelVisual
}

type modelVisual struct {
	width                      int
	height                     int
	bodyBorderFooterCompatible bool
	styleBody                  lipgloss.Style
}

func NewModel() *Model {
	styleBody := appstyle.NewAppStyle().
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(appstyle.AppBorderColor)

	return &Model{
		visual: &modelVisual{
			width:     80,
			height:    20,
			styleBody: styleBody,
		},
	}
}

func (p *Model) WithTitle(title string) *Model {
	p.title = title
	return p
}

func (p *Model) WithWidth(value int) *Model {
	p.visual.width = value
	return p
}

func (p *Model) WithHeight(value int) *Model {
	p.visual.height = value
	return p
}

func (p *Model) WithStyle(styleBody lipgloss.Style) *Model {
	p.visual.styleBody = styleBody
	return p
}

func (p *Model) WithFooter(actions ...press.KeyAction) *Model {
	p.footer = actionfooter.NewModel(actionfooter.FooterPanel, actions...)
	return p.WithBodyBorderFooterCompatible()
}

func (p *Model) WithBodyBorderFooterCompatible() *Model {
	p.visual.bodyBorderFooterCompatible = true
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

func (p *Model) AddLayer(bodyLayer *lipgloss.Layer) *Model {
	p.bodyLayers = append(p.bodyLayers, bodyLayer)
	return p
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
	var render stringutil.Builder

	styleBody := p.getComputedBodyStyle()

	if p.title != "" {
		// if title is present, then custom render the top border with embedded title
		render.WriteLn(p.renderTopBorderWithTitle())
	}

	render.WriteStyle(p.body.String(), styleBody)

	if p.footer != nil {
		p.footer.WithWidth(p.visual.width)
		render.Ln().WriteLn(p.footer.Render("")) // TODO: action context
	}

	return render.String()
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

	return stringutil.NewBuilder().
		WithStyle(styleBorderTop).
		Write(borderDef.TopLeft).
		WriteRepeat(borderDef.Top, sideLengthLeft).
		Write("[").
		WriteStyle(" ✧ ", stylePointedStart).
		WriteStyle(stringutil.Truncate(p.title, 60), styleTitle).
		WriteStyle(" ✧ ", stylePointedStart).
		Write("]").
		WriteRepeat(borderDef.Top, sideLengthRight).
		Write(borderDef.TopRight).
		String()
}

func (p *Model) getComputedBodyStyle() lipgloss.Style {

	computed := p.visual.styleBody.Width(p.visual.width)
	heightBody := p.visual.height

	if p.visual.bodyBorderFooterCompatible {
		if p.footer != nil {
			heightBody -= 2 // reduce body height to make room for the footer
		}

		// Configure body border to seamlessly connect with footer border that follows
		border, top, right, bottom, left := computed.GetBorder()
		border.BottomLeft = "├"
		border.BottomRight = "┤"
		computed = computed.Border(border, top, right, bottom, left)
	}

	if p.title != "" {
		// top border with title takes 1 line and is rendered before the body starts
		heightBody -= 1

		// disable default lipgloss border top, title is part of top border and requires custom rendering
		computed = computed.BorderTop(false)
	}

	return computed.Height(heightBody)
}
