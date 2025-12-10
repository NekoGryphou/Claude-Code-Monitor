package app

import (
	"fmt"
	"time"

	"claude-monitor/internal/api"
)

// Config holds runtime settings for the UI.
type Config struct {
	Token          string
	Credentials    string
	RefreshEvery   time.Duration
	HTTPClient     api.HTTPClient
	RequestTimeout time.Duration
}

// Validate ensures required fields are present.
func (c Config) Validate() error {
	if c.Token == "" {
		return fmt.Errorf("token required")
	}
	if c.HTTPClient == nil {
		return fmt.Errorf("http client required")
	}
	if c.RefreshEvery <= 0 {
		return fmt.Errorf("refresh interval must be positive")
	}
	if c.RequestTimeout <= 0 {
		return fmt.Errorf("request timeout must be positive")
	}
	return nil
}
