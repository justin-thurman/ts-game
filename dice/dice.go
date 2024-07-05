package dice

import (
	"fmt"
	"math/rand/v2"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type Dice struct {
	Number int
	Sides  int
}

func (d *Dice) Roll() int {
	var sum int
	for i := 1; i <= d.Number; i++ {
		sum += rand.IntN(d.Sides) + 1
	}
	return sum
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
