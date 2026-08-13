package appstyle

import (
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/compat"
	"github.com/charmbracelet/x/exp/charmtone"
	"github.com/randomvlad/trader-vlads/internal/util"
)

// TODO: requires more thought and structure how to centralize and manage app wide styles
var (
	AppWidth       = 120
	AppBorderColor = compat.AdaptiveColor{
		Light: lipgloss.Color("#04B575"),
		Dark:  lipgloss.Color("#04B575"),
	}

	// TODO: research: compat.AdaptiveColor{Light: lipgloss.Color("#f1f1f1"), Dark: lipgloss.Color("#cccccc")}
	// See: https://github.com/charmbracelet/lipgloss/releases
	AppTextColor = compat.AdaptiveColor{
		Light: lipgloss.Black,
		Dark:  lipgloss.BrightWhite,
	}

	StyleAppContainer = NewAppStyle()

	PopupWidth = 70
	StylePopup = NewAppStyle().
			Width(PopupWidth).
			Padding(1, 2).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForegroundBlend(
			charmtone.Cherry,
			charmtone.Charple,
			charmtone.Guac,
			charmtone.Charple,
			charmtone.Sriracha)

	TabViewHeight = 30
	StyleTabView  = NewAppStyle().
			Width(AppWidth).
			Height(TabViewHeight).
			Padding(0, 2).
			Border(lipgloss.RoundedBorder(), false, true, true, true).
			BorderForeground(AppBorderColor)

	StyleBadge = NewAppStyle().
			Padding(0, 2).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(AppBorderColor)

	StyleActionsBar = NewAppStyle().
			Width(AppWidth).
			Padding(0, 2).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(AppBorderColor)

	StyleTextFirstLetter = NewAppStyle().Bold(true).Underline(true)

	blendColors = util.ToColors("#8BF578", "#6CCB5B", "#55A147", "#407C35", "#55A147", "#6CCB5B", "#8BF578")
	StyleToast  = NewAppStyle().
			Width(50).
			Padding(2, 3).
			Align(lipgloss.Center).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForegroundBlend(blendColors...)
)

func NewAppStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(AppTextColor)
}
