package app

import "github.com/charmbracelet/lipgloss"

var (
	paletteMuted    = colorMuted
	paletteAccent   = colorAccent
	paletteAccentHi = colorAccentHi
	paletteError    = colorError

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
			Foreground(colorWhite).
			Background(paletteAccent).
			Padding(1, 3).
			MarginBottom(1).
			Bold(true)

	statusStyle = lipgloss.NewStyle().
			Foreground(paletteMuted).
			Italic(true)

	labelBaseStyle = lipgloss.NewStyle().
			Foreground(colorWhite).
			Bold(true)

	resetBaseStyle = lipgloss.NewStyle().
			Foreground(paletteMuted).
			Italic(true)

	remainBaseStyle = lipgloss.NewStyle().
			Foreground(paletteAccentHi)

	valueBaseStyle = lipgloss.NewStyle().
			Foreground(colorWhite).
			Bold(true)

	spinnerStyle = lipgloss.NewStyle().Foreground(paletteAccent)

	barBlockStyle     = lipgloss.NewStyle().MarginBottom(1)
	barLastBlockStyle = lipgloss.NewStyle()
	barSeparatorStyle = lipgloss.NewStyle().Width(1)
	barFillStyle      = lipgloss.NewStyle().
				Background(paletteAccent).
				Foreground(colorWhite)
	barEmptyStyle = lipgloss.NewStyle().
			Background(colorTrack).
			Foreground(colorWhite)

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
)
