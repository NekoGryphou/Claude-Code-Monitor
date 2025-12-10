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
)

const (
	apiURL        = "https://api.anthropic.com/api/oauth/usage"
	betaHeaderVal = "oauth-2025-04-20"
)

// WindowUsage holds utilization for a sliding window.
// It implements custom JSON parsing so we keep only the parsed timestamp.
type WindowUsage struct {
	Utilization *float64 `json:"utilization"`
	ResetsAt    *time.Time
}

type UsageResponse struct {
	FiveHour *WindowUsage `json:"five_hour"`
	SevenDay *WindowUsage `json:"seven_day"`
}

func FetchUsage(ctx context.Context, client HTTPClient, token string) (UsageResponse, error) {
	if strings.TrimSpace(token) == "" {
		return UsageResponse{}, errors.New("missing OAuth token")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return UsageResponse{}, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("anthropic-beta", betaHeaderVal)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return UsageResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4_096))
		return UsageResponse{}, fmt.Errorf("http %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var payload UsageResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return UsageResponse{}, err
	}
	return payload, nil
}

// HTTPClient matches http.Client's Do method for easy substitution in tests.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

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
		if ts, err := time.Parse(time.RFC3339Nano, *raw.ResetsAtRaw); err == nil {
			w.ResetsAt = &ts
		}
	}
	return nil
}
