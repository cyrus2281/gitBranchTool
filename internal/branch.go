package internal

import "fmt"

type Branch struct {
	Name  string `json:"name"`
	Alias string `json:"alias"`
	Note  string `json:"note"`
}

// Convert to print format
func (b *Branch) String() string {
	return fmt.Sprintf("%-20s\t%-20s\t%-20s", b.Name, b.Alias, b.Note)
}
