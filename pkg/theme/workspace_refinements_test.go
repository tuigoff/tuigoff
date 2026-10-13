package theme

import (
	"github.com/charmbracelet/x/ansi"
	"strings"
	"testing"
)

func TestWorkspaceChromePreservesNavigationAndSelectionGrammar(t *testing.T) {
	for _, dark := range []bool{false, true} {
		SetDark(dark)
		for _, width := range []int{80, 100, 140} {
			view := WorkspaceTopBar(width, "DataTug", "Project: "+strings.Repeat("long", 60)+" [F3] ¦ Chat", []MenuItem{{Label: "Server OFF [F8]"}, {Label: "Help [F1]"}})
			plain := ansi.Strip(view)
			if ansi.StringWidth(plain) != width || !strings.HasPrefix(plain, "DataTug ¦ Project:") || !strings.HasSuffix(plain, "Help [F1]") || !strings.Contains(plain, "Server OFF [F8]") {
				t.Fatalf("navigation lost at %d: %q", width, plain)
			}
			if strings.Contains(plain, "│") {
				t.Fatal("full-height separator returned")
			}
		}
		bg, fg := WorkspaceTopColors()
		if Contrast(bg, fg) < 4.5 {
			t.Fatal("top text contrast")
		}
		bg, fg = WorkspaceHeaderColors()
		if Contrast(bg, fg) < 4.5 {
			t.Fatal("header contrast")
		}
		for _, selected := range []bool{true, false} {
			row := WorkspaceSelectedRow("Customer", selected, 24)
			if ansi.StringWidth(ansi.Strip(row)) != 24 || strings.HasPrefix(ansi.Strip(row), "›") != selected {
				t.Fatalf("selection grammar: %q", row)
			}
		}
		if ansi.Strip(WorkspacePanelHeader("Inspect", "[F7]")) != "Inspect [F7]" {
			t.Fatal("header hint lost")
		}
		if ansi.Strip(WorkspacePanelHeader("Inspect", "")) != "Inspect" {
			t.Fatal("empty hint")
		}
	}
	SetDark(true)
	if WorkspaceTabs(nil, 0, 0, false) != "" {
		t.Fatal("zero width tabs")
	}
	for _, focused := range []bool{false, true} {
		strip := WorkspaceTabs([]string{"Table", "Charts", "Current row", "Query"}, 3, 14, focused)
		if ansi.StringWidth(ansi.Strip(strip)) != 14 || !strings.Contains(ansi.Strip(strip), "›Query") {
			t.Fatalf("active tab cropped: %q", strip)
		}
		if !WorkspaceTabStyle(true, focused).GetUnderline() {
			t.Fatal("active tab not underlined")
		}
		if WorkspaceTabStyle(false, focused).GetUnderline() {
			t.Fatal("inactive tab competes")
		}
	}
	WorkspaceTopBar(1, "", "", nil)
	WorkspaceTopBar(5, "Title", "", []MenuItem{{Label: "long"}})
	WorkspaceSelectedRow("long", false, 1)
	WorkspaceTabs([]string{"only"}, -1, 10, false)
}

func TestWorkspaceNarrowNavigationRetainsRightmostHelp(t *testing.T) {
	view := ansi.Strip(WorkspaceTopBar(22, "DataTug", "Project", []MenuItem{{Label: "Server OFF [F8]"}, {Label: "Help [F1]"}}))
	if !strings.HasSuffix(view, "Help [F1]") {
		t.Fatalf("permanent help was clipped: %q", view)
	}
}
