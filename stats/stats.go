// Package stats implements character and NPC statistics.
package stats

import (
	"fmt"
	"ts-game/items"
)

// Stats represents an entity's stats values.
type Stats struct {
	BaseStr     int
	Str         int
	StrModifier int
	BaseCon     int
	Con         int
	ConModifier int
}

// New creates a new Stats struct
func New(str, con int) *Stats {
	return &Stats{
		BaseStr:     str,
		Str:         str,
		StrModifier: calculateModifier(str),
		BaseCon:     con,
		Con:         con,
		ConModifier: calculateModifier(con),
	}
}

// String returns a string representing the stat values.
func (s *Stats) String() string {
	fmtString := `Stats [base (modifier)]:
  Strength:     %d (Base: %d) - Bonus to strength rolls: %d
  Constitution: %d (Base: %d) - Bonus to constitution rolls: %d`
	return fmt.Sprintf(fmtString, s.Str, s.BaseStr, s.StrModifier, s.Con, s.BaseCon, s.ConModifier)
}

func calculateModifier(statValue int) (modifier int) {
	return (statValue - 10) / 2
}

// AddStatsBonus adds a StatsBonus struct to the player's stats.
func (s *Stats) AddStatsBonus(b *items.StatsBonus) {
	s.Str = s.BaseStr + b.Str
	s.Con = s.BaseCon + b.Con
	s.StrModifier = calculateModifier(s.Str)
	s.ConModifier = calculateModifier(s.Con)
}
