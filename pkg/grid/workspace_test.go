package grid

import (
	tea "charm.land/bubbletea/v2"
	"github.com/tuigoff/tuigoff/pkg/uitest"
	"strings"
	"testing"
)

func TestWorkspaceSearchConsumesPrintableAndEscOneLayer(t *testing.T) {
	calls := 0
	m := New([]Column{{Name: "Name"}}, []Row{{Key: "one", Values: []any{"alpha"}}, {Key: "two", Values: []any{"beta"}}}, WithFilterDisabled(), WithKeyHandler(func(*Model, tea.KeyPressMsg) (tea.Cmd, bool) { calls++; return nil, true }))
	m.SetWorkspaceStyle(true)
	m.SetFocused(true)
	for _, text := range []string{"z", "r", "j", "e", "q", "s", "1", "ß"} {
		m.Update(uitest.Text(text))
		if !m.Editing() || m.table.GetCurrentFilter() != text || calls != 0 {
			t.Fatalf("printable %q escaped search: filter=%q calls=%d", text, m.table.GetCurrentFilter(), calls)
		}
		m.Update(uitest.Key("esc"))
		if !m.Editing() || m.table.GetCurrentFilter() != "" {
			t.Fatal("first Esc must clear without leaving search")
		}
		m.SetFocused(false)
		m.SetFocused(true)
		if !m.Editing() {
			t.Fatal("focus transfer lost empty active search")
		}
		m.Update(uitest.Key("esc"))
		if m.Editing() {
			t.Fatal("second Esc must leave search")
		}
	}
	m.Update(uitest.Key("/"))
	if !m.Editing() || m.table.GetCurrentFilter() != "" {
		t.Fatal("slash should open empty search")
	}
	m.Update(uitest.Key("esc"))
	m.Update(tea.KeyPressMsg{Code: 'B', Text: "B", Mod: tea.ModShift})
	if !m.Editing() || m.table.GetCurrentFilter() != "B" || calls != 0 {
		t.Fatal("shifted printable text invoked an action")
	}
	m.Update(uitest.Key("esc"))
	m.Update(uitest.Key("esc"))
	m.SetWorkspaceStyle(false)
	m.Update(uitest.Key("s"))
	if calls != 1 || m.Editing() {
		t.Fatal("turning workspace off did not restore legacy key/filter behavior")
	}
}

func TestWorkspaceViewLabelsDoNotAdvertiseSearchDigitsAsActions(t *testing.T) {
	m := New([]Column{{Name: "Name"}}, []Row{{Key: "one", Values: []any{"alpha"}}})
	for _, width := range []int{20, 30, 40, 50, 80} {
		legacy := m.viewLabels(width)
		m.SetWorkspaceStyle(true)
		for _, label := range m.viewLabels(width) {
			if strings.ContainsAny(label, "1234567890") {
				t.Fatalf("workspace advertises view number at width %d: %q", width, label)
			}
		}
		m.SetWorkspaceStyle(false)
		if strings.Join(m.viewLabels(width), "|") != strings.Join(legacy, "|") {
			t.Fatalf("legacy view shortcuts changed at width %d", width)
		}
	}
}

func TestWorkspaceSearchDoesNotChangeLegacyHandlerPrecedence(t *testing.T) {
	calls := 0
	m := New([]Column{{Name: "Name"}}, []Row{{Key: "one", Values: []any{"alpha"}}}, WithKeyHandler(func(*Model, tea.KeyPressMsg) (tea.Cmd, bool) { calls++; return nil, true }))
	m.SetFocused(true)
	m.Update(uitest.Key("s"))
	if calls != 1 || m.Editing() {
		t.Fatal("legacy product handler lost precedence")
	}
}

func TestWorkspaceStyleRestoresConfiguredOptions(t *testing.T) {
	for _, frameless := range []bool{false, true} {
		for _, filterOnType := range []bool{false, true} {
			opts := []Option{}
			if frameless {
				opts = append(opts, WithoutFrame())
			}
			if filterOnType {
				opts = append(opts, WithFilterOnType())
			}
			m := New([]Column{{Name: "Name"}}, []Row{{Key: "one", Values: []any{"alpha"}}}, opts...)
			m.SetWorkspaceStyle(false) // disabling an unused override is a no-op
			for i := 0; i < 2; i++ {
				m.SetWorkspaceStyle(true)
				m.SetWorkspaceStyle(true) // repeated enable cannot overwrite saved options
				m.SetWorkspaceStyle(false)
				m.SetWorkspaceStyle(false)
				if m.frameless != frameless || m.filterOnType != filterOnType || m.neutralSelection {
					t.Fatalf("configured frame=%v filter=%v changed to frame=%v filter=%v", frameless, filterOnType, m.frameless, m.filterOnType)
				}
			}
			m.SetFocused(true)
			m.Update(uitest.Text("a"))
			if m.Editing() != filterOnType {
				t.Fatalf("configured filter-on-type=%v was not restored", filterOnType)
			}
		}
	}
}
