package app

import "github.com/charmbracelet/lipgloss"

var (
	colorMuted    = lipgloss.AdaptiveColor{Light: "#b09a90", Dark: "#b09a90"} // warm muted
	colorTrack    = lipgloss.AdaptiveColor{Light: "#4a4a4a", Dark: "#3c3c3c"} // gray track background
	colorAccent   = lipgloss.AdaptiveColor{Light: "#d77757", Dark: "#d77757"} // requested base hue
	colorAccentHi = lipgloss.AdaptiveColor{Light: "#f0a889", Dark: "#f0a889"} // soft highlight
	colorError    = lipgloss.AdaptiveColor{Light: "#d15555", Dark: "#ff7b7b"}
	colorWhite    = lipgloss.AdaptiveColor{Light: "#ffffff", Dark: "#ffffff"}
)
