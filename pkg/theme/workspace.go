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
	if focused {
		title = lipgloss.NewStyle().Bold(true).Foreground(FocusColor()).Render(title)
	}
	width, height = max(1, width), max(1, height)
	lines := strings.Split(title+"\n"+content, "\n")
	for len(lines) < height {
		lines = append(lines, "")
	}
	lines = lines[:height]
	for i, line := range lines {
		line = ansi.Truncate(" "+line, max(0, width-1), "")
		line += strings.Repeat(" ", max(0, width-1-lipgloss.Width(line)))
		line = paintOver(line, bg, fg)
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
