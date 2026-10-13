package theme

import (
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"image/color"
	"strings"
)

// NeutralSurfaceColors keeps focus on a neutral surface; blue is reserved for
// the single active left marker. Existing consumers retain their default theme.
func NeutralSurfaceColors(active bool) (bg, fg color.Color) {
	if active {
		return pick(lipgloss.Color("#D8DDE3"), lipgloss.Color("#30343B")), TextColor()
	}
	return pick(lipgloss.Color("#EEF0F3"), lipgloss.Color("#202329")), TextColor()
}

// WorkspaceSurface fills a borderless, fixed-size pane with one optional left marker.
func WorkspaceSurface(content string, width, height int, active bool) string {
	width, height = max(1, width), max(1, height)
	bg, fg := NeutralSurfaceColors(active)
	lines := strings.Split(content, "\n")
	for len(lines) < height {
		lines = append(lines, "")
	}
	lines = lines[:height]
	for i, line := range lines {
		line = ansi.Truncate(line, max(0, width-2), "")
		line = paintOver(line, bg, fg)
		marker := " "
		if active {
			marker = lipgloss.NewStyle().Foreground(FocusColor()).Render("▌")
		}
		lines[i] = ansi.Truncate(marker+lipgloss.NewStyle().Width(max(0, width-1)).Background(bg).Foreground(fg).Render(" "+line), width, "")
	}
	return strings.Join(lines, "\n")
}

// NeutralCard is an opt-in transcript card with neutral focus contrast.
func NeutralCard(role Role, header, body string, width int, focused bool) string {
	body = ansi.Wrap(body, max(1, width-3), " ")
	if header != "" {
		body = header + "\n" + body
	}
	return WorkspaceSurface(body, width, len(strings.Split(body, "\n"))+2, focused)
}

// WorkspacePane marks an adjacent-to-content divider: Explorer on its right,
// Inspect on its left. The whole focused pane remains a lighter neutral tint.
func WorkspacePane(title, content string, width, height int, focused, left bool) string {
	bg, fg := NeutralSurfaceColors(focused)
	width, height = max(1, width), max(1, height)
	lines := strings.Split(title+"\n"+content, "\n")
	for len(lines) < height {
		lines = append(lines, "")
	}
	lines = lines[:height]
	for i, line := range lines {
		line = ansi.Truncate(" "+line, max(0, width-1), "")
		line += strings.Repeat(" ", max(0, width-1-lipgloss.Width(line)))
		if i == 0 {
			headerBG, headerFG := WorkspaceHeaderColors()
			line = paintOver(line, headerBG, headerFG)
		} else {
			line = paintOver(line, bg, fg)
		}
		divider := " "
		if focused {
			divider = lipgloss.NewStyle().Foreground(FocusColor()).Render("│")
		}
		if left {
			line += divider
		} else {
			line = divider + line
		}
		lines[i] = line
	}
	return strings.Join(lines, "\n")
}

// WorkspaceTopColors gives opt-in navigation a distinct surface.
func WorkspaceTopColors() (bg, fg color.Color) {
	return pick(lipgloss.Color("#CFD8E5"), lipgloss.Color("#141B27")), TextColor()
}

// WorkspaceHeaderColors separates a panel's quiet header from its content.
func WorkspaceHeaderColors() (bg, fg color.Color) {
	return pick(lipgloss.Color("#E0E5EB"), lipgloss.Color("#252B35")), MutedColor()
}

// WorkspaceTopBar reserves right navigation before truncating left context.
// Right entries are product content, in their requested order; no product key
// or listener state is inferred by the theme.
func WorkspaceTopBar(width int, title, context string, right []MenuItem) string {
	width = max(1, width)
	var items []string
	for _, item := range right {
		items = append(items, WorkspaceTabStyle(item.Active, false).Underline(false).Render(item.Label))
	}
	rhs := strings.Join(items, "  ")
	// Later entries have priority, allowing permanent help to survive even
	// when a small terminal cannot accommodate every navigation label.
	for len(items) > 1 && ansi.StringWidth(rhs) > max(0, width-1) {
		items = items[1:]
		rhs = strings.Join(items, "  ")
	}
	rhs = ansi.Truncate(rhs, max(0, width-1), "")
	available := max(0, width-ansi.StringWidth(rhs)-1)
	lhs := lipgloss.NewStyle().Bold(true).Foreground(AccentColor()).Render(title)
	if context != "" {
		lhs += " ¦ " + context
	}
	lhs = ansi.Truncate(lhs, available, "…")
	bg, fg := WorkspaceTopColors()
	return paintOver(lhs+strings.Repeat(" ", max(0, width-ansi.StringWidth(lhs)-ansi.StringWidth(rhs)))+rhs, bg, fg)
}

// WorkspacePanelHeader keeps titles bold and hints quiet, independently.
func WorkspacePanelHeader(title, hint string) string {
	text := lipgloss.NewStyle().Bold(true).Foreground(MutedColor()).Render(title)
	if hint != "" {
		text += " " + lipgloss.NewStyle().Foreground(MutedColor()).Faint(true).Render(hint)
	}
	return text
}

// WorkspaceSelectedRow uses a shape and mint/green text independently of gold
// column selection and the blue focus divider. It retains selection on blur.
func WorkspaceSelectedRow(text string, selected bool, width int) string {
	width = max(1, width)
	prefix := "  "
	style := lipgloss.NewStyle()
	if selected {
		prefix = "› "
		bg, _ := NeutralSurfaceColors(true)
		style = style.Bold(true).Background(bg).Foreground(pick(lipgloss.Color("#176548"), lipgloss.Color("#91DFBD")))
	}
	return style.Width(width).Render(ansi.Truncate(prefix+text, width, "…"))
}

// WorkspaceTabStyle is the shared grammar for result, inspector and chat tabs.
func WorkspaceTabStyle(active, focused bool) lipgloss.Style {
	style := lipgloss.NewStyle().Foreground(MutedColor())
	if active {
		style = style.Bold(true).Underline(true).Foreground(TextColor())
		if focused {
			style = style.Foreground(FocusColor())
		}
	}
	return style
}

// WorkspaceTabs renders a bounded row; the active item has both underline and
// a shape fallback, and remains visible when the whole strip cannot fit.
func WorkspaceTabs(labels []string, active, width int, focused bool) string {
	if width <= 0 {
		return ""
	}
	var parts []string
	start, end, col := 0, 0, 0
	for i, label := range labels {
		text := " " + label + " "
		if i == active {
			text = "›" + label + " "
			start, end = col, col+ansi.StringWidth(text)
		}
		col += ansi.StringWidth(text)
		parts = append(parts, WorkspaceTabStyle(i == active, focused).Render(text))
	}
	scroll := min(max(0, end-width), start)
	text := ansi.Cut(strings.Join(parts, ""), scroll, scroll+width)
	return text + strings.Repeat(" ", max(0, width-ansi.StringWidth(text)))
}
