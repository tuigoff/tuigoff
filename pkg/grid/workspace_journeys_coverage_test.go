package grid

import (
	"fmt"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
)

func TestWorkspaceSelectedRowLosesActiveEmphasisWhenFocusMovesAway(t *testing.T) {
	columns := []Column{{Name: "Name"}, {Name: "State"}}
	rows := []Row{{Key: "one", Values: []any{"Alpha", "ready"}}, {Key: "two", Values: []any{"Beta", "waiting"}}}
	m := New(columns, rows)
	m.SelectRow(0)
	active := m.View(48, true)
	if !containsPlain(active, "Alpha") || !containsPlain(active, "ready") {
		t.Fatalf("focused result no longer shows its selected row: %q", ansi.Strip(active))
	}
	m.SetFocused(false)
	inactive := m.View(48, false)
	if active == inactive {
		t.Fatal("selected row kept the same active styling after its grid lost focus")
	}
	activeLines := strings.Split(ansi.Strip(active), "\n")
	inactiveLines := strings.Split(ansi.Strip(inactive), "\n")
	if len(activeLines) != len(inactiveLines) || strings.Join(activeLines[1:len(activeLines)-1], "\n") != strings.Join(inactiveLines[1:len(inactiveLines)-1], "\n") {
		t.Fatalf("focus loss changed result rows:\nactive %q\ninactive %q", ansi.Strip(active), ansi.Strip(inactive))
	}
}

func TestWorkspaceSelectedRowUsesNeutralColorsAcrossFocusTransitions(t *testing.T) {
	columns := []Column{{Name: "Name"}, {Name: "State"}}
	rows := []Row{{Key: "one", Values: []any{"Alpha", "ready"}}, {Key: "two", Values: []any{"Beta", "waiting"}}}
	m := New(columns, rows)
	m.SetWorkspaceStyle(true)
	m.SelectRow(0)
	m.SetFocused(true)
	focusedView := m.View(48, true)
	focusedSelected := m.rowStyle(true)
	focusedUnselected := m.rowStyle(false)
	m.SetFocused(false)
	unfocusedView := m.View(48, false)
	unfocusedSelected := m.rowStyle(true)
	if !containsPlain(focusedView, "Alpha") || !containsPlain(unfocusedView, "Alpha") {
		t.Fatal("workspace focus transition hid the selected row")
	}
	if fmt.Sprint(focusedSelected.GetForeground()) != fmt.Sprint(unfocusedSelected.GetForeground()) ||
		fmt.Sprint(focusedSelected.GetBackground()) != fmt.Sprint(unfocusedSelected.GetBackground()) {
		t.Fatal("selected workspace row colors changed with keyboard focus")
	}
	if fmt.Sprint(focusedSelected.GetForeground()) == fmt.Sprint(focusedUnselected.GetForeground()) &&
		fmt.Sprint(focusedSelected.GetBackground()) == fmt.Sprint(focusedUnselected.GetBackground()) {
		t.Fatal("neutral selection styling does not distinguish the selected row")
	}
}

func TestWorkspaceLazySourceSwallowsPrintableInputWithoutFilteringOrFetching(t *testing.T) {
	source := &countingSource{n: 200}
	m := newSourceGrid(source)
	m.SetWorkspaceStyle(true)
	rowsLoadedAtStart := source.rowCalls
	keyHandlerCalls := 0
	m.SetKeyHandler(func(*Model, tea.KeyPressMsg) (tea.Cmd, bool) {
		keyHandlerCalls++
		return nil, true
	})
	selectedBefore := m.CurrentIndex()

	updated, cmd := m.Update(tea.KeyPressMsg{Text: "x", Code: 'x'})
	if cmd != nil || updated != m {
		t.Fatalf("printable source input returned model/cmd = %p/%v", updated, cmd)
	}
	if keyHandlerCalls != 0 {
		t.Fatalf("workspace printable input fell through to the product key handler %d times", keyHandlerCalls)
	}
	if got := source.rowCalls; got != rowsLoadedAtStart {
		t.Fatalf("typing fetched %d extra source rows (from %d to %d)", got-rowsLoadedAtStart, rowsLoadedAtStart, got)
	}
	if got := m.CurrentIndex(); got != selectedBefore {
		t.Fatalf("printable input moved lazy-source selection from %d to %d", selectedBefore, got)
	}
}

func containsPlain(view, text string) bool {
	return strings.Contains(ansi.Strip(view), text)
}
