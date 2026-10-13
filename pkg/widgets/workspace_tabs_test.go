package widgets

import (
	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
	"strings"
	"testing"
)

func TestWorkspaceTabsScrollClickAndKeepActiveShape(t *testing.T) {
	tabs := NewTabs("chat", TabsStyle{Workspace: true}, Tab{ID: "1", Title: "First"}, Tab{ID: "2", Title: "Second"}, Tab{ID: "3", Title: "Current"})
	tabs.SetActive(2)
	tabs.SetSize(12, 1)
	tabs.Focus()
	view := ansi.Strip(tabs.View())
	if !strings.Contains(view, "›Current") || !tabs.tabStyle(2).GetUnderline() {
		t.Fatalf("active tab missing: %q", view)
	}
	if tabs.tabStyle(0).GetUnderline() {
		t.Fatal("inactive underline")
	}
	tabs.Blur()
	if !tabs.tabStyle(2).GetUnderline() {
		t.Fatal("blur lost selected tab")
	}
	spans, scroll := tabs.layout()
	tabs, cmd := tabs.Update(tea.MouseClickMsg{Button: tea.MouseLeft, X: spans[1].start - scroll, Y: 0})
	if cmd == nil || tabs.Active() != 1 {
		t.Fatal("rendered tab click did not select")
	}
}
