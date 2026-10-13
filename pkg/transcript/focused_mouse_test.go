package transcript

import (
	tea "charm.land/bubbletea/v2"
	"testing"
)

type focusedMouseBlock struct {
	fakeBlock
	click tea.MouseClickMsg
}

func (b *focusedMouseBlock) Update(msg tea.Msg) (Block, tea.Cmd) {
	b.updates++
	b.click = msg.(tea.MouseClickMsg)
	return b, nil
}
func TestUpdateFocusedMouseDoesNotBroadcastOrReplaceRetainedEntry(t *testing.T) {
	m := New()
	m.SetSize(40, 10)
	first, last := &focusedMouseBlock{}, &focusedMouseBlock{}
	m.Append(Entry{ID: "first", Block: first})
	m.Append(Entry{ID: "last", Block: last})
	m.Focus(1)
	click := tea.MouseClickMsg{Button: tea.MouseLeft, X: 36, Y: 2}
	m.UpdateFocused(click)
	if first.updates != 0 || last.updates != 1 || last.click != click || m.BlockAt(1) != last {
		t.Fatal("mouse was broadcast or retained block lost")
	}
	m.Blur()
	m.UpdateFocused(click)
	if last.updates != 1 {
		t.Fatal("unfocused input delivered")
	}
}
