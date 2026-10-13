package transcript

import (
	"strings"
	"testing"

	"charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/tuigoff/tuigoff/pkg/entity"
)

type workspaceJourneyFrameBlock struct{ views int }

func (b *workspaceJourneyFrameBlock) View(width int, focused bool) string {
	b.views++
	return "standalone result"
}
func (*workspaceJourneyFrameBlock) Update(tea.Msg) (Block, tea.Cmd) { return nil, nil }
func (*workspaceJourneyFrameBlock) Focusable() bool                 { return true }
func (*workspaceJourneyFrameBlock) SelfFramed() bool                { return true }

type workspaceJourneyEntityBlock struct{ ref *entity.Ref }

func (*workspaceJourneyEntityBlock) View(width int, focused bool) string { return "entity result" }
func (*workspaceJourneyEntityBlock) Update(tea.Msg) (Block, tea.Cmd)     { return nil, nil }
func (*workspaceJourneyEntityBlock) Focusable() bool                     { return true }
func (b *workspaceJourneyEntityBlock) Current() *entity.Ref              { return b.ref }

func TestWorkspaceTranscriptRendersCachedEntryKindsAndRetainsSelectedBlock(t *testing.T) {
	markdownRenders := 0
	m := New(WithMarkdownRenderer(func(text string, width int) string {
		markdownRenders++
		return "rendered: " + text
	}))
	m.SetSize(48, 24)
	frame := &workspaceJourneyFrameBlock{}
	ref := &entity.Ref{Type: "record", Keys: map[string]string{"id": "42"}, Title: "Selected record"}
	entityBlock := &workspaceJourneyEntityBlock{ref: ref}
	m.Append(Entry{Role: RoleUser, Text: "ordinary question"})
	m.Append(Entry{Role: RoleAssistant, Text: "**answer**", Markdown: true})
	m.Append(Entry{Role: RoleSystem, Text: "saved status"})
	m.Append(Entry{Role: RoleAssistant, Block: frame})
	m.Append(Entry{Role: RoleAssistant, Block: entityBlock})

	first := ansi.Strip(m.View())
	second := ansi.Strip(m.View())
	if !strings.Contains(first, "ordinary question") || !strings.Contains(first, "rendered: **answer**") ||
		!strings.Contains(first, "saved status") || !strings.Contains(first, "standalone result") || !strings.Contains(first, "entity result") {
		t.Fatalf("workspace transcript omitted an entry kind: %q", first)
	}
	if second != first || markdownRenders != 1 || frame.views != 1 {
		t.Fatalf("unchanged cached entries rendered again: markdown=%d frame=%d", markdownRenders, frame.views)
	}

	m.SetNeutralWorkspace(true)
	workspaceView := ansi.Strip(m.View())
	if markdownRenders != 2 || frame.views != 2 || !strings.Contains(workspaceView, "standalone result") {
		t.Fatalf("workspace opt-in did not invalidate and redraw cached entries: markdown=%d frame=%d", markdownRenders, frame.views)
	}
	m.Focus(2) // only user entries and interactive blocks create focus stops
	if got := m.BlockAt(1); got != frame {
		t.Fatalf("BlockAt(1) = %T, want the self-framed result", got)
	}
	if got := m.BlockAt(2); got != entityBlock {
		t.Fatalf("BlockAt(2) = %T, want the selected entity block", got)
	}
	if got := m.SelectedRef(2); got != ref {
		t.Fatalf("SelectedRef(2) = %v, want retained entity ref %v", got, ref)
	}
	if got := m.SelectedRef(1); got != nil {
		t.Fatalf("SelectedRef(1) = %v for a non-entity block, want nil", got)
	}
	if m.BlockAt(-1) != nil || m.BlockAt(99) != nil || m.SelectedRef(-1) != nil || m.SelectedRef(0) != nil {
		t.Fatal("out-of-range or non-entity stops returned a block reference")
	}
	m.RestoreScrollOffset(-10)
	if got := m.ScrollOffset(); got != 0 {
		t.Fatalf("negative restored history offset = %d, want clamped zero", got)
	}
}
