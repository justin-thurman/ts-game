package items

import (
	"cmp"
	"embed"
	"errors"
	"fmt"
	"log/slog"
	"slices"

	"gopkg.in/yaml.v3"
)

//go:embed itemdata
var itemdata embed.FS

var Weapons []*weapon

func Load() error {
	data, err := itemdata.ReadDir("itemdata")
	if err != nil {
		return err
	}

	for _, dirEntry := range data {
		if dirEntry.IsDir() {
			return errors.New("item loader does not (yet) support nested directories inside itemdata")
		}
		fileInfo, err := dirEntry.Info()
		if err != nil {
			return err
		}
		data, err := itemdata.ReadFile("itemdata/" + fileInfo.Name())
		if err != nil {
			return err
		}
		if err := yaml.Unmarshal(data, &Weapons); err != nil {
			return err
		}
	}
	slices.SortFunc(Weapons, func(a, b *weapon) int {
		return cmp.Compare(a.Id, b.Id)
	})
	return nil
}

func FindWeaponById(id int) (*weapon, error) {
	expectedWeaponIdx := id - 1 // Weapon IDs start at 1
	if expectedWeaponIdx >= len(Weapons) {
		return nil, fmt.Errorf("weapon Id %d outside bounds of Weapons array", id)
	}
	targetWeapon := Weapons[expectedWeaponIdx]
	if targetWeapon.Id == id {
		return targetWeapon, nil
	}
	slog.Error("Weapon Idx in Weapons array does not match Weapon id", "searchedForId", id, "foundWeaponId", targetWeapon.Id)
	// Fall back to a binary search
	target := &weapon{Id: id}
	idx, found := slices.BinarySearchFunc(Weapons, target, func(a, b *weapon) int {
		return cmp.Compare(a.Id, b.Id)
	})
	if !found {
		return nil, fmt.Errorf("weapon with Id %d not found", id)
	}
	return Weapons[idx], nil
}
