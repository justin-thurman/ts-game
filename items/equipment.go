// Package items implements items and equipment.
package items

import (
	"fmt"
	"ts-game/dice"

	_ "gopkg.in/yaml.v3"
)

// EquipInfo represents a character's currently equipped gear.
type EquipInfo struct {
	body       *armor
	legs       *armor
	helm       *armor
	mainWeapon *weapon
}

// Damage returns the amount of damage dealt by a single swing with the character's current equipment.
func (einfo *EquipInfo) Damage() int {
	return einfo.mainWeapon.damage()
}

// String returns a display of the player's equipment slots and any equipped items.
func (einfo *EquipInfo) String() string {
	fmtString := `Equipment:
  Body: %s
  Legs: %s
  Helm: %s
  Main Weapon: %s`
	body := "Nothing"
	if einfo.body != nil {
		body = einfo.body.String()
	}
	legs := "Nothing"
	if einfo.legs != nil {
		legs = einfo.legs.String()
	}
	helm := "Nothing"
	if einfo.helm != nil {
		helm = einfo.helm.String()
	}
	mainWeapon := "Nothing"
	if einfo.mainWeapon != nil {
		mainWeapon = einfo.mainWeapon.String()
	}
	return fmt.Sprintf(fmtString, body, legs, helm, mainWeapon)
}

type armor struct {
	Name           string   `yaml:"name"`
	Slot           string   `yaml:"slot"` // TODO: probably make this a type
	TargetingNames []string `yaml:"targetingNames"`
	Id             int      `yaml:"id"`
}

// String returns the armor's name.
func (a *armor) String() string {
	return a.Name
}

// Equip equips the armor to the provided EquipInfo instance.
func (a *armor) Equip(equipInfo *EquipInfo) {
	// TODO: handle unequip, putting back in inventory
	switch a.Slot {
	case "body":
		equipInfo.body = a
	case "legs":
		equipInfo.legs = a
	case "helm":
		equipInfo.helm = a
	}
}

type weapon struct {
	Name           string    `yaml:"name"`
	TargetingNames []string  `yaml:"targetingNames"`
	Id             int       `yaml:"id"`
	DamageDice     dice.Dice `yaml:"damageDice"`
}

// String returns the weapon's name.
func (w *weapon) String() string {
	return w.Name
}

// Equip equips the weapon to the provided EquipInfo instance.
func (w *weapon) Equip(equipInfo *EquipInfo) {
	// TODO: handle unequip
	// TODO: handle two handed weapons, dual wield
	equipInfo.mainWeapon = w
}

func (w *weapon) damage() int {
	return w.DamageDice.Roll()
}
