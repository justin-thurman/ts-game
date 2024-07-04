package room

import (
	"cmp"
	"embed"
	"errors"
	"fmt"
	"log/slog"
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
		for _, si := range z.SpecialMobs {
			si.Initialize()
		}
	}
	return nil
}

func FindRoomById(id int) (*Room, error) {
	expectedRoomIdx := id - 1 // Room IDs start at 1
	if expectedRoomIdx >= len(Rooms) {
		return nil, fmt.Errorf("room ID %d outside bounds of Rooms array", id)
	}
	targetRoom := Rooms[expectedRoomIdx]
	if targetRoom.Id == id {
		return targetRoom, nil
	}
	slog.Error("Room idx in Rooms array does not match room id", "searchedForId", id, "foundRoomId", targetRoom.Id)
	// Fall back to a binary search
	target := &Room{Id: id}
	idx, found := slices.BinarySearchFunc(Rooms, target, func(a, b *Room) int {
		return cmp.Compare(a.Id, b.Id)
	})
	if !found {
		return nil, fmt.Errorf("room with ID %d not found", id)
	}
	return Rooms[idx], nil
}
