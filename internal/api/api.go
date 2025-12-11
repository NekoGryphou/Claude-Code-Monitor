package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"claude-monitor/internal/consts"
)

const (
	apiURL = "https://api.anthropic.com/api/oauth/usage"
)

// HTTPError captures structured details from non-2xx API responses.
type HTTPError struct {
	Status     int
	Body       string
	RetryAfter time.Duration
}

func (e HTTPError) Error() string {
	return fmt.Sprintf(consts.TextHTTPErrorFmt, e.Status, e.Body)
}

// WindowUsage holds utilization for a sliding window and optional reset time.
type WindowUsage struct {
	Utilization *float64 `json:"utilization"`
	ResetsAt    *time.Time
}

// UsageResponse contains rolling usage windows returned by the API.
type UsageResponse struct {
	FiveHour *WindowUsage `json:"five_hour"`
	SevenDay *WindowUsage `json:"seven_day"`
}

// FetchUsage requests utilization data with the given HTTP client and token.
//
// Parameters:
//
//	ctx        - request context; cancellation/timeout is respected.
//	client     - HTTP client to execute the request.
//	token      - OAuth access token.
//	betaHeader - anthropic-beta header value required by the API.
//
// Returns:
//
//	UsageResponse - parsed utilization windows.
//	error         - non-nil on missing token, request failure, or decode error.
func FetchUsage(ctx context.Context, client HTTPClient, token, betaHeader string) (UsageResponse, error) {
	if strings.TrimSpace(token) == "" {
		return UsageResponse{}, errors.New(consts.ErrMissingToken)
	}
	if strings.TrimSpace(betaHeader) == "" {
		return UsageResponse{}, errors.New(consts.ErrBetaHeaderRequired)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return UsageResponse{}, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("anthropic-beta", betaHeader)
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return UsageResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4_096))
		bodyText := strings.TrimSpace(string(body))
		if betaHeader == consts.DefaultBetaName {
			bodyText = bodyText + " (beta header may be outdated; set -beta-header or ANTHROPIC_BETA_HEADER)"
		}
		return UsageResponse{}, HTTPError{
			Status:     resp.StatusCode,
			Body:       bodyText,
			RetryAfter: parseRetryAfter(resp.Header.Get("Retry-After")),
		}
	}

	var payload UsageResponse
	dec := json.NewDecoder(io.LimitReader(resp.Body, 32<<10))
	if err := dec.Decode(&payload); err != nil {
		return UsageResponse{}, err
	}
	return payload, nil
}

// HTTPClient defines the minimal interface required to execute HTTP requests.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// parseRetryAfter converts Retry-After header into a duration when possible.
func parseRetryAfter(raw string) time.Duration {
	if raw == "" {
		return 0
	}
	if secs, err := time.ParseDuration(strings.TrimSpace(raw) + "s"); err == nil {
		return secs
	}
	if ts, err := time.Parse(http.TimeFormat, raw); err == nil {
		return time.Until(ts)
	}
	return 0
}

// UnmarshalJSON parses utilization and optional reset time from API payloads.
//
// Parameters:
//
//	data - raw JSON for a single window usage object.
//
// Returns:
//
//	error - non-nil on malformed JSON; nil otherwise.
func (w *WindowUsage) UnmarshalJSON(data []byte) error {
	var raw struct {
		Utilization *float64 `json:"utilization"`
		ResetsAtRaw *string  `json:"resets_at"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	w.Utilization = raw.Utilization
	if raw.ResetsAtRaw != nil {
		ts, err := time.Parse(time.RFC3339Nano, *raw.ResetsAtRaw)
		if err != nil {
			return fmt.Errorf("resets_at parse: %w", err)
		}
		w.ResetsAt = &ts
	}
	return nil
}
