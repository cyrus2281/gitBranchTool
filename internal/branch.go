package internal

import "fmt"

type Branch struct {
	Name  string `json:"name"`
	Alias string `json:"alias"`
	Note  string `json:"note"`
}

// Convert to string
// Delimiter optional, default is "|"
func (b *Branch) String(delimiter ...string) string {
	d := "|"
	if len(delimiter) > 0 {
		d = delimiter[0]
	}
	return b.Name + d + b.Alias + d + b.Note
}

// Convert to print format
func (b *Branch) Print() string {
	return fmt.Sprintf("%-20s\t%-20s\t%-20s", b.Name, b.Alias, b.Note)
}
