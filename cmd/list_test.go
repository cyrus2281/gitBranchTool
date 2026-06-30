package cmd

import (
	"testing"

	"github.com/cyrus2281/gitBranchTool/internal"
)

func TestAppendUnregisteredBranches(t *testing.T) {
	repoDir := initTestRepo(t)
	chdir(t, repoDir)
	git := internal.Git{}

	// initTestRepo creates an initial commit on the default branch.
	defaultBranch, _ := git.GetCurrentBranch()

	git.CreateNewBranch("feature/registered")
	git.CreateNewBranch("feature/loose")

	registered := []internal.Branch{
		{Name: "feature/registered", Alias: "reg", Note: "note"},
	}

	result := appendUnregisteredBranches(&git, registered)

	// Registered branch must come first and keep its alias/note.
	if len(result) == 0 || result[0].Name != "feature/registered" || result[0].Alias != "reg" {
		t.Fatalf("expected registered branch first with alias, got %+v", result)
	}

	byName := map[string]internal.Branch{}
	for _, b := range result {
		byName[b.Name] = b
	}

	// The loose branch must appear with empty alias/note.
	loose, ok := byName["feature/loose"]
	if !ok {
		t.Fatal("expected unregistered branch feature/loose to be listed")
	}
	if loose.Alias != "" || loose.Note != "" {
		t.Errorf("unregistered branch should have empty alias/note, got %+v", loose)
	}

	// The default branch is also unregistered and must be listed.
	if _, ok := byName[defaultBranch]; !ok {
		t.Errorf("expected default branch %q to be listed", defaultBranch)
	}

	// The registered branch must not be duplicated.
	count := 0
	for _, b := range result {
		if b.Name == "feature/registered" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("registered branch should appear once, got %d", count)
	}
}
