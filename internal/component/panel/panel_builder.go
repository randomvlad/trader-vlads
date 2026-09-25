package panel

import (
	"charm.land/lipgloss/v2"
	"github.com/randomvlad/trader-vlads/internal/appmod/keybind"
	"github.com/randomvlad/trader-vlads/internal/appstyle"
	"github.com/randomvlad/trader-vlads/internal/component/actionfooter"
)

func NewPanelBuilder() *PanelBuilder {
	styleBody := appstyle.NewAppStyle().
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(appstyle.AppBorderColor)

	builder := &PanelBuilder{}
	return builder.StyleBody(styleBody)
}

type PanelBuilder struct {
	title                                 string
	keyBinder                             *keybind.KeyBinder
	actionContextLeft, actionContextRight string
	styleBody                             lipgloss.Style
	width, height                         int
	bodyBorderFooterCompatible            bool
}

func (b *PanelBuilder) Title(title string) *PanelBuilder {
	b.title = title
	return b
}

func (b *PanelBuilder) Size(width, height int) *PanelBuilder {
	b.width = width
	b.height = height
	return b
}

func (b *PanelBuilder) StyleBody(style lipgloss.Style) *PanelBuilder {
	b.styleBody = style
	return b
}

func (b *PanelBuilder) Footer(keyBinder *keybind.KeyBinder, actionContextLeft, actionContextRight string) *PanelBuilder {
	b.keyBinder = keyBinder
	b.bodyBorderFooterCompatible = true
	b.actionContextLeft = actionContextLeft
	b.actionContextRight = actionContextRight
	return b
}

func (b *PanelBuilder) BodyBorderFooterCompatible() *PanelBuilder {
	b.bodyBorderFooterCompatible = true
	return b
}

func (p *PanelBuilder) Build() *Panel {
	footer := actionfooter.NewPanelFooter(p.keyBinder, p.width)

	return &Panel{
		footer:             &footer,
		ActionContextLeft:  p.actionContextLeft,
		ActionContextRight: p.actionContextRight,
		styleConfig: &panelStyleConfig{
			width:                      p.width,
			height:                     p.height,
			styleBody:                  p.styleBody,
			bodyBorderFooterCompatible: true,
		},
	}
}
