package app

import tea "github.com/charmbracelet/bubbletea"

// Run bootstraps the Bubble Tea program.
func Run(cfg Config) error {
	p := tea.NewProgram(initialModel(cfg), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
