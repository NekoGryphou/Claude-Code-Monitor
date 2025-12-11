# Claude Monitor

Terminal dashboard that keeps an eye on your Anthropic OAuth usage. It polls the official usage endpoint and shows five-hour and seven-day utilization with a small Bubble Tea UI.

## Features
- Live utilization bars for the 5‑hour and 7‑day windows, with reset time and remaining window shown underneath.
- Auto-refreshes on a timer; press `r` to fetch immediately, `q` or `ctrl+c` to exit.
- Compact lipgloss styling, spinner while loading, and friendly “last updated” text.
- Reads your OAuth token from an env var or the same credentials file used by the Claude desktop app.
- Optional high-contrast / colorless modes via `CLAUDE_MONITOR_HIGH_CONTRAST=1` or `NO_COLOR`.

## Requirements
- Go 1.22 or newer.
- Network access to `https://api.anthropic.com/api/oauth/usage`.
- An OAuth access token supplied via `ANTHROPIC_OAUTH_TOKEN` or a credentials file.
- Anthropic beta header value; defaults to `oauth-2025-04-20` but is configurable (see flags below). This header changes over time—expect to override it when Anthropic rotates betas.

## Getting the token
- **Env var (preferred for CI):** export `ANTHROPIC_OAUTH_TOKEN="your-token"`.
- **Credentials file:** default path is `~/.claude/.credentials.json` (override with `-creds`). Expected shape:
```json
{
  "claudeAiOauth": {
    "accessToken": "your-token-here"
  }
}
```

## Build & Run
- Quick run without installing: `go run ./cmd/usage`
- Build a reusable binary: `go build -o bin/claude-monitor ./cmd/usage`
- Install to `$GOBIN`: `go install ./cmd/usage`

### Precedence & behavior
- OAuth token: `ANTHROPIC_OAUTH_TOKEN` wins; otherwise the credentials file is read.
- Timeout: `ANTHROPIC_HTTP_TIMEOUT` overrides the flag default; invalid values are ignored with a warning and the 8s default is used.
- Beta header: `ANTHROPIC_BETA_HEADER` overrides the compiled default; set this explicitly if the API starts returning 401/403 with the baked-in value.

### CLI flags
- `-interval` poll cadence, e.g. `15s` or `1m` (default `30s`)
- `-creds` path to credentials JSON when not using `ANTHROPIC_OAUTH_TOKEN` (default `~/.claude/.credentials.json`)
- `-http-timeout` request timeout (default 8s; overrideable via `ANTHROPIC_HTTP_TIMEOUT`)
- `-beta-header` Anthropic beta header value (default `oauth-2025-04-20`; overrideable via `ANTHROPIC_BETA_HEADER`)

Requests time out using the configured HTTP timeout (or the refresh interval, whichever is shorter) to avoid overlapping polls.

Example: `./bin/claude-monitor -interval 20s`

> Heads up: the baked-in beta header will expire when Anthropic rotates betas. Prefer setting `ANTHROPIC_BETA_HEADER` or `-beta-header` explicitly, especially if you see 401/403 responses.

## Reading the UI
- “Current” is the rolling 5‑hour utilization; “Weekly” is the rolling 7‑day utilization.
- Bars clamp between 0–100%. If the API omits a window, that row is hidden.
- Reset timestamps are shown in your local time with a “left” indicator until the window rolls over.

## Credentials permissions
If you use the credentials file (`~/.claude/.credentials.json` by default), it must be owner-only readable (`chmod 600`). The tool refuses to load world- or group-readable files to avoid leaking OAuth tokens.

## Project layout
- `cmd/usage` — CLI entrypoint and flag parsing.
- `internal/app` — Bubble Tea model, view, styling, and layout helpers.
- `internal/api` — Minimal client for the Anthropic OAuth usage endpoint.
- `internal/auth` — Credential resolution from env or credentials file.
- `internal/utils` — Small helpers for math, time formatting, etc.

## Troubleshooting
- “token error”: env var missing or credentials file unreadable/empty.
- “http <code>”: API rejected the request (check token validity and beta header requirements). If you see 401s with the default beta value, supply a current header via `-beta-header` or `ANTHROPIC_BETA_HEADER`.
- The UI will keep running after an error and retry on the next interval; press `r` to retry immediately.
- Persistent 429/5xx: the app backs off exponentially up to 8× the interval; consider increasing interval or updating the beta header.
