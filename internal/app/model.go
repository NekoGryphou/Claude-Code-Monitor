package app

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"claude-monitor/internal/api"
	"claude-monitor/internal/consts"
	"claude-monitor/internal/utils"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// chartRow represents a single usage window row rendered in the UI.
type chartRow struct {
	label   string
	percent float64
	reset   string
	remain  string
}

// usageMsg wraps usage data or an error returned from the API request.
type usageMsg struct {
	data api.UsageResponse
	err  error
}

// model holds all Bubble Tea state for the application.
type model struct {
	cfg         Config
	baseCtx     context.Context
	width       int
	rows        []chartRow
	lastUpdated time.Time
	loading     bool
	err         error
	sp          spinner.Model
	cancel      context.CancelFunc
	failures    int
}

// tickMsg signals that the refresh interval elapsed.
type tickMsg struct{}

// initialModel builds the starting model state for the UI.
//
// Params:
//   - cfg: validated Config used to seed model state and refresh cadence.
//
// Returns:
//   - a model with spinner initialized, loading set, and default width.
func initialModel(ctx context.Context, cfg Config) model {
	return model{
		cfg:         cfg,
		baseCtx:     ctx,
		width:       80,
		loading:     false,
		lastUpdated: time.Time{},
		sp:          newSpinner(),
	}
}

// newSpinner constructs a spinner with the application's base styling.
func newSpinner() spinner.Model {
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = spinnerStyle
	return sp
}

// Init registers the initial commands (usage fetch and spinner tick).
//
// Returns:
//
//	tea.Cmd - batch command to kick off data fetch and spinner.
func (m model) Init() tea.Cmd {
	return tea.Batch(tickCmd(0), m.sp.Tick)
}

// Update routes incoming messages to state handlers and returns the next command.
//
// Params:
//   - msg: Bubble Tea message to process (window size, spinner tick, data, key, or tick).
//
// Returns:
//   - updated model reflecting the message.
//   - a command to continue program flow (or nil).
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		return m, nil
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.sp, cmd = m.sp.Update(msg)
		if m.loading {
			return m, cmd
		}
		return m, nil
	case tickMsg:
		return m.handleTick()
	case usageMsg:
		return m.handleUsage(msg)
	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

// fetchUsageCmd queries usage with the configured HTTP client and token.
//
// Params:
//   - cfg: provides HTTP client, timeout, and token used for the request.
//   - ctx: deadline/timeout context for the call.
//
// Returns:
//   - a command that fetches usage and emits usageMsg containing data or error.
func fetchUsageCmd(cfg Config, ctx context.Context, cancel context.CancelFunc) tea.Cmd {
	return func() tea.Msg {
		defer cancel()
		data, err := api.FetchUsage(ctx, cfg.HTTPClient, cfg.Token, cfg.BetaHeader)
		return usageMsg{data: data, err: err}
	}
}

// tickCmd schedules the next periodic refresh.
//
// Params:
//   - interval: duration before the next tick fires.
//
// Returns:
//   - a command that will send tickMsg after interval elapses.
func tickCmd(interval time.Duration) tea.Cmd {
	return tea.Tick(interval, func(time.Time) tea.Msg { return tickMsg{} })
}

// effectiveTimeout picks the smaller of the HTTP client timeout and refresh
// interval (when the interval is positive) to keep each request bounded by the
// poll cadence. This prevents a single slow request from blocking multiple
// refresh cycles.
func effectiveTimeout(cfg Config) time.Duration {
	timeout := cfg.HTTPClient.Timeout
	if cfg.RefreshEvery > 0 && cfg.RefreshEvery < timeout {
		return cfg.RefreshEvery
	}
	return timeout
}

// newRequestContext builds a request-scoped context bounded by the smaller of
// the HTTP client timeout and the refresh interval.
func (m model) newRequestContext() (context.Context, context.CancelFunc) {
	timeout := effectiveTimeout(m.cfg)
	if timeout <= 0 {
		timeout = m.cfg.HTTPClient.Timeout
	}
	if timeout <= 0 {
		timeout = 2 * time.Second
	}
	parent := m.baseCtx
	if parent == nil {
		parent = context.Background()
	}
	return context.WithTimeout(parent, timeout)
}

func retryInterval(base time.Duration, failures int, retryAfter time.Duration) time.Duration {
	if retryAfter > 0 {
		return retryAfter
	}
	if failures <= 0 || base <= 0 {
		return base
	}
	step := failures
	if step > 3 {
		step = 3
	}
	backoff := base * time.Duration(1<<step)
	maxBackoff := base * 8
	if backoff > maxBackoff {
		backoff = maxBackoff
	}
	jitter := time.Duration(rand.Int63n(int64(backoff / 5))) // up to +20%
	return backoff + jitter
}

// handleTick triggers a refresh cycle.
//
// Returns:
//   - the model with loading set to true.
//   - a batch command to refetch usage and continue spinner animation.
func (m model) handleTick() (tea.Model, tea.Cmd) {
	// Avoid overlapping fetches when a previous request is still in flight.
	if m.loading {
		return m, nil
	}
	return m.startFetch()
}

// handleUsage ingests fetched usage data or records an error.
//
// Params:
//   - msg: usageMsg carrying API data or an error.
//
// Returns:
//   - the updated model with rows/lastUpdated or error set.
//   - a command scheduling the next tick based on refresh interval.
func (m model) handleUsage(msg usageMsg) (tea.Model, tea.Cmd) {
	if msg.err != nil {
		m.err = msg.err
		m.loading = false
		m.failures++
		var retryAfter time.Duration
		var httpErr api.HTTPError
		if errors.As(msg.err, &httpErr) {
			retryAfter = httpErr.RetryAfter
		}
		return m, tickCmd(retryInterval(m.cfg.RefreshEvery, m.failures, retryAfter))
	}
	m.failures = 0
	m.err = nil
	m.rows = buildRows(msg.data)
	m.lastUpdated = time.Now()
	m.loading = false
	if len(m.rows) == 0 {
		m.err = errors.New(consts.TextNoData)
	}
	return m, tickCmd(m.cfg.RefreshEvery)
}

// handleKey processes user input shortcuts.
//
// Params:
//   - msg: key message containing the pressed key.
//
// Returns:
//   - the model (possibly reset to loading).
//   - a command to quit, refetch, or no-op based on the key.
func (m model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case consts.HelpQuitKey, consts.HelpQuitCtrlKey:
		if m.cancel != nil {
			m.cancel()
		}
		return m, tea.Quit
	case consts.HelpRefreshKey, consts.HelpRefreshKeyUpper:
		if m.loading {
			return m, nil
		}
		return m.startFetch()
	default:
		return m, nil
	}
}

// startFetch marks the model as loading and fires a new usage request while
// continuing the spinner animation.
func (m model) startFetch() (tea.Model, tea.Cmd) {
	m.loading = true
	m.err = nil
	if m.cancel != nil {
		m.cancel()
	}
	ctx, cancel := m.newRequestContext()
	m.cancel = cancel
	return m, tea.Batch(fetchUsageCmd(m.cfg, ctx, cancel), m.sp.Tick)
}

// buildRows constructs chart rows for each usage window.
//
// Params:
//   - u: usage response from the API.
//
// Returns:
//   - slice of chartRow for windows that contain utilization data.
func buildRows(u api.UsageResponse) []chartRow {
	return buildChartRows([]struct {
		label string
		win   *api.WindowUsage
	}{
		{label: consts.LabelCurrent, win: u.FiveHour},
		{label: consts.LabelWeekly, win: u.SevenDay},
	})
}

// buildChartRows normalizes window usage items into chartRow slices.
//
// Params:
//   - items: labeled window usage pointers.
//
// Returns:
//   - chart rows containing clamped utilization and formatted reset/remain text.
func buildChartRows(items []struct {
	label string
	win   *api.WindowUsage
}) []chartRow {
	rows := make([]chartRow, 0, len(items))
	for _, item := range items {
		if item.win == nil || item.win.Utilization == nil {
			continue
		}
		reset, remain := "", ""
		if item.win.ResetsAt != nil {
			reset, remain = utils.FormatReset(*item.win.ResetsAt)
		}
		rows = append(rows, chartRow{
			label:   item.label,
			percent: utils.Clamp(*item.win.Utilization, 0, 100),
			reset:   reset,
			remain:  remain,
		})
	}
	return rows
}
