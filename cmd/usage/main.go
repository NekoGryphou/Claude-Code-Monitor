package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"claude-monitor/internal/app"
	"claude-monitor/internal/auth"
	"claude-monitor/internal/consts"
)

// defaultHTTPTimeout is used when no -http-timeout flag or ANTHROPIC_HTTP_TIMEOUT
// environment variable is provided.
const defaultHTTPTimeout = 8 * time.Second

// main parses CLI flags (including beta header and HTTP timeout), resolves the OAuth
// token, builds Config, and starts the UI. It exits with a non-zero status if
// configuration, token resolution, or program execution fails.
func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	refresh := flag.Duration(consts.FlagIntervalName, 30*time.Second, consts.FlagIntervalHelp)
	credPath := flag.String(consts.FlagCredsName, auth.DefaultCredPath(), consts.FlagCredsHelp)

	timeoutDefault, timeoutWarning := loadTimeoutDefault()
	httpTimeout := flag.Duration(consts.FlagTimeoutName, timeoutDefault, consts.FlagTimeoutHelp)
	betaDefault := loadBetaDefault()
	betaHeader := flag.String(consts.FlagBetaName, betaDefault, consts.FlagBetaHelp)
	flag.Parse()

	if timeoutWarning != "" {
		fmt.Fprintln(os.Stderr, timeoutWarning)
	}
	if strings.TrimSpace(*betaHeader) == consts.DefaultBetaName {
		fmt.Fprintln(os.Stderr, "warning: using baked-in beta header; override -beta-header or ANTHROPIC_BETA_HEADER when Anthropic rotates betas")
	}
	token, err := auth.ResolveToken(*credPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, consts.TextTokenErrorFmt+"\n", err)
		os.Exit(1)
	}

	client := newHTTPClient(*httpTimeout)

	cfg := app.Config{
		Token:        token,
		Credentials:  *credPath,
		RefreshEvery: *refresh,
		HTTPClient:   client,
		BetaHeader:   strings.TrimSpace(*betaHeader),
	}

	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, consts.TextConfigErrFmt+"\n", err)
		os.Exit(1)
	}

	if err := app.Run(ctx, cfg); err != nil {
		fmt.Fprintf(os.Stderr, consts.TextAppErrFmt+"\n", err)
		os.Exit(1)
	}
}

func loadTimeoutDefault() (time.Duration, string) {
	if v := strings.TrimSpace(os.Getenv(consts.EnvHTTPTimeout)); v != "" {
		if d, err := time.ParseDuration(v); err == nil && d > 0 {
			return d, ""
		}
		return defaultHTTPTimeout, fmt.Sprintf("warning: invalid %s value %q; using default %s", consts.EnvHTTPTimeout, v, defaultHTTPTimeout)
	}
	return defaultHTTPTimeout, ""
}

// loadBetaDefault returns the beta header from environment or the baked-in
// default that ships with the binary.
func loadBetaDefault() string {
	if v := strings.TrimSpace(os.Getenv(consts.EnvBetaHeader)); v != "" {
		return v
	}
	return consts.DefaultBetaName
}

func newHTTPClient(timeout time.Duration) *http.Client {
	tr, ok := http.DefaultTransport.(*http.Transport)
	if ok && tr != nil {
		tr = tr.Clone()
		tr.MaxIdleConns = 32
		tr.MaxIdleConnsPerHost = 8
		tr.IdleConnTimeout = 30 * time.Second
		if timeout > 0 {
			tr.ResponseHeaderTimeout = timeout
		}
		return &http.Client{Timeout: timeout, Transport: tr}
	}
	return &http.Client{Timeout: timeout}
}
