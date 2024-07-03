package room

import (
	"cmp"
	"embed"
	"errors"
	"fmt"
	"slices"

	"gopkg.in/yaml.v3"
)

//go:embed roomdata
var roomdata embed.FS

var (
	Rooms []*Room
	Zones []*Zone
)

func Load() error {
	data, err := roomdata.ReadDir("roomdata")
	if err != nil {
		return err
	}

	for _, dirEntry := range data {
		if dirEntry.IsDir() {
			return errors.New("room loader does not (yet) support nested directories inside roomdata")
		}
		fileInfo, err := dirEntry.Info()
		if err != nil {
			return err
		}
		data, err := roomdata.ReadFile("roomdata/" + fileInfo.Name())
		if err != nil {
			return err
		}

		var zone Zone
		if err := yaml.Unmarshal(data, &zone); err != nil {
			return err
		}
		Zones = append(Zones, &zone)
		Rooms = append(Rooms, zone.Rooms...)
		for _, r := range zone.Rooms {
			r.zone = &zone
		}
	}
	slices.SortFunc(Rooms, func(a, b *Room) int {
		return cmp.Compare(a.Id, b.Id)
	})
	slices.SortFunc(Zones, func(a, b *Zone) int {
		return cmp.Compare(a.Name, b.Name)
	})
	for _, r := range Rooms {
		r.initialize()
	}
	for _, z := range Zones {
		z.initialize()
	}
	return nil
}

func FindRoomById(id int) (*Room, error) {
	target := &Room{Id: id}
	idx, found := slices.BinarySearchFunc(Rooms, target, func(a, b *Room) int {
		return cmp.Compare(a.Id, b.Id)
	})
	if !found {
		return nil, fmt.Errorf("room with ID %d not found", id)
	}
	return Rooms[idx], nil
}
