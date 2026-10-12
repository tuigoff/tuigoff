package grid

import (
	"reflect"
	"testing"

	"github.com/tuigoff/tuigoff/pkg/uitest"
)

func TestWorkspaceReorderRetainsSearchSelectionAndViewport(t *testing.T) {
	columns := []Column{{Name: "Id"}, {Name: "Long description"}, {Name: "Last column"}}
	rows := []Row{{Key: "a", Values: []any{2, "same", "two"}}, {Key: "b", Values: []any{1, "same", "one"}}}
	m := New(columns, rows)
	m.SetWorkspaceStyle(true)
	m.SetWidth(26)
	m.SetFocused(true)
	m.SelectColumn(2)
	m.SelectRow(1)
	for _, r := range "same" {
		m.Update(uitest.Text(string(r)))
	}
	if m.table.GetCurrentFilter() != "same" {
		t.Fatalf("search setup: %q", m.table.GetCurrentFilter())
	}
	m.SelectRow(1)
	offset := m.ColumnOffset()
	if err := m.ReorderRows([]string{"b", "a"}); err != nil {
		t.Fatal(err)
	}
	current, _ := m.CurrentRow()
	if current.Key != "b" || !m.Editing() || m.table.GetCurrentFilter() != "same" || m.ColumnOffset() != offset || m.Rows()[0].Key != "b" {
		t.Fatalf("local reorder: current=%q editing=%v filter=%q offset=%d/%d first=%q", current.Key, m.Editing(), m.table.GetCurrentFilter(), m.ColumnOffset(), offset, m.Rows()[0].Key)
	}
	before := append([]Row(nil), m.Rows()...)
	for _, bad := range [][]string{{"a"}, {"a", "a"}, {"a", "unknown"}} {
		if err := m.ReorderRows(bad); err == nil || !reflect.DeepEqual(before, m.Rows()) {
			t.Fatal("invalid permutation partially applied")
		}
	}
	m.Update(uitest.Key("esc"))
	if !m.Editing() || m.table.GetCurrentFilter() != "" {
		t.Fatal("first Esc did not clear retained search")
	}
	m.Update(uitest.Key("esc"))
	if m.Editing() {
		t.Fatal("second Esc did not leave search")
	}
	if err := m.ReorderRows([]string{"a", "b"}); err != nil {
		t.Fatal(err)
	}
	current, _ = m.CurrentRow()
	if current.Key != "b" {
		t.Fatal("reset lost source row")
	}
	duplicates := New(columns, []Row{{Key: "x"}, {Key: "x"}})
	if err := duplicates.ReorderRows([]string{"x", "x"}); err == nil {
		t.Fatal("ambiguous row identity accepted")
	}
}
