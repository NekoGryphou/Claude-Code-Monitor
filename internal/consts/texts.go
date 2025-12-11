package consts

// User-facing copy and keys used across the application.
const (
	// HeaderTitle is the banner title shown in the UI header.
	HeaderTitle = "Claude OAuth Usage"
	// TextNoData indicates no utilization was returned.
	TextNoData = "No utilization data available."
	// TextSeparatorDot is the middle dot separator for inline lists.
	TextSeparatorDot = " · "
	// TextIntervalFmt formats the refresh interval in the footer.
	TextIntervalFmt = "interval %s"
	// TextStatusFetch shows while a request is running.
	TextStatusFetch = "%s fetching latest…"
	// TextStatusWaiting shows before the first sample arrives.
	TextStatusWaiting = "waiting for first sample…"
	// TextErrorFmt prefixes fatal errors from main.
	TextErrorFmt = "error: %v"
	// TextHTTPErrorFmt formats HTTP status and body on failure.
	TextHTTPErrorFmt = "http %d: %s"
	// TextTokenErrorFmt formats token resolution errors.
	TextTokenErrorFmt = "token error: %v"
	// TextConfigErrFmt formats configuration validation errors.
	TextConfigErrFmt = "config error: %v"
	// TextAppErrFmt formats Bubble Tea runtime errors.
	TextAppErrFmt = "app error: %v"

	// LabelCurrent is the row label for 5-hour usage.
	LabelCurrent = "Current"
	// LabelWeekly is the row label for 7-day usage.
	LabelWeekly = "Weekly"

	// EnvBetaHeader names the env var for the Anthropic beta header.
	EnvBetaHeader = "ANTHROPIC_BETA_HEADER"
	// EnvHTTPTimeout names the env var for request timeout.
	EnvHTTPTimeout = "ANTHROPIC_HTTP_TIMEOUT"
	// DefaultBetaName is the baked-in default beta header value.
	DefaultBetaName = "oauth-2025-04-20"

	// MarkArt is the ASCII logotype block in the header.
	MarkArt = " ▐▛███▜▌ \n▝▜█████▛▘\n  ▘▘ ▝▝"

	// FlagIntervalName is the CLI flag name for refresh cadence.
	FlagIntervalName = "interval"
	// FlagCredsName is the CLI flag name for credentials path.
	FlagCredsName = "creds"
	// FlagTimeoutName is the CLI flag name for HTTP timeout.
	FlagTimeoutName = "http-timeout"
	// FlagBetaName is the CLI flag name for beta header value.
	FlagBetaName = "beta-header"
	// FlagIntervalHelp describes the interval flag.
	FlagIntervalHelp = "poll interval (e.g. 15s, 1m)"
	// FlagCredsHelp describes the creds flag.
	FlagCredsHelp = "path to credentials JSON (uses ANTHROPIC_OAUTH_TOKEN if set)"
	// FlagTimeoutHelp describes the HTTP timeout flag.
	FlagTimeoutHelp = "HTTP timeout (e.g. 5s, 2s)"
	// FlagBetaHelp describes the beta header flag.
	FlagBetaHelp = "Anthropic beta header value"

	// HelpRefreshKey is the lowercase key to refresh now.
	HelpRefreshKey = "r"
	// HelpRefreshKeyUpper is the uppercase key to refresh now.
	HelpRefreshKeyUpper = "R"
	// HelpQuitKey is the key to quit the app.
	HelpQuitKey = "q"
	// HelpQuitCtrlKey is the ctrl key combo to quit.
	HelpQuitCtrlKey = "ctrl+c"
	// HelpRefreshDesc describes the refresh shortcut.
	HelpRefreshDesc = "refresh now"
	// HelpQuitDesc describes the quit shortcut.
	HelpQuitDesc = "quit"
	// HelpKeyJoiner joins multiple keys in help text.
	HelpKeyJoiner = "/"
	// HelpSpacer separates key and description.
	HelpSpacer = " "

	// SpinnerSamplePercent provides width for spinner/value alignment.
	SpinnerSamplePercent = "100.0%"
	// PercentFmt formats utilization percentages.
	PercentFmt = "%5.1f%%"

	// EnvTokenName names the env var for the OAuth token.
	EnvTokenName = "ANTHROPIC_OAUTH_TOKEN"
	// DefaultCredRelPath is the default credentials path relative to home.
	DefaultCredRelPath = ".claude/.credentials.json"
	// TildePrefix marks a path that should expand to the home directory.
	TildePrefix = "~"
	// ResetTimeLayout formats reset timestamps for display.
	ResetTimeLayout = "15:04 MST Jan 02"

	// ErrTokenRequired signals missing token in config.
	ErrTokenRequired = "token required"
	// ErrHTTPClientRequired signals missing HTTP client in config.
	ErrHTTPClientRequired = "http client required"
	// ErrHTTPClientTimeout signals invalid HTTP client timeout.
	ErrHTTPClientTimeout = "http client timeout must be positive"
	// ErrRefreshInterval signals invalid refresh cadence.
	ErrRefreshInterval = "refresh interval must be positive"
	// ErrMissingToken signals missing OAuth token before request.
	ErrMissingToken = "missing OAuth token"
	// ErrBetaHeaderRequired signals missing beta header value.
	ErrBetaHeaderRequired = "beta header required"
	// ErrReadCredentialsFmt formats credential read failures.
	ErrReadCredentialsFmt = "read credentials: %w"
	// ErrParseCredentialsFmt formats credential parse failures.
	ErrParseCredentialsFmt = "parse credentials: %w"
	// ErrEmptyAccessToken signals empty token inside the credentials file.
	ErrEmptyAccessToken = "accessToken empty in credentials file"

	// TextRequestTimedOut is shown when a request exceeds its deadline.
	TextRequestTimedOut = "request timed out"
	// TextRequestCanceled is shown when a request is canceled.
	TextRequestCanceled = "request canceled"
	// TextSkeletonReset is placeholder reset text in the loading skeleton.
	TextSkeletonReset = "resets at …"
	// TextSkeletonLeft is placeholder remaining text in the loading skeleton.
	TextSkeletonLeft = "... left"
)

// Time and duration formatting strings used in the UI.
const (
	// TextSecondsFmt formats seconds for “updated ago” messages.
	TextSecondsFmt = "%ds"
	// TextLtMinute indicates a duration under one minute.
	TextLtMinute = "less than a minute"
	// TextMinutesFmt formats whole minutes.
	TextMinutesFmt = "%dm"
	// TextHourExact formats whole hours.
	TextHourExact = "%dh"
	// TextHourMinute formats hours and remaining minutes.
	TextHourMinute = "%dh%dm"
	// TextDaysFmt formats whole days.
	TextDaysFmt = "%dd"
	// TextDaysHours formats days plus hours.
	TextDaysHours = "%dd %dh"
	// TextResetAtFmt formats the reset timestamp.
	TextResetAtFmt = "resets at %s"
	// TextResetSoon indicates a window has already rolled over.
	TextResetSoon = "resets soon"
	// TextRemainFmt formats remaining time in a window.
	TextRemainFmt = "%s left"
	// TextUpdatedNow indicates a very recent update.
	TextUpdatedNow = "updated right now"
	// TextUpdatedAgo formats time since last update.
	TextUpdatedAgo = "updated %s ago"
)
