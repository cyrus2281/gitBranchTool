package internal

type Git struct {
	currentBranch  string
	isGitRepo      bool
	hasInitialized bool
}

// initialize checks if current directory is a git repository and sets the current branch
func (g *Git) initialize() {
	// TODO: Implement this function
	g.hasInitialized = true
	g.isGitRepo = true
	g.currentBranch = "main"
}

// IsGitRepo returns true if the current directory is a git repository
func (g *Git) IsGitRepo() bool {
	if !g.hasInitialized {
		g.initialize()
	}
	return g.isGitRepo
}

// GetCurrentBranch returns the current branch of the git repository
func (g *Git) GetCurrentBranch() string {
	if !g.hasInitialized {
		g.initialize()
	}
	return g.currentBranch
}

type GitInterface interface {
	IsGitRepo() bool
	GetCurrentBranch() string
}
