package consts

import "github.com/charmbracelet/lipgloss"

// Palette used across the TUI. Adaptive colors keep light/dark parity.
var (
	// ColorMuted is the subdued text color for labels and metadata.
	ColorMuted = lipgloss.AdaptiveColor{Light: "#b09a90", Dark: "#b09a90"}
	// ColorTrack fills the empty portion of progress bars.
	ColorTrack = lipgloss.AdaptiveColor{Light: "#4a4a4a", Dark: "#3c3c3c"}
	// ColorAccent is the primary accent for headers and filled bars.
	ColorAccent = lipgloss.AdaptiveColor{Light: "#d77757", Dark: "#d77757"}
	// ColorAccentHi highlights hot accents such as keybinds.
	ColorAccentHi = lipgloss.AdaptiveColor{Light: "#f0a889", Dark: "#f0a889"}
	// ColorError tints error borders and text.
	ColorError = lipgloss.AdaptiveColor{Light: "#d15555", Dark: "#ff7b7b"}
	// ColorWhite is the base foreground on colored backgrounds.
	ColorWhite = lipgloss.AdaptiveColor{Light: "#ffffff", Dark: "#ffffff"}
)
