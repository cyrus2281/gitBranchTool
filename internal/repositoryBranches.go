package internal

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/cyrus2281/go-logger"
)

const DEFAULT_BRANCH = "main"

type RepositoryBranches struct {
	RepositoryName string
	StoreDirectory string
	branches       []Branch
	defaultBranch  string
	loaded         bool
}

type repositoryBranchesJson struct {
	Branches      []Branch `json:"branches"`
	DefaultBranch string   `json:"defaultBranch"`
}

// Loads the persist file
func (s *RepositoryBranches) load() {
	s.loaded = true
	repoStorePath := filepath.Join(s.StoreDirectory, s.RepositoryName+".json")
	// Check if the file exists
	if _, err := os.Stat(repoStorePath); os.IsNotExist(err) {
		// File does not exist
		s.branches = []Branch{}
		return
	}
	// Read the file
	content, err := os.ReadFile(repoStorePath)
	if err != nil {
		// Error reading the file
		logger.Fatalln(err)
	}
	// Parse the JSON
	jsonData := repositoryBranchesJson{}
	json.Unmarshal(content, &jsonData)
	s.branches = jsonData.Branches
	s.defaultBranch = jsonData.DefaultBranch
}

// Saves the persist file
func (s *RepositoryBranches) save() {
	repoStorePath := filepath.Join(s.StoreDirectory, s.RepositoryName+".json")
	branchName := s.defaultBranch
	if branchName == "" {
		branchName = DEFAULT_BRANCH
	}
	jsonData := repositoryBranchesJson{s.branches, branchName}
	jsonDataBytes, err := json.Marshal(jsonData)
	logger.CheckFatalln(err)
	err = os.WriteFile(repoStorePath, jsonDataBytes, 0644)
	logger.CheckFatalln(err)
}

func (s *RepositoryBranches) SetDefaultBranch(branchName string) {
	if !s.loaded {
		s.load()
	}
	s.defaultBranch = branchName
	s.save()
}

func (s *RepositoryBranches) GetDefaultBranch() string {
	if !s.loaded {
		s.load()
	}
	if s.defaultBranch == "" {
		s.defaultBranch = DEFAULT_BRANCH
		s.save()
	}
	return s.defaultBranch
}

func (s *RepositoryBranches) GetBranches() []Branch {
	if !s.loaded {
		s.load()
	}
	return s.branches
}

func (s *RepositoryBranches) AddBranch(branch Branch) {
	if !s.loaded {
		s.load()
	}
	s.branches = append(s.branches, branch)
	s.save()
}

func (s *RepositoryBranches) BranchExists(branch Branch) bool {
	if !s.loaded {
		s.load()
	}
	for _, b := range s.branches {
		if b.Name == branch.Name {
			return true
		}
	}
	return false
}

func (s *RepositoryBranches) AliasExists(alias string) bool {
	if !s.loaded {
		s.load()
	}
	for _, b := range s.branches {
		if b.Alias == alias {
			return true
		}
	}
	return false
}

func (s *RepositoryBranches) BranchWithAliasExists(alias string) bool {
	if !s.loaded {
		s.load()
	}
	for _, b := range s.branches {
		if b.Name == alias {
			return true
		}
	}
	return false
}

func (s *RepositoryBranches) GetBranchByAlias(alias string) (Branch, bool) {
	if !s.loaded {
		s.load()
	}
	for _, b := range s.branches {
		if b.Alias == alias {
			return b, true
		}
	}
	return Branch{}, false
}

func (s *RepositoryBranches) GetBranchByName(name string) (Branch, bool) {
	if !s.loaded {
		s.load()
	}
	for _, b := range s.branches {
		if b.Name == name {
			return b, true
		}
	}
	return Branch{}, false
}

func (s *RepositoryBranches) GetBranchByNameOrAlias(name string) (Branch, bool) {
	if !s.loaded {
		s.load()
	}
	for _, b := range s.branches {
		if b.Name == name || b.Alias == name {
			return b, true
		}
	}
	return Branch{}, false
}

func (s *RepositoryBranches) RemoveBranch(branch Branch) {
	if !s.loaded {
		s.load()
	}
	for index, b := range s.branches {
		if b.Name == branch.Name {
			s.branches = append(s.branches[:index], s.branches[index+1:]...)
			s.save()
			return
		}
	}
}

func (s *RepositoryBranches) UpdateBranch(branch Branch) {
	if !s.loaded {
		s.load()
	}
	for index, b := range s.branches {
		if b.Name == branch.Name {
			s.branches[index] = branch
			s.save()
			return
		}
	}
}
