package cmd

import (
	"strings"
	"testing"
)

// cobraGeneratedCommands are commands Cobra adds automatically; they do not
// carry a manual annotation and are excluded from the guide.
var cobraGeneratedCommands = map[string]bool{
	"help":       true,
	"completion": true,
}

func TestBuildManualDefaultOnlyImportant(t *testing.T) {
	out := buildManual(false)

	for _, name := range importantCommands {
		if !strings.Contains(out, "## "+name) {
			t.Errorf("default manual is missing important command %q", name)
		}
	}

	// Non-important commands must not appear without --full.
	for _, name := range []string{"addAlias", "rename", "removeEntry", "upgrade"} {
		if strings.Contains(out, "## "+name) {
			t.Errorf("default manual should not include non-important command %q", name)
		}
	}
}

func TestBuildManualFullIncludesAllDocumented(t *testing.T) {
	out := buildManual(true)

	for _, c := range rootCmd.Commands() {
		name := c.Name()
		if cobraGeneratedCommands[name] {
			continue
		}
		// Hidden helpers (e.g. _ps) carry a manual but are not printed.
		if c.Hidden {
			if strings.Contains(out, "## "+name) {
				t.Errorf("full manual should not include hidden command %q", name)
			}
			continue
		}
		if !strings.Contains(out, "## "+name) {
			t.Errorf("full manual is missing command %q", name)
		}
	}
}

func TestEveryTopLevelCommandHasManual(t *testing.T) {
	for _, c := range rootCmd.Commands() {
		if cobraGeneratedCommands[c.Name()] {
			continue
		}
		if strings.TrimSpace(c.Annotations[manualAnnotation]) == "" {
			t.Errorf("command %q is missing a %q annotation", c.Name(), manualAnnotation)
		}
	}
}
