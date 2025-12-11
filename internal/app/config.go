package app

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"claude-monitor/internal/consts"
)

// Config holds runtime options for the TUI.
type Config struct {
	// Token is the OAuth bearer token used for API requests.
	Token string
	// Credentials is the path to the credentials file (for diagnostics).
	Credentials string
	// RefreshEvery controls the polling interval for usage fetches.
	RefreshEvery time.Duration
	// HTTPClient executes API requests; its Timeout is also used for per-call deadlines.
	HTTPClient *http.Client
	// BetaHeader carries the anthropic-beta header value required by the API.
	BetaHeader string
}

// Validate ensures the configuration is usable before running the UI.
//
// Parameters:
//   - none (receiver-based validation).
//
// Returns:
//   - nil when all required fields are present and values are positive.
//   - an error describing the first missing or invalid field.
func (c Config) Validate() error {
	if c.Token == "" {
		return fmt.Errorf(consts.ErrTokenRequired)
	}
	if c.HTTPClient == nil {
		return fmt.Errorf(consts.ErrHTTPClientRequired)
	}
	if c.HTTPClient.Timeout <= 0 {
		return fmt.Errorf(consts.ErrHTTPClientTimeout)
	}
	const minRefresh = 200 * time.Millisecond
	if strings.TrimSpace(c.BetaHeader) == "" {
		return fmt.Errorf(consts.ErrBetaHeaderRequired)
	}
	if c.RefreshEvery <= 0 {
		return fmt.Errorf(consts.ErrRefreshInterval)
	}
	if c.RefreshEvery < minRefresh {
		return fmt.Errorf("refresh interval too small; must be at least %s", minRefresh)
	}
	return nil
}
