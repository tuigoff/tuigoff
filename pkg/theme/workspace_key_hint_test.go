package theme

import (
	"regexp"
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
)

// A key hint is dimmed by its colour only. SGR faint on top of it makes the key
// unreadable in terminals that honour faint, such as Windows Terminal.
func TestWorkspaceKeyHintHasNoFaintAttribute(t *testing.T) {
	sgr := regexp.MustCompile("\x1b\\[([0-9;]*)m")
	trueColor := regexp.MustCompile(`[34]8;2;\d+;\d+;\d+|[34]8;5;\d+`)
	for _, dark := range []bool{true, false} {
		SetDark(dark)
		hint := WorkspaceKeyHint("F3")
		if ansi.Strip(hint) != "F3" {
			t.Fatalf("hint text changed: %q", hint)
		}
		for _, match := range sgr.FindAllStringSubmatch(hint, -1) {
			for _, param := range strings.Split(trueColor.ReplaceAllString(match[1], ""), ";") {
				if param == "2" {
					t.Fatalf("key hint uses SGR faint: %q", hint)
				}
			}
		}
	}
}
