package room

import (
	"embed"
	"errors"

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
	}
	for _, r := range Rooms {
		r.initialize()
	}
	return nil
}
