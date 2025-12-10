package app

import (
	"fmt"
	"strings"
	"time"

	"claude-monitor/internal/texts"
	"claude-monitor/internal/util"

	"github.com/charmbracelet/lipgloss"
)

const (
	minContainerWidth = 60
	horizontalPadding = 2
)

var (
	headerStatic   = renderHeader()
	helpStatic     = renderHelp()
	helpTextWidth  = lipgloss.Width(helpStatic)
	valueTextWidth = lipgloss.Width(valueBaseStyle.Render("100.0%"))
)

type layout struct {
	containerWidth int
	contentWidth   int
}

type barMetrics struct {
	labelWidth int
	valueWidth int
	barWidth   int
}

func newLayout(windowWidth int) layout {
	containerWidth := util.Max(windowWidth, minContainerWidth)
	contentWidth := util.Max(containerWidth-(horizontalPadding*2), minContainerWidth-(horizontalPadding*2))
	return layout{
		containerWidth: containerWidth,
		contentWidth:   contentWidth,
	}
}

func (m model) View() string {
	frame := newLayout(m.width)

	header := headerStatic
	body := renderBody(frame, m)
	footer := renderFooter(frame.contentWidth, helpStatic, helpTextWidth, renderStatus(m), m.cfg.RefreshEvery)

	content := lipgloss.JoinVertical(lipgloss.Left, header, body, footer)

	return pageStyle.
		Width(frame.containerWidth).
		Render(content)
}

func renderTitle() string {
	return headerStyle.
		Render(texts.HeaderTitle)
}

func renderMark() string {
	return markStyle.
		Render(texts.MarkArt)
}

func renderHeader() string {
	mark := renderMark()

	return lipgloss.JoinHorizontal(lipgloss.Top, mark, renderTitle())
}

func renderBody(frame layout, m model) string {
	sections := make([]string, 0, 3)

	switch {
	case len(m.rows) == 0 && m.err == nil:
		sections = append(sections, texts.TextNoData)
	default:
		const chartFrame = 2 /*borders*/ + 4              /*padding*/
		chartWidth := util.Max(14, frame.contentWidth-2)  // keep some inset to avoid wrapping
		innerWidth := util.Max(12, chartWidth-chartFrame) // inner area after border/padding
		sections = append(sections,
			chartBoxStyle.Width(chartWidth).Render(renderBars(m.rows, innerWidth)))
	}

	if m.err != nil {
		sections = append(sections, errorBox(fmt.Sprintf(texts.TextErrorFmt, m.err), frame.contentWidth))
	}

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func errorBox(msg string, width int) string {
	boxStyle := errorBoxStyleBase.Width(util.Max(40, width))

	content := errorContentStyleBase.
		Width(util.Max(38, width-4)).
		Render(msg)

	return boxStyle.Render(content)
}

func renderStatus(m model) string {
	switch {
	case m.loading:
		return fmt.Sprintf(texts.TextStatusFetch, spinnerStyle.Render(spinnerFrames[m.spinIndex]))
	case !m.lastUpdated.IsZero():
		return util.HumanTime(m.lastUpdated)
	default:
		return texts.TextStatusWaiting
	}
}

func renderHelp() string {
	shortcuts := []struct {
		keys []string
		desc string
	}{
		{keys: []string{"r"}, desc: "refresh now"},
		{keys: []string{"q", "ctrl+c"}, desc: "quit"},
	}

	parts := make([]string, 0, len(shortcuts))
	for _, sc := range shortcuts {
		key := helpKeyStyle.Render(strings.Join(sc.keys, "/"))
		desc := helpDescStyle.Render(sc.desc)
		parts = append(parts, lipgloss.JoinHorizontal(lipgloss.Top, key, " ", desc))
	}
	return helpStyle.Render(strings.Join(parts, texts.TextSeparatorDot))
}

func renderFooter(width int, helpText string, helpWidth int, status string, interval time.Duration) string {
	intervalText := fmt.Sprintf(texts.TextIntervalFmt, interval)
	rightContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		statusStyle.Render(status),
		separatorStyle.Render(texts.TextSeparatorDot),
		statusStyle.Render(intervalText),
	)
	rightW := lipgloss.Width(rightContent)
	lineWidth := util.Max(width, helpWidth+rightW)
	right := lipgloss.PlaceHorizontal(lineWidth-helpWidth, lipgloss.Right, rightContent)

	return footerStyle.
		Width(lineWidth).
		Render(lipgloss.JoinHorizontal(lipgloss.Top, helpText, right))
}

func renderBars(rows []chartRow, totalWidth int) string {
	if len(rows) == 0 {
		return ""
	}

	// Layout: [label] [bar filling width] [right-aligned percent]; meta lines are indented under the bar.
	metrics := computeBarMetrics(totalWidth, rows, valueTextWidth)

	blocks := make([]string, 0, len(rows))
	separator := barSeparatorStyle.Render(" ")
	labelStyle := labelBaseStyle.Width(metrics.labelWidth)
	valueStyle := valueBaseStyle.
		Width(metrics.valueWidth).
		AlignHorizontal(lipgloss.Right)
	space := "  "
	metaStyle := metaBaseStyle.MarginLeft(metrics.labelWidth + 1)

	for i, r := range rows {
		percent := util.Clamp(r.percent, 0, 100)
		bar := renderProgressBar(metrics.barWidth, percent/100)
		value := valueStyle.Render(fmt.Sprintf("%5.1f%%", percent))

		left := lipgloss.JoinHorizontal(lipgloss.Top,
			labelStyle.Render(r.label),
			separator,
			bar,
		)
		line := lipgloss.JoinHorizontal(lipgloss.Top, left+space, value)

		metaParts := []string{}
		if r.reset != "" {
			metaParts = append(metaParts, resetBaseStyle.Render(r.reset))
		}
		if r.remain != "" {
			metaParts = append(metaParts, remainBaseStyle.Render(r.remain))
		}

		var rendered string
		if len(metaParts) > 0 {
			// Second line: optional reset/remaining info, indented to align with bar start.
			metaLine := metaStyle.Render(strings.Join(metaParts, texts.TextSeparatorDot))
			rendered = lipgloss.JoinVertical(lipgloss.Left, line, metaLine)
		} else {
			rendered = line
		}

		if i == len(rows)-1 {
			blocks = append(blocks, barLastBlockStyle.Render(rendered))
		} else {
			blocks = append(blocks, barBlockStyle.Render(rendered))
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, blocks...)
}

func longestLabel(rows []chartRow, floor int) int {
	maxLen := floor
	for _, r := range rows {
		if w := lipgloss.Width(r.label); w > maxLen {
			maxLen = w
		}
	}
	return maxLen
}

// renderProgressBar draws a fixed-width bar without shared state or extra allocations.
func renderProgressBar(width int, pct float64) string {
	if width <= 0 {
		return ""
	}
	fill := int(pct * float64(width))
	if fill < 0 {
		fill = 0
	} else if fill > width {
		fill = width
	}

	// preallocate once, then slice without extra allocs
	bar := strings.Builder{}
	bar.Grow(width)
	for i := 0; i < fill; i++ {
		bar.WriteByte(' ')
	}
	filled := barFillStyle.Width(fill).Render(bar.String())

	bar.Reset()
	bar.Grow(width - fill)
	for i := 0; i < width-fill; i++ {
		bar.WriteByte(' ')
	}
	empty := barEmptyStyle.Width(width - fill).Render(bar.String())

	return filled + empty
}

// computeBarMetrics derives fixed widths for labels, values, and bars to keep rows aligned.
func computeBarMetrics(totalWidth int, rows []chartRow, valueWidth int) barMetrics {
	const minBarWidth = 12
	labelWidth := longestLabel(rows, 8)
	barWidth := util.Max(minBarWidth, totalWidth-labelWidth-valueWidth-3) // label + gap + bar + value
	return barMetrics{
		labelWidth: labelWidth,
		valueWidth: valueWidth,
		barWidth:   barWidth,
	}
}
