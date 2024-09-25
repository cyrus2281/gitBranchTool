package internal

type Branch struct {
	name  string
	alias string
	note  string
}

// Convert to string
// Delimiter optional, default is "|"
func (b *Branch) String(delimiter ...string) string {
	d := "|"
	if len(delimiter) > 0 {
		d = delimiter[0]
	}
	return b.name + d + b.alias + d + b.note
}
