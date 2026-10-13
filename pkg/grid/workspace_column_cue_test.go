package grid

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/tuigoff/tuigoff/pkg/theme"
	"github.com/tuigoff/tuigoff/pkg/uitest"
)

func TestWorkspaceSelectedColumnOverridesProductGoldAndMovesWithArrows(t *testing.T) {
	m := New([]Column{{Name: "InvoiceId"}, {Name: "CustomerId"}}, []Row{{Key: "one", Values: []any{101, 29}}},
		WithTitle("Invoices · 1 loaded row"), WithCellStyle(func(Row, int, any) lipgloss.Style {
			return lipgloss.NewStyle().Bold(true).Foreground(theme.AccentColor())
		}))
	m.SetWorkspaceStyle(true)
	m.SetSize(70, 8)
	m.SetFocused(true)
	view := m.View(70, true)
	lines := strings.Split(view, "\n")
	if !strings.Contains(lines[1], "38;2;255;239;128") || !strings.Contains(lines[2], "38;2;255;239;128") || !strings.Contains(ansi.Strip(lines[1]), "InvoiceId") || !strings.Contains(ansi.Strip(lines[2]), "101") {
		t.Fatalf("selected PK foreground was hidden by product/header styles: %q", view)
	}
	if !strings.Contains(ansi.Strip(view), "Invoices · 1 loaded row") {
		t.Fatal("frameless result lost its identity")
	}
	m.Update(uitest.Key("right"))
	if m.SelectedColumn() != 1 || !strings.Contains(m.header(1), "38;2;255;239;128") || strings.Contains(m.header(0), "\x1b") {
		t.Fatal("selected foreground did not follow the column")
	}
	m.SetWorkspaceStyle(false)
	if strings.Contains(m.header(1), "\x1b") {
		t.Fatal("workspace header cue leaked to legacy rendering")
	}
}
