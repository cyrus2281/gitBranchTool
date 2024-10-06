package internal

import (
	"github.com/spf13/cobra"
)

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func GetBranchesAndAliases() ([]string, cobra.ShellCompDirective) {
	completions := []string{}
	git := Git{}
	if git.IsGitRepo() {
		repoBranches := GetRepositoryBranches()
		for _, branch := range repoBranches.GetBranches() {
			completions = append(completions, branch.Name, branch.Alias)
		}
	}
	return completions, cobra.ShellCompDirectiveNoFileComp
}

func GetBranches() ([]string, cobra.ShellCompDirective) {
	completions := []string{}
	git := Git{}
	if git.IsGitRepo() {
		repoBranches := GetRepositoryBranches()
		for _, branch := range repoBranches.GetBranches() {
			completions = append(completions, branch.Name)
		}
	}
	return completions, cobra.ShellCompDirectiveNoFileComp
}

func GetAliases() ([]string, cobra.ShellCompDirective) {
	completions := []string{}
	git := Git{}
	if git.IsGitRepo() {
		repoBranches := GetRepositoryBranches()
		for _, branch := range repoBranches.GetBranches() {
			completions = append(completions, branch.Alias)
		}
	}
	return completions, cobra.ShellCompDirectiveNoFileComp
}

func GetGitBranches() ([]string, cobra.ShellCompDirective) {
	git := Git{}
	if git.IsGitRepo() {
		branches, err := git.GetBranches()
		if err != nil {
			return []string{}, cobra.ShellCompDirectiveError
		}
		return branches, cobra.ShellCompDirectiveNoFileComp
	}
	return []string{}, cobra.ShellCompDirectiveNoFileComp
}

func GetAllBranchesAndAliases() ([]string, cobra.ShellCompDirective) {
	completions, _ := GetGitBranches()
	git := Git{}
	if git.IsGitRepo() {
		repoBranches := GetRepositoryBranches()
		for _, branch := range repoBranches.GetBranches() {
			completions = append(completions, branch.Alias)
			// Add branch name if it's not already in the completions
			if !contains(completions, branch.Name) {
				completions = append(completions, branch.Name)
			}
		}
	}
	return completions, cobra.ShellCompDirectiveNoFileComp
}
