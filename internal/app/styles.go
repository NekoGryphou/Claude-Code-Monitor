package app

import (
	"os"
	"strings"

	"claude-monitor/internal/consts"

	"github.com/charmbracelet/lipgloss"
)

var (
	paletteMuted    lipgloss.TerminalColor = consts.ColorMuted
	paletteAccent   lipgloss.TerminalColor = consts.ColorAccent
	paletteAccentHi lipgloss.TerminalColor = consts.ColorAccentHi
	paletteError    lipgloss.TerminalColor = consts.ColorError

	pageStyle             lipgloss.Style
	chartBoxStyle         lipgloss.Style
	markStyle             lipgloss.Style
	headerStyle           lipgloss.Style
	statusStyle           lipgloss.Style
	labelBaseStyle        lipgloss.Style
	resetBaseStyle        lipgloss.Style
	remainBaseStyle       lipgloss.Style
	valueBaseStyle        lipgloss.Style
	spinnerStyle          lipgloss.Style
	barBlockStyle         lipgloss.Style
	barLastBlockStyle     lipgloss.Style
	barSeparatorStyle     lipgloss.Style
	barFillStyle          lipgloss.Style
	barEmptyStyle         lipgloss.Style
	metaBaseStyle         lipgloss.Style
	footerStyle           lipgloss.Style
	helpStyle             lipgloss.Style
	helpKeyStyle          lipgloss.Style
	helpDescStyle         lipgloss.Style
	separatorStyle        lipgloss.Style
	errorBoxStyleBase     lipgloss.Style
	errorContentStyleBase lipgloss.Style
	skeletonLabelStyle    lipgloss.Style
	skeletonBarFillStyle  lipgloss.Style
	skeletonBarEmptyStyle lipgloss.Style
	skeletonMetaStyle     lipgloss.Style
)

func init() {
	applyThemeOverrides()
	initStyles()
}

func applyThemeOverrides() {
	if _, ok := os.LookupEnv("NO_COLOR"); ok {
		blank := lipgloss.Color("")
		paletteMuted = blank
		paletteAccent = blank
		paletteAccentHi = blank
		paletteError = blank
		return
	}
	if v, ok := os.LookupEnv("CLAUDE_MONITOR_HIGH_CONTRAST"); ok && isTruthy(v) {
		paletteMuted = lipgloss.AdaptiveColor{Light: "#8a8f98", Dark: "#c3c7cf"}
		paletteAccent = lipgloss.AdaptiveColor{Light: "#ff6b3d", Dark: "#ff8a50"}
		paletteAccentHi = lipgloss.AdaptiveColor{Light: "#ffd7c2", Dark: "#ffe1cf"}
		paletteError = lipgloss.AdaptiveColor{Light: "#ff4d4f", Dark: "#ff7b84"}
	}
}

func isTruthy(v string) bool {
	v = strings.TrimSpace(strings.ToLower(v))
	return v == "1" || v == "true" || v == "yes" || v == "on"
}

func initStyles() {
	pageStyle = lipgloss.NewStyle().
		Padding(1, 2)

	chartBoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(paletteAccent).
		Padding(1, 2)

	markStyle = lipgloss.NewStyle().
		Foreground(paletteAccent).
		MarginBottom(1).
		MarginRight(2)

	headerStyle = lipgloss.NewStyle().
		Foreground(consts.ColorWhite).
		Background(paletteAccent).
		Padding(1, 3).
		MarginBottom(1).
		Bold(true)

	statusStyle = lipgloss.NewStyle().
		Foreground(paletteMuted).
		Italic(true)

	labelBaseStyle = lipgloss.NewStyle().
		Foreground(consts.ColorWhite).
		Bold(true)

	resetBaseStyle = lipgloss.NewStyle().
		Foreground(paletteMuted).
		Italic(true)

	remainBaseStyle = lipgloss.NewStyle().
		Foreground(paletteAccentHi)

	valueBaseStyle = lipgloss.NewStyle().
		Foreground(consts.ColorWhite).
		Bold(true)

	spinnerStyle = lipgloss.NewStyle().Foreground(paletteAccent)

	barBlockStyle = lipgloss.NewStyle().MarginBottom(1)
	barLastBlockStyle = lipgloss.NewStyle()
	barSeparatorStyle = lipgloss.NewStyle().Width(1)
	barFillStyle = lipgloss.NewStyle().
		Background(paletteAccent).
		Foreground(consts.ColorWhite)
	barEmptyStyle = lipgloss.NewStyle().
		Background(consts.ColorTrack).
		Foreground(consts.ColorWhite)

	metaBaseStyle = lipgloss.NewStyle()

	footerStyle = lipgloss.NewStyle().
		MarginTop(1).
		Foreground(paletteMuted)

	helpStyle = lipgloss.NewStyle().
		Foreground(paletteMuted)

	helpKeyStyle = lipgloss.NewStyle().
		Foreground(paletteAccentHi).
		Bold(true)

	helpDescStyle = lipgloss.NewStyle().
		Foreground(paletteMuted)

	separatorStyle = lipgloss.NewStyle()

	errorBoxStyleBase = lipgloss.NewStyle().
		Foreground(paletteError).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(paletteError).
		Padding(0, 1)

	errorContentStyleBase = lipgloss.NewStyle()

	skeletonLabelStyle = lipgloss.NewStyle().
		Foreground(paletteMuted).
		Bold(true)
	skeletonBarFillStyle = lipgloss.NewStyle().Background(paletteMuted)
	skeletonBarEmptyStyle = lipgloss.NewStyle().Background(consts.ColorTrack)
	skeletonMetaStyle = lipgloss.NewStyle().
		Foreground(paletteMuted).
		Background(consts.ColorTrack)
}
