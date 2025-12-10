package util

import (
	"claude-monitor/internal/texts"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type credentialsFile struct {
	ClaudeOauth struct {
		AccessToken string `json:"accessToken"`
	} `json:"claudeAiOauth"`
}

// ResolveToken finds the token from env or credentials file.
func ResolveToken(credPath string) (string, error) {
	if token := strings.TrimSpace(os.Getenv("ANTHROPIC_OAUTH_TOKEN")); token != "" {
		return token, nil
	}
	path := credPath
	if strings.HasPrefix(path, "~") {
		if home, err := os.UserHomeDir(); err == nil {
			path = filepath.Join(home, strings.TrimPrefix(path, "~"))
		}
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read credentials: %w", err)
	}
	var creds credentialsFile
	if err := json.Unmarshal(content, &creds); err != nil {
		return "", fmt.Errorf("parse credentials: %w", err)
	}
	token := strings.TrimSpace(creds.ClaudeOauth.AccessToken)
	if token == "" {
		return "", errors.New("accessToken empty in credentials file")
	}
	return token, nil
}

// DefaultCredPath returns the default location for credentials.
func DefaultCredPath() string {
	if home, err := os.UserHomeDir(); err == nil {
		return filepath.Join(home, ".claude", ".credentials.json")
	}
	return ".claude/.credentials.json"
}

// HumanTime returns human friendly update time strings.
func HumanTime(t time.Time) string {
	diff := time.Since(t)
	if diff < 5*time.Second {
		return texts.TextUpdatedNow
	}
	if diff < time.Minute {
		secs := int(diff.Seconds())
		if secs < 10 {
			secs = 10
		} else {
			secs = (secs / 10) * 10
		}
		return fmt.Sprintf(texts.TextUpdatedAgo, fmt.Sprintf(texts.TextSecondsFmt, secs))
	}
	if diff < time.Hour {
		return fmt.Sprintf(texts.TextUpdatedAgo, fmt.Sprintf(texts.TextMinutesFmt, int(diff.Minutes())))
	}
	return fmt.Sprintf(texts.TextUpdatedAgo, fmt.Sprintf(texts.TextHourExact, int(diff.Hours())))
}

// FriendlyDuration renders durations in compact form.
func FriendlyDuration(d time.Duration) string {
	if d < time.Minute {
		return texts.TextLtMinute
	}
	if d < time.Hour {
		return fmt.Sprintf(texts.TextMinutesFmt, int(d.Minutes()))
	}
	if d < 24*time.Hour {
		h := int(d.Hours())
		m := int(d.Minutes()) % 60
		if m == 0 {
			return fmt.Sprintf(texts.TextHourExact, h)
		}
		return fmt.Sprintf(texts.TextHourMinute, h, m)
	}
	days := int(d.Hours()) / 24
	h := int(d.Hours()) % 24
	if h == 0 {
		return fmt.Sprintf(texts.TextDaysFmt, days)
	}
	return fmt.Sprintf(texts.TextDaysHours, days, h)
}

// FormatReset produces reset timestamp and remaining duration strings.
func FormatReset(t time.Time) (string, string) {
	if t.IsZero() {
		return "", ""
	}
	local := t.In(time.Local)
	reset := fmt.Sprintf(texts.TextResetAtFmt, local.Format("15:04 MST Jan 02"))
	remain := ""
	if d := time.Until(local); d > 0 {
		remain = fmt.Sprintf(texts.TextRemainFmt, FriendlyDuration(d))
	} else {
		remain = texts.TextResetSoon
	}
	return reset, remain
}

// Clamp constrains v to [minVal, maxVal].
func Clamp(v, minVal, maxVal float64) float64 {
	if v < minVal {
		return minVal
	}
	if v > maxVal {
		return maxVal
	}
	return v
}

// Max returns max of two ints.
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
