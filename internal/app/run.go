package app

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
)

// Run starts the Bubble Tea program configured with cfg.
//
// Returns:
//   - nil on a clean shutdown.
//   - an error if the Bubble Tea program fails to start or run.
func Run(ctx context.Context, cfg Config) error {
	p := tea.NewProgram(initialModel(ctx, cfg), tea.WithAltScreen(), tea.WithContext(ctx))
	_, err := p.Run()
	return err
}
