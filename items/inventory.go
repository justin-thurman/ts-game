package items

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"sync"
)

// Inventory represents a player's inventory.
type Inventory struct {
	armor   []*armor
	weapons []*weapon
	// TODO: other types, consumes, etc.
	mu sync.Mutex
}

// NewInventory creates a new Inventory struct.
func NewInventory(weaponIds []int) *Inventory {
	inv := &Inventory{
		armor:   make([]*armor, 0),
		weapons: make([]*weapon, 0),
	}
	for _, id := range weaponIds {
		weap, err := FindWeaponById(id)
		if err != nil {
			slog.Error("Error building new inventory: " + err.Error())
			continue
		}
		inv.weapons = append(inv.weapons, weap)
	}
	return inv
}

// String returns a list of all the items in the player's inventory.
func (i *Inventory) String() string {
	var out strings.Builder
	out.WriteString("Inventory:")
	for _, a := range i.armor {
		_, err := out.WriteString("\n" + a.String())
		if err != nil {
			slog.Error("Failed building inventory string: " + err.Error())
		}
	}
	for _, w := range i.weapons {
		_, err := out.WriteString("\n" + w.String())
		if err != nil {
			slog.Error("Failed building inventory string: " + err.Error())
		}
	}
	return out.String()
}

// Wear attempts to wear an item.
func (i *Inventory) Wear(itemName string, einfo *EquipInfo) (message string) {
	i.mu.Lock()
	defer i.mu.Unlock()
	var outMessage strings.Builder
	item := i.findArmorByName(itemName)
	if item == nil {
		return "Wear what?"
	}
	switch item.Slot {
	case "body":
		if einfo.body != nil {
			i.addArmor(einfo.body)
			outMessage.WriteString(fmt.Sprintf("You remove %s.\n", einfo.body.String()))
		}
	case "legs":
		if einfo.legs != nil {
			i.addArmor(einfo.legs)
			outMessage.WriteString(fmt.Sprintf("You remove %s.\n", einfo.legs.String()))
		}
	case "helm":
		if einfo.helm != nil {
			i.addArmor(einfo.helm)
			outMessage.WriteString(fmt.Sprintf("You remove %s.\n", einfo.helm.String()))
		}
	}
	item.Equip(einfo)
	i.removeArmor(item)
	outMessage.WriteString(fmt.Sprintf("You wear %s.", item.String()))
	return outMessage.String()
}

// Wield attempts to wield a weapon
func (i *Inventory) Wield(itemName string, einfo *EquipInfo) (message string) {
	i.mu.Lock()
	defer i.mu.Unlock()
	var outMessage strings.Builder
	item := i.findWeaponByName(itemName)
	if item == nil {
		return "Wear what?"
	}
	if einfo.mainWeapon != nil {
		i.addWeapon(einfo.mainWeapon)
		outMessage.WriteString(fmt.Sprintf("You stop wielding %s.\n", einfo.mainWeapon.String()))
	}
	item.Equip(einfo)
	i.removeWeapon(item)
	outMessage.WriteString(fmt.Sprintf("You wield %s.", item.String()))
	return outMessage.String()
}

func (i *Inventory) addWeapon(w *weapon) {
	i.weapons = append(i.weapons, w)
}

func (i *Inventory) removeWeapon(w *weapon) {
	for idx, weap := range i.weapons {
		if weap == w {
			i.weapons = slices.Delete(i.weapons, idx, idx+1)
		}
	}
}

func (i *Inventory) addArmor(a *armor) {
	i.armor = append(i.armor, a)
}

func (i *Inventory) removeArmor(a *armor) {
	for idx, arm := range i.armor {
		if arm == a {
			i.armor = slices.Delete(i.armor, idx, idx+1)
		}
	}
}

func (i *Inventory) findWeaponByName(itemName string) *weapon {
	for _, weap := range i.weapons {
		for _, targetingName := range weap.TargetingNames {
			if itemName == targetingName {
				return weap
			}
		}
	}
	return nil
}

func (i *Inventory) findArmorByName(itemName string) *armor {
	for _, arm := range i.armor {
		for _, targetingName := range arm.TargetingNames {
			if itemName == targetingName {
				return arm
			}
		}
	}
	return nil
}
