package widgets_test

import (
	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/tuigoff/tuigoff/pkg/uitest"
	"github.com/tuigoff/tuigoff/pkg/widgets"
	"strings"
	"testing"
)

func TestChatFrameTabsKeepScrolledUnicodeCloseHitbox(t *testing.T) {
	tabs := widgets.NewTabs("sessions", widgets.TabsStyle{ChatFrame: true}, widgets.Tab{Title: "First"}, widgets.Tab{Title: "Current界", Closable: true})
	tabs.SetActive(1)
	tabs.SetSize(12, 2)
	tabs.Focus()
	lines := strings.Split(ansi.Strip(tabs.View()), "\n")
	if len(lines) != 2 || !strings.Contains(lines[0], "Current界") {
		t.Fatal("active session lost on narrow strip")
	}
	for _, line := range lines {
		if ansi.StringWidth(line) != 12 {
			t.Fatal("geometry")
		}
	}
	x := ansi.StringWidth(strings.Split(lines[0], "×")[0])
	_, cmd := tabs.Update(tea.MouseClickMsg{Button: tea.MouseLeft, X: x, Y: 0})
	msgs := uitest.Msgs(cmd)
	if len(msgs) != 1 {
		t.Fatal("close missed")
	}
	if event, ok := msgs[0].(widgets.TabCloseMsg); !ok || event.Index != 1 {
		t.Fatal("close bypasses product chooser")
	}
	_, cmd = tabs.Update(tea.MouseClickMsg{Button: tea.MouseLeft, X: x, Y: 1})
	if cmd != nil {
		t.Fatal("seam accepted click")
	}
}
