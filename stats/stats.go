// Package stats implements character and NPC statistics.
package stats

import (
	"fmt"
)

// Stats represents an entity's stats values.
type Stats struct {
	BaseStr     int
	Str         int
	StrModifier int
	BaseCon     int
	Con         int
	ConModifier int
	DamageRoll  int
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
  Constitution: %d (Base: %d) - Bonus to constitution rolls: %d
  Damage roll:  %d`
	return fmt.Sprintf(fmtString, s.Str, s.BaseStr, s.StrModifier, s.Con, s.BaseCon, s.ConModifier, s.DamageRoll)
}

func calculateModifier(statValue int) (modifier int) {
	return (statValue - 10) / 2
}

// AddStatsBonus adds StatsBonus structs to the entity's stats.
func (s *Stats) AddStatsBonus(b ...*StatsBonus) {
	for _, bonus := range b {
		if bonus == nil {
			continue
		}
		s.Str = s.BaseStr + bonus.Str
		s.Con = s.BaseCon + bonus.Con
		s.DamageRoll = s.DamageRoll + bonus.Damage
	}
	s.StrModifier = calculateModifier(s.Str)
	s.ConModifier = calculateModifier(s.Con)
}

// StatsBonus represents an increase to stats, i.e., from items or buffs.
type StatsBonus struct {
	Str    int `yaml:"str"`
	Con    int `yaml:"con"`
	Damage int `yaml:"damage"`
}

// Add adds the values from StatusBonus structs to another.
func (s *StatsBonus) Add(others ...*StatsBonus) {
	for _, other := range others {
		if other == nil {
			continue
		}
		s.Str += other.Str
		s.Con += other.Con
		s.Damage += other.Damage
	}
}
