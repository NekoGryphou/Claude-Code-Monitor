package app

import (
	"context"
	"time"

	"claude-monitor/internal/api"
	"claude-monitor/internal/texts"
	"claude-monitor/internal/util"

	tea "github.com/charmbracelet/bubbletea"
)

type chartRow struct {
	label   string
	percent float64
	reset   string
	remain  string
}

type usageMsg struct {
	data api.UsageResponse
	err  error
}

type tickMsg struct{}
type spinMsg struct{}

const spinnerInterval = 120 * time.Millisecond

var spinnerFrames = []string{"⠋", "⠙", "⠸", "⠴", "⠦", "⠇"}

type model struct {
	cfg         Config
	width       int
	rows        []chartRow
	lastUpdated time.Time
	loading     bool
	err         error
	spinIndex   int
}

func initialModel(cfg Config) model {
	return model{
		cfg:         cfg,
		width:       80,
		loading:     true,
		lastUpdated: time.Time{},
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(fetchUsageCmd(m.cfg), spinCmd())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		return m, nil
	case spinMsg:
		if !m.loading {
			return m, nil
		}
		m.spinIndex = (m.spinIndex + 1) % len(spinnerFrames)
		return m, spinCmd()
	case tickMsg:
		return m.handleTick()
	case usageMsg:
		return m.handleUsage(msg)
	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func fetchUsageCmd(cfg Config) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), cfg.RequestTimeout)
		defer cancel()
		data, err := api.FetchUsage(ctx, cfg.HTTPClient, cfg.Token)
		return usageMsg{data: data, err: err}
	}
}

func tickCmd(interval time.Duration) tea.Cmd {
	return tea.Tick(interval, func(time.Time) tea.Msg { return tickMsg{} })
}

func spinCmd() tea.Cmd {
	return tea.Tick(spinnerInterval, func(time.Time) tea.Msg { return spinMsg{} })
}

func (m model) handleTick() (tea.Model, tea.Cmd) {
	m.loading = true
	return m, tea.Batch(fetchUsageCmd(m.cfg), spinCmd())
}

func (m model) handleUsage(msg usageMsg) (tea.Model, tea.Cmd) {
	if msg.err != nil {
		m.err = msg.err
		m.loading = false
		return m, tickCmd(m.cfg.RefreshEvery)
	}
	m.err = nil
	m.rows = buildRows(msg.data)
	m.lastUpdated = time.Now()
	m.loading = false
	return m, tickCmd(m.cfg.RefreshEvery)
}

func (m model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "r", "R":
		m.loading = true
		return m, tea.Batch(fetchUsageCmd(m.cfg), spinCmd())
	default:
		return m, nil
	}
}

func buildRows(u api.UsageResponse) []chartRow {
	return buildChartRows([]struct {
		label string
		win   *api.WindowUsage
	}{
		{label: texts.LabelCurrent, win: u.FiveHour},
		{label: texts.LabelWeekly, win: u.SevenDay},
	})
}

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
			reset, remain = util.FormatReset(*item.win.ResetsAt)
		}
		rows = append(rows, chartRow{
			label:   item.label,
			percent: util.Clamp(*item.win.Utilization, 0, 100),
			reset:   reset,
			remain:  remain,
		})
	}
	return rows
}
