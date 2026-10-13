package grid

import (
	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/tuigoff/tuigoff/pkg/uitest"
	"strings"
	"testing"
)

func TestTabViewCyclingWrapsAndPreservesFilterRowColumn(t *testing.T) {
	columns, rows := sampleRows()
	m := New(columns, rows, WithExtraViews(CardView("Current row"), ExtraView{Label: "Query", Render: func(*Model, int, int) string { return "actual query" }}))
	m.SetFocused(true)
	m.SelectRow(1)
	m.SelectColumn(1)
	m.SetTabViewCycling(true)
	for _, key := range []string{"tab", "tab", "tab", "shift+tab"} {
		m.Update(uitest.Key(key))
	}
	if m.ActiveView() != 2 || m.CurrentIndex() != 1 || m.SelectedColumn() != 1 {
		t.Fatal("local tab cycle lost selection or wrap")
	}
	_, strip := m.switcher(80)
	if !strings.Contains(ansi.Strip(strip), "›3 Query") {
		t.Fatalf("selected Query shape missing: %q", strip)
	}
	m.SetTabViewCycling(false)
	m.ShowView(ViewTable)
	m.Update(uitest.Key("tab"))
	if m.ActiveView() != ViewTable {
		t.Fatal("legacy Tab changed views")
	}
	m.SetTabViewCycling(true)
	m.Update(uitest.Key("/"))
	m.Update(uitest.Key("tab"))
	if m.ActiveView() != ViewTable {
		t.Fatal("tab stole filter ownership")
	}
	m.Update(uitest.Key("esc"))
	m.SetKeyHandler(func(*Model, tea.KeyPressMsg) (tea.Cmd, bool) { return nil, true })
	m.Update(uitest.Key("tab"))
	if m.ActiveView() != ViewTable {
		t.Fatal("tab stole product ownership")
	}
	m.SetKeyHandler(nil)
	m.ShowView(2)
	m.extraViews[1].Update = func(*Model, tea.KeyPressMsg) (tea.Cmd, bool) { return nil, true }
	m.Update(uitest.Key("shift+tab"))
	if m.ActiveView() != 2 {
		t.Fatal("tab stole secondary editor")
	}
	empty := New(nil, nil)
	empty.SetTabViewCycling(true)
	empty.Update(uitest.Key("shift+tab"))
	if empty.ActiveView() != ViewTable {
		t.Fatal("single tab wrap")
	}
}
