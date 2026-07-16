package screen

import (
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/exp/charmtone"
)

// TODO: requires more thought on where to keep style instructions
var (
	style = lipgloss.NewStyle().
		Width(70).
		Padding(1, 2).
		Foreground(lipgloss.White).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForegroundBlend(
			charmtone.Cherry,
			charmtone.Charple,
			charmtone.Guac,
			charmtone.Charple,
			charmtone.Sriracha,
		)
)
