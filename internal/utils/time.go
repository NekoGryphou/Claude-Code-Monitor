package utils

import (
	"fmt"
	"time"

	"claude-monitor/internal/consts"
)

// HumanTime renders how long ago a timestamp occurred.
//
// Parameters:
//   - t: timestamp to compare against now.
//
// Returns:
//   - human-friendly "updated X ago" text.
func HumanTime(t time.Time) string {
	diff := time.Since(t)
	switch {
	case diff < 5*time.Second:
		return consts.TextUpdatedNow
	case diff < time.Minute:
		secs := int(diff.Seconds())
		if secs < 10 {
			secs = 10
		} else {
			secs = (secs / 10) * 10
		}
		return fmt.Sprintf(consts.TextUpdatedAgo, fmt.Sprintf(consts.TextSecondsFmt, secs))
	case diff < time.Hour:
		return fmt.Sprintf(consts.TextUpdatedAgo, fmt.Sprintf(consts.TextMinutesFmt, int(diff.Minutes())))
	default:
		return fmt.Sprintf(consts.TextUpdatedAgo, fmt.Sprintf(consts.TextHourExact, int(diff.Hours())))
	}
}

// FriendlyDuration formats durations for display.
//
// Parameters:
//   - d: duration to format.
//
// Returns:
//   - concise string representing the duration (e.g., 5m, 2h30m).
func FriendlyDuration(d time.Duration) string {
	switch {
	case d < time.Minute:
		return consts.TextLtMinute
	case d < time.Hour:
		return fmt.Sprintf(consts.TextMinutesFmt, int(d.Minutes()))
	case d < 24*time.Hour:
		h := int(d.Hours())
		m := int(d.Minutes()) % 60
		if m == 0 {
			return fmt.Sprintf(consts.TextHourExact, h)
		}
		return fmt.Sprintf(consts.TextHourMinute, h, m)
	default:
		days := int(d.Hours()) / 24
		h := int(d.Hours()) % 24
		if h == 0 {
			return fmt.Sprintf(consts.TextDaysFmt, days)
		}
		return fmt.Sprintf(consts.TextDaysHours, days, h)
	}
}

// FormatReset builds user-facing reset and remaining strings.
//
// Parameters:
//   - t: reset timestamp.
//
// Returns:
//   - reset string and remaining-duration string (or empty strings when zero time).
func FormatReset(t time.Time) (string, string) {
	if t.IsZero() {
		return "", ""
	}

	local := t.In(time.Local)
	reset := fmt.Sprintf(consts.TextResetAtFmt, local.Format(consts.ResetTimeLayout))

	var remain string
	if d := time.Until(local); d > 0 {
		remain = fmt.Sprintf(consts.TextRemainFmt, FriendlyDuration(d))
	} else {
		remain = consts.TextResetSoon
	}

	return reset, remain
}
