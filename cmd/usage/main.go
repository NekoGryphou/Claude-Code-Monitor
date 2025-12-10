package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"claude-monitor/internal/app"
	"claude-monitor/internal/texts"
	"claude-monitor/internal/util"
)

func main() {
	refresh := flag.Duration("interval", 30*time.Second, texts.FlagIntervalHelp)
	credPath := flag.String("creds", util.DefaultCredPath(), texts.FlagCredsHelp)
	timeout := flag.Duration("timeout", 8*time.Second, texts.FlagTimeoutHelp)
	flag.Parse()

	token, err := util.ResolveToken(*credPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "token error: %v\n", err)
		os.Exit(1)
	}

	cfg := app.Config{
		Token:          token,
		Credentials:    *credPath,
		RefreshEvery:   *refresh,
		HTTPClient:     &http.Client{Timeout: *timeout},
		RequestTimeout: *timeout,
	}

	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "config error: %v\n", err)
		os.Exit(1)
	}

	if err := app.Run(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "app error: %v\n", err)
		os.Exit(1)
	}
}
