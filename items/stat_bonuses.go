package items

import (
	_ "gopkg.in/yaml.v3"
)

// StatsBonus represents an increase to stats on an item or items.
type StatsBonus struct {
	Str    int `yaml:"str"`
	Con    int `yaml:"con"`
	Damage int `yaml:"damage"`
}

// Add adds the values from one StatusBonus to another.
func (s *StatsBonus) Add(other *StatsBonus) {
	if other == nil {
		return
	}
	s.Str += other.Str
	s.Con += other.Con
	s.Damage += other.Damage
}
