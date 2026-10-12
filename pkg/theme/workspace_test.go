package theme

import (
	"github.com/charmbracelet/x/ansi"
	"strings"
	"testing"
)

func TestWorkspacePanesHaveOnlyContentFacingDivider(t *testing.T) {
	for _, width := range []int{1, 12, 40} {
		for _, left := range []bool{false, true} {
			for _, active := range []bool{false, true} {
				view := WorkspacePane("Title", "content", width, 4, active, left)
				for _, line := range strings.Split(ansi.Strip(view), "\n") {
					if ansi.StringWidth(line) != width {
						t.Fatalf("width=%d: %q", width, line)
					}
					if active && ((left && !strings.HasSuffix(line, "│")) || (!left && !strings.HasPrefix(line, "│"))) {
						t.Fatal("divider is not content-facing")
					}
					if strings.ContainsAny(line, "┌┐└┘─▌") || (!active && strings.Contains(line, "│")) {
						t.Fatal("unexpected pane outline/marker")
					}
				}
				if active {
					inactiveBG, inactiveFG := NeutralSurfaceColors(false)
					if strings.Contains(view, fillSGR(inactiveBG, inactiveFG)) {
						t.Fatal("inactive surface overrides focused pane")
					}
				}
			}
		}
	}
}

func TestWorkspaceSurfaceMarkerAndDimensions(t *testing.T) {
	for _, width := range []int{1, 4, 40} {
		for _, active := range []bool{false, true} {
			view := ansi.Strip(WorkspaceSurface("hello\nworld", width, 3, active))
			for _, line := range strings.Split(view, "\n") {
				if ansi.StringWidth(line) != width || (active && !strings.HasPrefix(line, "▌")) || (!active && strings.HasPrefix(line, "▌")) {
					t.Fatalf("bad surface %q", line)
				}
			}
		}
	}
}
