package internal

type Store struct {
	RepositoryName string
	Branches       []Branch
}

// Loads the persist file
func (s *Store) Load() {
	// TODO: Implement this function
}

// Saves the persist file
func (s *Store) Save() {
	// TODO: Implement this function
}
