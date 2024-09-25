package internal

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type RepositoryBranches struct {
	RepositoryName string
	StoreDirectory string
	branches       []Branch
	loaded         bool
}

type repositoryBranchesJson struct {
	Branches []Branch `json:"branches"`
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
		panic(err)
	}
	// Parse the JSON
	jsonData := repositoryBranchesJson{}
	json.Unmarshal(content, &jsonData)
	s.branches = jsonData.Branches
}

// Saves the persist file
func (s *RepositoryBranches) save() {
	repoStorePath := filepath.Join(s.StoreDirectory, s.RepositoryName+".json")
	jsonData := repositoryBranchesJson{s.branches}
	jsonDataBytes, err := json.Marshal(jsonData)
	if err != nil {
		// Error marshalling the JSON
		panic(err)
	}
	err = os.WriteFile(repoStorePath, jsonDataBytes, 0644)
	if err != nil {
		// Error writing the file
		panic(err)
	}
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
