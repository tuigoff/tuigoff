package theme

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
)

func TestWorkspaceNeutralCardKeepsHeaderAndWrappedBodyInsideSurface(t *testing.T) {
	withDark(t, false, func() {
		card := NeutralCard(RoleUser, "DETAILS", "alpha beta gamma delta", 12, true)
		plain := ansi.Strip(card)
		if !strings.Contains(plain, "DETAILS") || !strings.Contains(plain, "alpha") || !strings.Contains(plain, "delta") {
			t.Fatalf("card lost its header or wrapped content: %q", plain)
		}
		lines := strings.Split(card, "\n")
		if len(lines) < 4 {
			t.Fatalf("wrapped card has only %d lines: %q", len(lines), plain)
		}
		for i, line := range lines {
			if got := ansi.StringWidth(line); got > 12 {
				t.Fatalf("card line %d is %d cells wide, want <= 12: %q", i, got, ansi.Strip(line))
			}
		}
	})

	withoutHeader := NeutralCard(RoleSystem, "", "status", 20, false)
	if got := ansi.Strip(withoutHeader); !strings.Contains(got, "status") {
		t.Fatalf("headerless card lost its content: %q", got)
	}
}
