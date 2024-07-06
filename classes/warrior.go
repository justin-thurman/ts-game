package classes

import (
	"ts-game/dice"
	"ts-game/stats"
)

type Warrior struct{}

func (w *Warrior) StartingStats() *stats.Stats {
	return stats.New(16, 16)
}

func (w *Warrior) HitDice() *dice.Dice {
	return &dice.Dice{Number: 1, Sides: 10}
}
