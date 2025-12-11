package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"claude-monitor/internal/api"
	"claude-monitor/internal/consts"
)

// formatError prettifies API or transport errors for display.
//
// Parameters:
//   - err: error returned from usage fetch.
//
// Returns:
//   - concise, user-facing error message.
func formatError(err error) string {
	msg := strings.TrimSpace(err.Error())

	if errors.Is(err, context.DeadlineExceeded) {
		return consts.TextRequestTimedOut
	}
	if errors.Is(err, context.Canceled) {
		return consts.TextRequestCanceled
	}

	status := ""
	body := msg

	var httpErr api.HTTPError
	if errors.As(err, &httpErr) {
		status = fmt.Sprintf("http %d", httpErr.Status)
		body = strings.TrimSpace(httpErr.Body)
	} else if parts := strings.SplitN(msg, ":", 2); len(parts) == 2 && strings.HasPrefix(parts[0], "http ") {
		status = strings.TrimSpace(parts[0])
		body = strings.TrimSpace(parts[1])
	}

	var payload struct {
		Error struct {
			Type      string `json:"type"`
			Message   string `json:"message"`
			ErrorCode string `json:"error_code"`
		} `json:"error"`
		Message   string `json:"message"`
		ErrorCode string `json:"error_code"`
	}

	summary := body
	code := ""

	if len(body) > 0 && (body[0] == '{' || body[0] == '[') {
		if err := json.Unmarshal([]byte(body), &payload); err == nil {
			switch {
			case payload.Error.Message != "":
				summary = payload.Error.Message
				code = payload.Error.ErrorCode
			case payload.Message != "":
				summary = payload.Message
				code = payload.ErrorCode
			}
		}
	}

	parts := make([]string, 0, 3)
	if status != "" {
		parts = append(parts, status)
	}
	if summary != "" {
		parts = append(parts, summary)
	}
	if code != "" {
		parts = append(parts, fmt.Sprintf("code: %s", code))
	}

	result := strings.Join(parts, " · ")
	return truncateString(result, 120)
}

// truncateString shortens s to maxLen runes, appending ellipsis when trimmed.
//
// Parameters:
//   - s: string to trim.
//   - maxLen: maximum rune count.
//
// Returns:
//   - possibly truncated string with ellipsis when needed.
func truncateString(s string, maxLen int) string {
	if maxLen <= 0 {
		return ""
	}
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	if maxLen <= 1 {
		return string(runes[:maxLen])
	}
	return string(runes[:maxLen-1]) + "…"
}
