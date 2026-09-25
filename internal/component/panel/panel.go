package panel

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/randomvlad/trader-vlads/internal/appstyle"
	"github.com/randomvlad/trader-vlads/internal/component/actionfooter"
	"github.com/randomvlad/trader-vlads/internal/util/stringutil"
)

type Panel struct {
	ActionContextLeft  string
	ActionContextRight string
	title              string
	body               stringutil.Builder
	bodyLayers         []*lipgloss.Layer
	footer             *actionfooter.ActionFooter
	styleConfig        *panelStyleConfig
}

type panelStyleConfig struct {
	width                      int
	height                     int
	bodyBorderFooterCompatible bool
	styleBody                  lipgloss.Style
}

func (p *Panel) Write(bodyContent string) *Panel {
	p.body.Write(bodyContent)
	return p
}

func (p *Panel) WriteLn(bodyContent string) *Panel {
	return p.Write(bodyContent).AddLn()
}

func (p *Panel) AddLn() *Panel {
	p.body.Ln()
	return p
}

func (p *Panel) AddLayer(bodyLayer *lipgloss.Layer) *Panel {
	p.bodyLayers = append(p.bodyLayers, bodyLayer)
	return p
}

func (p *Panel) Render() string {
	var view stringutil.Builder

	if p.title != "" {
		// if title is present, then custom render the top border with embedded title
		view.WriteLn(p.renderTitleInsideTopBorder())
	}

	styleBody := p.getComputedBodyStyle()
	view.WriteStyle(p.body.String(), styleBody)

	if p.footer != nil {
		view.Ln().
			WriteLn(p.footer.Render(p.ActionContextLeft, p.ActionContextRight))
	}

	compositor := lipgloss.NewCompositor()
	compositor.AddLayers(lipgloss.NewLayer(view.String())) // main layer

	// additional layers that are typically small popups
	for _, layer := range p.bodyLayers {
		compositor.AddLayers(layer)
	}

	return compositor.Render()
}

func (p *Panel) RenderStyle(style lipgloss.Style) string {
	return style.Render(p.Render())
}

func (p *Panel) RenderTeaView() tea.View {
	return tea.NewView(p.Render())
}

func (p *Panel) renderTitleInsideTopBorder() string {

	styleBorderTop := appstyle.NewAppStyle().Foreground(appstyle.AppBorderColor)
	styleTitle := appstyle.NewAppStyle()
	stylePointedStart := appstyle.NewAppStyle().Foreground(lipgloss.Color("#F54927"))

	// 2 characters for border corners, 4 for left and right padding each
	borderToFill := p.styleConfig.width - 2 - 4 - 4 - len(p.title)
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

func (p *Panel) getComputedBodyStyle() lipgloss.Style {

	computed := p.styleConfig.styleBody.Width(p.styleConfig.width)
	heightBody := p.styleConfig.height

	if p.styleConfig.bodyBorderFooterCompatible {
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
