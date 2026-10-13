package theme

import (
	"github.com/charmbracelet/x/ansi"
	"strings"
	"testing"
)

func TestPanelPolishStructuredChromeAndStableSelection(t *testing.T) {
	prevDark := Dark
	t.Cleanup(func() { SetDark(prevDark) })
	for _, width := range []int{140, 100, 80, 22, 9, 1} {
		view := WorkspaceChromeBar(width, WorkspaceChrome{Title: "DataTug", Context: []WorkspaceContext{{Label: "Project", Value: strings.Repeat("long", 30), Key: "F3"}}, Actions: []WorkspaceChromeAction{{Key: "F8", Label: "Server", Status: "LIVE"}, {Key: "F1", Label: "Help"}}})
		if !strings.HasPrefix(view, ansi.ResetStyle) {
			t.Fatal("header inherits terminal underline")
		}
		if ansi.StringWidth(view) != width {
			t.Fatalf("width%d: %q", width, view)
		}
		if width >= 9 && !strings.Contains(ansi.Strip(view), "F1 Help") {
			t.Fatal("help lost")
		}
		if strings.Contains(ansi.Strip(view), "[") {
			t.Fatal("bracketed key")
		}
	}
	selected := WorkspaceSelectedRowPlain("  Orders", true, 20)
	plain := WorkspaceSelectedRowPlain("  Orders", false, 20)
	if ansi.Strip(selected) != ansi.Strip(plain) || strings.Contains(ansi.Strip(selected), "›") {
		t.Fatal("selected entity moved")
	}
	if selected == plain {
		t.Fatal("selection styling lost")
	}
	for _, dark := range []bool{true, false} {
		SetDark(dark)
		if WorkspaceHintColor() == MutedColor() || WorkspaceHintColor() == TextColor() {
			t.Fatal("key hint color indistinguishable from label")
		}
		if WorkspaceChatTabStyle(true, true).GetForeground() != FocusColor() {
			t.Fatal("tab focus tint missing")
		}
		if WorkspaceChatTabStyle(true, false).GetBackground() != WorkspaceFrameColor() {
			t.Fatal("active tab does not match frame")
		}
		if WorkspaceChatTabStyle(false, false).GetBackground() == WorkspaceFrameColor() {
			t.Fatal("inactive tab not darker")
		}
	}
}
func TestPanelPolishFrameRetainsOnlyExistingMarkerColumn(t *testing.T) {
	for _, width := range []int{1, 2, 10, 70} {
		content := WorkspaceSurface("name", width, 4, true)
		view := WorkspaceChatFrame(content, width, 6, 3, true)
		lines := strings.Split(ansi.Strip(view), "\n")
		if len(lines) != 6 {
			t.Fatal("height")
		}
		for y, line := range lines {
			if ansi.StringWidth(line) != width {
				t.Fatalf("%d width %q", width, line)
			}
			if y < 5 && !strings.HasPrefix(line, "▎") {
				t.Fatal("marker did not reuse frame column")
			}
		}
	}
}

func TestPanelPolishFrameHasQuietCornersWithoutFocus(t *testing.T) {
	view := ansi.Strip(WorkspaceChatFrame("╭────────╮\n  name    ", 10, 4, 3, false))
	lines := strings.Split(view, "\n")
	if lines[0] != "╭────────╮" || !strings.HasPrefix(lines[1], "│ name") {
		t.Fatalf("corners/marker inset %q", view)
	}
	if strings.Contains(view, "▎") {
		t.Fatal("inactive frame looks focused")
	}
}
