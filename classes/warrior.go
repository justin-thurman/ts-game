package classes

import (
	"log/slog"
	"ts-game/dice"
	"ts-game/items"
	"ts-game/stats"
)

// Warrior represents the Warrior class.
type Warrior struct{}

// StartingStats returns a Stats struct representing the warrior's starting stats.
func (w *Warrior) StartingStats() *stats.Stats {
	return stats.New(16, 16)
}

// StartingEquipment returns an EquipInfo representing the warrior's starting equipment.
func (w *Warrior) StartingEquipment() *items.EquipInfo {
	shortSword, err := items.FindWeaponById(1)
	if err != nil {
		slog.Error(err.Error())
	}
	einfo := items.EquipInfo{}
	shortSword.Equip(&einfo)
	return &einfo
}

// HitDice returns a Dice struct representing the warrior's hit dice.
func (w *Warrior) HitDice() *dice.Dice {
	return &dice.Dice{Number: 1, Sides: 10}
}

// String returns the class name.
func (w *Warrior) String() string {
	return "Warrior"
}
