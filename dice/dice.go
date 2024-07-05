// Package dice provides dice rolling capabilities.
package dice

import (
	"fmt"
	"math"
	"math/rand/v2"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// Dice is the representation a Number of N-Sided dice.
type Dice struct {
	Number int
	Sides  int
}

// Roll returns the value obtained by rolling the Dice.
func (d *Dice) Roll() int {
	var sum int
	for i := 1; i <= d.Number; i++ {
		sum += rand.IntN(d.Sides) + 1
	}
	return sum
}

// AverageN rolls the Dice N times and returns the average roll value, rounded.
func (d *Dice) AverageN(n int) int {
	if n == 0 {
		return 0
	}
	if n == 1 {
		return d.Roll()
	}
	var sum int
	for i := 0; i < n; i++ {
		sum += d.Roll()
	}
	return int(math.Round(float64(sum) / float64(n)))
}

// Max returns the maximum possible roll.
func (d *Dice) Max() int {
	return d.Number * d.Sides
}

func (d *Dice) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.ScalarNode {
		return fmt.Errorf("expected a scalar value")
	}

	parts := strings.Split(value.Value, "d")
	if len(parts) != 2 {
		return fmt.Errorf("invalid dice format; expected <number>d<sides>")
	}

	number, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("invalid number of dice: %v", err)
	}

	sides, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("invalid number of sides: %v", err)
	}

	d.Number = number
	d.Sides = sides

	return nil
}
