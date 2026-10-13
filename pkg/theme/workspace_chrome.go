package theme

import (
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"image/color"
	"strings"
)

// WorkspaceContext separates metadata labels, values and optional key hints.
type WorkspaceContext struct{ Label, Value, Key string }

// WorkspaceChromeAction is navigation with a distinct shortcut and optional state.
type WorkspaceChromeAction struct {
	Key, Label, Status string
	Active             bool
}

// WorkspaceChrome describes opt-in product chrome without interpreting label text.
type WorkspaceChrome struct {
	Title   string
	Context []WorkspaceContext
	Actions []WorkspaceChromeAction
}

// WorkspaceKeyHint keeps shortcuts visually separate from their labels.
func WorkspaceKeyHint(key string) string {
	return lipgloss.NewStyle().Foreground(WorkspaceHintColor()).Faint(true).Render(key)
}

// WorkspaceHintColor remains visibly dimmer than labels without relying on terminal faint support.
func WorkspaceHintColor() color.Color {
	return pick(lipgloss.Color("#7A8491"), lipgloss.Color("#707D91"))
}

// WorkspaceMetadata renders descriptive information, never as a shortcut.
func WorkspaceMetadata(item WorkspaceContext) string {
	text := ""
	if item.Label != "" {
		text = lipgloss.NewStyle().Foreground(MutedColor()).Render(item.Label+":") + " "
	}
	text += lipgloss.NewStyle().Foreground(TextColor()).Render(item.Value)
	if item.Key != "" {
		text += " " + WorkspaceKeyHint(item.Key)
	}
	return text
}

// WorkspaceChromeBar reserves the trailing actions, keeping help visible when narrow.
func WorkspaceChromeBar(width int, chrome WorkspaceChrome) string {
	width = max(1, width)
	items := make([]string, 0, len(chrome.Actions))
	for _, item := range chrome.Actions {
		label := lipgloss.NewStyle().Foreground(TextColor()).Bold(item.Active).Render(item.Label)
		if item.Key != "" {
			label = WorkspaceKeyHint(item.Key) + " " + label
		}
		if item.Status != "" {
			label += " " + lipgloss.NewStyle().Foreground(WorkspaceMintColor()).Render(item.Status)
		}
		items = append(items, label)
	}
	rhs := strings.Join(items, "  ")
	for len(items) > 1 && ansi.StringWidth(rhs) > max(0, width-2) {
		items = items[1:]
		rhs = strings.Join(items, "  ")
	}
	rhs = ansi.Truncate(rhs, max(0, width-1), "")
	lhs := " " + lipgloss.NewStyle().Bold(true).Foreground(AccentColor()).Render(chrome.Title)
	for _, item := range chrome.Context {
		lhs += " " + WorkspaceKeyHint("¦") + " " + WorkspaceMetadata(item)
	}
	lhs = ansi.Truncate(lhs, max(0, width-ansi.StringWidth(rhs)-1), "…")
	bg, fg := WorkspaceTopColors()
	return ansi.ResetStyle + paintOver(lhs+strings.Repeat(" ", max(0, width-ansi.StringWidth(lhs)-ansi.StringWidth(rhs)-1))+rhs+" ", bg, fg)
}

// WorkspaceMintColor gives entity names and selection a consistent mint tint.
func WorkspaceMintColor() color.Color {
	return pick(lipgloss.Color("#176548"), lipgloss.Color("#91DFBD"))
}

// WorkspaceSelectedRowPlain keeps selected and unselected entity text at the same column.
func WorkspaceSelectedRowPlain(text string, selected bool, width int) string {
	width = max(1, width)
	style := lipgloss.NewStyle()
	if selected {
		bg, _ := NeutralSurfaceColors(true)
		style = style.Bold(true).Background(bg).Foreground(WorkspaceMintColor())
	}
	return style.Width(width).Render(ansi.Truncate(text, width, "…"))
}

// WorkspaceChatTabStyle joins active tabs to the neutral chat frame.
func WorkspaceChatTabStyle(active, focused bool) lipgloss.Style {
	bg, fg := WorkspaceTopColors()
	if active {
		bg, fg = WorkspaceFrameColor(), TextColor()
	}
	style := lipgloss.NewStyle().Background(bg).Foreground(fg).Bold(active)
	if active && focused {
		style = style.Foreground(FocusColor())
	}
	return style
}

// WorkspaceFrameColor joins the active tab background to the frame strokes.
func WorkspaceFrameColor() color.Color {
	return pick(lipgloss.Color("#A0B1C4"), lipgloss.Color("#4A596F"))
}

// WorkspaceFrameRule is the quiet one-cell frame grammar.
func WorkspaceFrameRule(text string) string {
	return lipgloss.NewStyle().Foreground(WorkspaceFrameColor()).Render(text)
}

// WorkspaceChatFrame reuses column zero for card and composer focus markers.
// content already includes its own marker/spacer columns. No nested padding is added.
func WorkspaceChatFrame(content string, width, height, composerTop int, composerFocused bool) string {
	width, height = max(1, width), max(1, height)
	lines := strings.Split(content, "\n")
	for len(lines) < height {
		lines = append(lines, "")
	}
	lines = lines[:height]
	for y, line := range lines {
		line = ansi.Truncate(line, width, "")
		first := ansi.Strip(ansi.Cut(line, 0, 1))
		edge := WorkspaceFrameRule("│")
		if first == "╭" {
			edge = WorkspaceFrameRule("╭")
		}
		if first == "▌" || first == "▎" || (composerFocused && y >= composerTop && y < height-1) {
			edge = lipgloss.NewStyle().Foreground(FocusColor()).Render("▎")
		}
		middle := ansi.Cut(line, 1, max(1, width-1))
		middle += strings.Repeat(" ", max(0, width-2-ansi.StringWidth(middle)))
		if width == 1 {
			lines[y] = edge
			continue
		}
		bg, fg := NeutralSurfaceColors(false)
		middle = paintOver(middle, bg, fg)
		right := WorkspaceFrameRule("│")
		if ansi.Strip(ansi.Cut(line, width-1, width)) == "╮" {
			right = WorkspaceFrameRule("╮")
		}
		lines[y] = edge + middle + right
	}
	if height > 1 {
		lines[height-1] = WorkspaceFrameRule("╰" + strings.Repeat("─", max(0, width-2)) + "╯")
		if width == 1 {
			lines[height-1] = WorkspaceFrameRule("╰")
		}
	}
	return strings.Join(lines, "\n")
}
