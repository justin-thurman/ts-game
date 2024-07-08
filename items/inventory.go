package items

import (
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

// New creates a new Inventory struct.
func New(itemIds []int) *Inventory {
	inv := &Inventory{
		armor:   make([]*armor, 0),
		weapons: make([]*weapon, 0),
	}
	// TODO: Get items by ID, then need to  figure out how to add them to the right slice; introspect type? ehhhhh
	return inv
}

// String returns a list of all the items in the player's inventory.
func (i *Inventory) String() string {
	var out strings.Builder
	out.WriteString("Inventory:")
	for _, a := range i.armor {
		_, err := out.WriteString(a.String())
		if err != nil {
			slog.Error("Failed building inventory string: " + err.Error())
		}
	}
	for _, w := range i.weapons {
		_, err := out.WriteString(w.String())
		if err != nil {
			slog.Error("Failed building inventory string: " + err.Error())
		}
	}
	return out.String()
}

// Wear attempts to wear an item.
func (i *Inventory) Wear(itemName string, einfo *EquipInfo) (message string) {
	var outMessage strings.Builder
	var item *armor
outerLoop:
	for _, each := range i.armor {
		for _, targetingName := range each.TargetingNames {
			if itemName == targetingName {
				item = each
				break outerLoop
			}
		}
	}
	if item == nil {
		return "Wear what?"
	}
	switch item.Slot {
	case "body":
		if einfo.body != nil {
			i.addArmor(einfo.body)
			outMessage.WriteString("You remove " + einfo.body.String())
		}
	case "legs":
		if einfo.legs != nil {
			i.addArmor(einfo.legs)
			outMessage.WriteString("You remove " + einfo.legs.String())
		}
	case "helm":
		if einfo.helm != nil {
			i.addArmor(einfo.helm)
			outMessage.WriteString("You remove " + einfo.helm.String())
		}
	}
	item.Equip(einfo)
	i.addArmor(item)
	outMessage.WriteString("You wear " + item.String())
	return outMessage.String()
}

func (i *Inventory) addWeapon(w *weapon) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.weapons = append(i.weapons, w)
}

func (i *Inventory) removeWeapon(w *weapon) {
	i.mu.Lock()
	defer i.mu.Unlock()
	for idx, weap := range i.weapons {
		if weap == w {
			i.weapons = slices.Delete(i.weapons, idx, idx+1)
		}
	}
}

func (i *Inventory) addArmor(a *armor) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.armor = append(i.armor, a)
}

func (i *Inventory) removeArmor(a *armor) {
	i.mu.Lock()
	defer i.mu.Unlock()
	for idx, arm := range i.armor {
		if arm == a {
			i.armor = slices.Delete(i.armor, idx, idx+1)
		}
	}
}
