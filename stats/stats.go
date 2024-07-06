// Package stats implements character and NPC statistics.
package stats

import "fmt"

// Stats represents an entity's stats values.
type Stats struct {
	Str         int
	StrModifier int
	Con         int
	ConModifier int
}

// New creates a new Stats struct
func New(str, con int) *Stats {
	return &Stats{
		Str:         str,
		StrModifier: calculateModifier(str),
		Con:         con,
		ConModifier: calculateModifier(con),
	}
}

// String returns a string representing the stat values.
func (s *Stats) String() string {
	fmtString := `Stats [base (modifier)]:
  Strength:     %d (%d)
  Constitution: %d (%d)`
	return fmt.Sprintf(fmtString, s.Str, s.StrModifier, s.Con, s.ConModifier)
}

func calculateModifier(statValue int) (modifier int) {
	return (statValue - 10) / 2
}
