package texts

// UI labels and headers.
const (
	HeaderTitle       = "Claude code Usage"
	TextNoData        = "No utilization data available."
	TextSeparatorDot  = " · "
	TextIntervalFmt   = "interval %s"
	TextStatusFetch   = "%s fetching latest…"
	TextStatusWaiting = "waiting for first sample…"
	TextErrorFmt      = "error: %v"

	LabelCurrent = "Current"
	LabelWeekly  = "Weekly"

	MarkArt = " ▐▛███▜▌ \n▝▜█████▛▘\n  ▘▘ ▝▝"

	FlagIntervalHelp = "poll interval (e.g. 15s, 1m)"
	FlagCredsHelp    = "path to credentials JSON (uses ANTHROPIC_OAUTH_TOKEN if set)"
	FlagTimeoutHelp  = "per-request timeout"
)

// Time and reset phrasing.
const (
	TextSecondsFmt = "%ds"
	TextLtMinute   = "less than a minute"
	TextMinutesFmt = "%dm"
	TextHourExact  = "%dh"
	TextHourMinute = "%dh%dm"
	TextDaysFmt    = "%dd"
	TextDaysHours  = "%dd %dh"
	TextResetAtFmt = "resets at %s"
	TextResetSoon  = "resets soon"
	TextRemainFmt  = "%s left"
	TextUpdatedNow = "updated right now"
	TextUpdatedAgo = "updated %s ago"
)
