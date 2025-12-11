package app

import (
	"fmt"
	"math"
	"strings"
	"sync"
	"time"

	"claude-monitor/internal/consts"
	"claude-monitor/internal/utils"

	"github.com/charmbracelet/lipgloss"
)

const (
	minContainerWidth = 20
	horizontalPadding = 2
)

var (
	headerOnce     sync.Once
	headerStatic   string
	helpOnce       sync.Once
	helpStatic     string
	helpTextWidth  int
	valueTextOnce  sync.Once
	valueTextWidth int
)

// layout represents calculated widths for the outer container and inner content.
type layout struct {
	containerWidth int
	contentWidth   int
}

// barMetrics captures computed widths for label, value, and bar columns.
type barMetrics struct {
	labelWidth int
	valueWidth int
	barWidth   int
}

// barRenderOptions bundles styling and formatting hooks for rendering bars.
type barRenderOptions struct {
	labelStyle     lipgloss.Style
	valueStyle     lipgloss.Style
	resetStyle     *lipgloss.Style
	remainStyle    *lipgloss.Style
	barFillStyle   lipgloss.Style
	barEmptyStyle  lipgloss.Style
	valueFormatter func(float64) string
	metaBuilder    func(r chartRow, metrics barMetrics, opt barRenderOptions) string
}

func headerCached() string {
	headerOnce.Do(func() {
		headerStatic = renderHeader()
	})
	return headerStatic
}

func helpCached() (string, int) {
	helpOnce.Do(func() {
		helpStatic = renderHelp()
		helpTextWidth = lipgloss.Width(helpStatic)
	})
	return helpStatic, helpTextWidth
}

func valueWidthCached() int {
	valueTextOnce.Do(func() {
		valueTextWidth = lipgloss.Width(valueBaseStyle.Render(consts.SpinnerSamplePercent))
	})
	return valueTextWidth
}

// newLayout computes container and content widths based on the window size.
//
// Parameters:
//
//	windowWidth - current terminal width in cells.
//
// Returns:
//
//	layout - sizing information used when rendering the frame.
func newLayout(windowWidth int) layout {
	containerWidth := utils.Max(windowWidth, minContainerWidth)
	contentWidth := utils.Max(containerWidth-(horizontalPadding*2), 8)
	return layout{
		containerWidth: containerWidth,
		contentWidth:   contentWidth,
	}
}

// View renders the full TUI frame for the current model state.
//
// Returns:
//
//	string - ANSI-styled layout containing header, body, and footer.
func (m model) View() string {
	frame := newLayout(m.width)

	header := headerCached()
	body := renderBody(frame, m)
	helpText, helpWidth := helpCached()
	footer := renderFooter(frame.contentWidth, helpText, helpWidth, renderStatus(m), m.cfg.RefreshEvery)

	content := lipgloss.JoinVertical(lipgloss.Left, header, body, footer)

	return pageStyle.
		Width(frame.containerWidth).
		Render(content)
}

// renderTitle renders the header title text with its styling.
func renderTitle() string {
	return headerStyle.
		Render(consts.HeaderTitle)
}

// renderMark returns the ASCII art logo block.
func renderMark() string {
	return markStyle.
		Render(consts.MarkArt)
}

// renderHeader composes the mark and title for the page header.
func renderHeader() string {
	mark := renderMark()

	return lipgloss.JoinHorizontal(lipgloss.Top, mark, renderTitle())
}

// renderBody produces the chart area and any error state.
//
// Parameters:
//
//	frame - layout sizing constraints.
//	m     - current model containing rows and potential error.
//
// Returns:
//
//	string - rendered body content.
func renderBody(frame layout, m model) string {
	if m.err != nil {
		return errorBox(formatError(m.err), frame.contentWidth)
	}

	sections := make([]string, 0, 2)

	switch {
	case len(m.rows) == 0:
		if m.loading || m.lastUpdated.IsZero() {
			sections = append(sections, renderSkeleton(frame.contentWidth))
		} else {
			sections = append(sections, consts.TextNoData)
		}
	default:
		const chartFrame = 6
		chartWidth := utils.Max(12, frame.contentWidth-2)
		innerWidth := utils.Max(4, chartWidth-chartFrame)
		sections = append(sections,
			chartBoxStyle.Width(chartWidth).Render(renderBars(m.rows, innerWidth)))
	}

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// errorBox renders a bordered error container sized to the provided width.
//
// Parameters:
//
//	msg   - error message to display.
//	width - available width for the error container.
//
// Returns:
//
//	string - rendered error box.
func errorBox(msg string, width int) string {
	boxW := utils.Max(16, width)
	contentW := utils.Max(8, boxW-4)

	content := errorContentStyleBase.
		Width(contentW).
		Render(msg)

	return errorBoxStyleBase.
		Width(boxW).
		Render(content)
}

// renderStatus builds the status line content.
//
// Parameters:
//
//	m - current model containing loading/error timing info.
//
// Returns:
//
//	string - spinner text, humanized timestamp, or waiting message.
func renderStatus(m model) string {
	switch {
	case m.loading:
		return fmt.Sprintf(consts.TextStatusFetch, m.sp.View())
	case !m.lastUpdated.IsZero():
		return utils.HumanTime(m.lastUpdated)
	default:
		return consts.TextStatusWaiting
	}
}

// renderHelp builds the keybinding legend.
//
// Returns:
//
//	string - inline help showing the available shortcuts.
func renderHelp() string {
	shortcuts := []struct {
		keys []string
		desc string
	}{
		{keys: []string{consts.HelpRefreshKey}, desc: consts.HelpRefreshDesc},
		{keys: []string{consts.HelpQuitKey, consts.HelpQuitCtrlKey}, desc: consts.HelpQuitDesc},
	}

	parts := make([]string, 0, len(shortcuts))
	for _, sc := range shortcuts {
		key := helpKeyStyle.Render(strings.Join(sc.keys, consts.HelpKeyJoiner))
		desc := helpDescStyle.Render(sc.desc)
		parts = append(parts, lipgloss.JoinHorizontal(lipgloss.Top, key, consts.HelpSpacer, desc))
	}
	return helpStyle.Render(strings.Join(parts, consts.TextSeparatorDot))
}

// renderFooter composes help text and right-aligned status/interval line.
//
// Parameters:
//
//	width      - total available footer width.
//	helpText   - pre-rendered help legend.
//	helpWidth  - width of the help legend.
//	status     - status string (fetching or last updated).
//	interval   - refresh interval to display.
//
// Returns:
//
//	string - rendered footer line.
func renderFooter(width int, helpText string, helpWidth int, status string, interval time.Duration) string {
	intervalText := fmt.Sprintf(consts.TextIntervalFmt, interval)
	rightContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		statusStyle.Render(status),
		separatorStyle.Render(consts.TextSeparatorDot),
		statusStyle.Render(intervalText),
	)
	rightW := lipgloss.Width(rightContent)
	lineWidth := utils.Max(width, helpWidth+rightW)
	right := lipgloss.PlaceHorizontal(lineWidth-helpWidth, lipgloss.Right, rightContent)

	return footerStyle.
		Width(lineWidth).
		Render(lipgloss.JoinHorizontal(lipgloss.Top, helpText, right))
}

// renderBars draws utilization bars and optional metadata lines.
//
// Parameters:
//
//	rows       - chart data to render.
//	totalWidth - width available for labels, bar, and value.
//
// Returns:
//
//	string - vertical composition of rendered bars.
func renderBars(rows []chartRow, totalWidth int) string {
	return renderBarsWithOptions(rows, totalWidth, barRenderOptions{
		labelStyle:    labelBaseStyle,
		valueStyle:    valueBaseStyle,
		resetStyle:    &resetBaseStyle,
		remainStyle:   &remainBaseStyle,
		barFillStyle:  barFillStyle,
		barEmptyStyle: barEmptyStyle,
		valueFormatter: func(p float64) string {
			return fmt.Sprintf(consts.PercentFmt, p)
		},
		metaBuilder: defaultMetaBuilder,
	})
}

// renderBarsWithOptions renders bars with custom styles, value formatting,
// and metadata composition.
//
// Parameters:
//   - rows: data to render.
//   - totalWidth: available width for labels, bars, and values.
//   - opt: styling and formatting hooks.
//
// Returns:
//   - vertical ANSI-rendered string of the bar blocks.
func renderBarsWithOptions(rows []chartRow, totalWidth int, opt barRenderOptions) string {
	if len(rows) == 0 {
		return ""
	}

	if opt.valueFormatter == nil {
		opt.valueFormatter = func(p float64) string {
			return fmt.Sprintf(consts.PercentFmt, p)
		}
	}
	if opt.metaBuilder == nil {
		opt.metaBuilder = defaultMetaBuilder
	}

	metrics := computeBarMetrics(totalWidth, rows, valueWidthCached())

	blocks := make([]string, 0, len(rows))
	separator := barSeparatorStyle.Render(" ")
	labelStyle := opt.labelStyle.Width(metrics.labelWidth)
	valueStyle := opt.valueStyle.
		Width(metrics.valueWidth).
		AlignHorizontal(lipgloss.Right)
	space := "  "
	metaStyle := metaBaseStyle.MarginLeft(metrics.labelWidth + 1)

	for i, r := range rows {
		percent := utils.Clamp(r.percent, 0, 100)
		bar := renderProgressBarStyled(metrics.barWidth, percent/100, opt.barFillStyle, opt.barEmptyStyle)
		value := valueStyle.Render(opt.valueFormatter(percent))

		left := lipgloss.JoinHorizontal(lipgloss.Top,
			labelStyle.Render(r.label),
			separator,
			bar,
		)
		line := lipgloss.JoinHorizontal(lipgloss.Top, left+space, value)

		meta := opt.metaBuilder(r, metrics, opt)

		rendered := line
		if meta != "" {
			metaLine := metaStyle.Render(meta)
			rendered = lipgloss.JoinVertical(lipgloss.Left, line, metaLine)
		}

		if i == len(rows)-1 {
			blocks = append(blocks, barLastBlockStyle.Render(rendered))
		} else {
			blocks = append(blocks, barBlockStyle.Render(rendered))
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, blocks...)
}

// defaultMetaBuilder joins reset/remaining strings for a chart row.
//
// Parameters:
//   - r: chart row containing reset/remain text.
//   - metrics: computed widths for label/value/bar.
//   - opt: styling options.
//
// Returns:
//   - combined metadata string or empty when none.
func defaultMetaBuilder(r chartRow, metrics barMetrics, opt barRenderOptions) string {
	metaParts := []string{}
	if r.reset != "" {
		resetStyle := pickStyle(opt.resetStyle, resetBaseStyle)
		metaParts = append(metaParts, resetStyle.Render(r.reset))
	}
	if r.remain != "" {
		remainStyle := pickStyle(opt.remainStyle, remainBaseStyle)
		metaParts = append(metaParts, remainStyle.Render(r.remain))
	}
	return strings.Join(metaParts, consts.TextSeparatorDot)
}

// renderSkeleton builds placeholder chart rows while data is loading.
//
// Parameters:
//   - contentWidth: available width for the skeleton chart.
//
// Returns:
//   - rendered skeleton chart string.
func renderSkeleton(contentWidth int) string {
	const chartFrame = 6
	chartWidth := utils.Max(12, contentWidth-2)
	innerWidth := utils.Max(4, chartWidth-chartFrame)

	placeholderRows := []chartRow{
		{label: consts.LabelCurrent, percent: 65, reset: consts.TextSkeletonReset, remain: consts.TextSkeletonLeft},
		{label: consts.LabelWeekly, percent: 40, reset: consts.TextSkeletonReset, remain: consts.TextSkeletonLeft},
	}

	chart := renderBarsWithOptions(placeholderRows, innerWidth, barRenderOptions{
		labelStyle:     skeletonLabelStyle,
		valueStyle:     valueBaseStyle.Foreground(paletteMuted),
		resetStyle:     ptrStyle(resetBaseStyle.Foreground(paletteMuted)),
		remainStyle:    ptrStyle(remainBaseStyle.Foreground(paletteMuted)),
		barFillStyle:   skeletonBarFillStyle,
		barEmptyStyle:  skeletonBarEmptyStyle,
		valueFormatter: func(float64) string { return "···" },
		metaBuilder: func(r chartRow, metrics barMetrics, opt barRenderOptions) string {
			return renderSkeletonMeta(r, metrics, opt)
		},
	})

	return chartBoxStyle.Width(chartWidth).Render(chart)
}

// renderSkeletonMeta renders placeholder meta rows with muted styling.
//
// Parameters:
//   - r: placeholder row values.
//   - metrics: bar sizing info.
//   - opt: styling hooks.
//
// Returns:
//   - rendered meta string for the skeleton rows.
func renderSkeletonMeta(r chartRow, metrics barMetrics, opt barRenderOptions) string {
	if metrics.barWidth <= 0 {
		return ""
	}
	reset := truncateWidth(r.reset, metrics.barWidth)
	remain := truncateWidth(r.remain, metrics.barWidth)
	parts := []string{}
	if reset != "" {
		resetStyle := pickStyle(opt.resetStyle, skeletonMetaStyle)
		parts = append(parts, resetStyle.Render(reset))
	}
	if remain != "" {
		remainStyle := pickStyle(opt.remainStyle, skeletonMetaStyle)
		parts = append(parts, remainStyle.Render(remain))
	}
	return strings.Join(parts, consts.TextSeparatorDot)
}

// pickStyle returns the provided style or a fallback when nil.
//
// Parameters:
//   - s: optional style pointer.
//   - fallback: style to return when s is nil.
//
// Returns:
//   - resolved style.
func pickStyle(s *lipgloss.Style, fallback lipgloss.Style) lipgloss.Style {
	if s == nil {
		return fallback
	}
	return *s
}

// ptrStyle returns a pointer to the provided style.
//
// Parameters:
//   - s: style value.
//
// Returns:
//   - pointer to the provided style.
func ptrStyle(s lipgloss.Style) *lipgloss.Style {
	return &s
}

// longestLabel returns the widest label width, honoring a minimum.
//
// Parameters:
//
//	rows  - chart rows to measure.
//	floor - minimum width to enforce.
//
// Returns:
//
//	int - maximum label width.
func longestLabel(rows []chartRow, floor int) int {
	maxLen := floor
	for _, r := range rows {
		if w := lipgloss.Width(r.label); w > maxLen {
			maxLen = w
		}
	}
	return maxLen
}

// renderProgressBar draws a filled/empty bar representing pct across width.
//
// Parameters:
//
//	width - total cells available.
//	pct   - fractional progress in [0,1].
//
// Returns:
//
//	string - rendered bar with filled and empty segments.
func renderProgressBar(width int, pct float64) string {
	return renderProgressBarStyled(width, pct, barFillStyle, barEmptyStyle)
}

// renderProgressBarStyled draws a filled/empty bar with custom styles.
//
// Parameters:
//   - width: total cells for the bar.
//   - pct: progress fraction [0,1].
//   - fillStyle: style for the filled portion.
//   - emptyStyle: style for the empty portion.
//
// Returns:
//   - rendered bar string.
func renderProgressBarStyled(width int, pct float64, fillStyle, emptyStyle lipgloss.Style) string {
	if width <= 0 {
		return ""
	}
	fill := int(math.Round(pct * float64(width)))
	if fill < 0 {
		fill = 0
	} else if fill > width {
		fill = width
	}

	filled := fillStyle.Width(fill).Render(strings.Repeat(" ", fill))
	empty := emptyStyle.Width(width - fill).Render(strings.Repeat(" ", width-fill))

	return filled + empty
}

// computeBarMetrics derives widths for label, bar, and value columns.
//
// Parameters:
//
//	totalWidth - available width for the bar row.
//	rows       - data rows to measure label width.
//	valueWidth - width reserved for the percentage column.
//
// Returns:
//
//	barMetrics - computed widths for label, value, and bar sections.
func computeBarMetrics(totalWidth int, rows []chartRow, valueWidth int) barMetrics {
	const (
		minLabelWidth = 6
		minBarWidth   = 8
		// separator (1) + space (2) between label/bar/value
		layoutSpacing = 3
	)

	available := totalWidth - valueWidth - layoutSpacing
	if available < 1 {
		available = 1
	}

	desiredLabel := longestLabel(rows, minLabelWidth)
	maxLabel := available - minBarWidth
	if maxLabel < 0 {
		maxLabel = 0
	}
	labelWidth := desiredLabel
	if labelWidth > maxLabel {
		labelWidth = maxLabel
	}
	barWidth := available - labelWidth
	if barWidth < 1 {
		barWidth = 1
	}
	if labelWidth == 0 && available > 1 {
		labelWidth = 1
		barWidth = available - labelWidth
		if barWidth < 1 {
			barWidth = 1
		}
	}

	return barMetrics{
		labelWidth: labelWidth,
		valueWidth: valueWidth,
		barWidth:   barWidth,
	}
}

// truncateWidth shortens text to fit within maxWidth cells, respecting runes.
//
// Parameters:
//   - text: string to truncate.
//   - maxWidth: maximum width in cells.
//
// Returns:
//   - original text if it fits; otherwise a trimmed version.
func truncateWidth(text string, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}
	if lipgloss.Width(text) <= maxWidth {
		return text
	}
	runes := []rune(text)
	for lipgloss.Width(string(runes)) > maxWidth {
		runes = runes[:len(runes)-1]
	}
	return string(runes)
}
