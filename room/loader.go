package room

import (
	"embed"
	"errors"

	"gopkg.in/yaml.v3"
)

//go:embed roomdata
var roomdata embed.FS

func Load() ([]Zone, error) {
	data, err := roomdata.ReadDir("roomdata")
	if err != nil {
		return nil, err
	}

	var zones []Zone

	for _, dirEntry := range data {
		if dirEntry.IsDir() {
			return nil, errors.New("room loader does not (yet) support nested directories inside roomdata")
		}
		fileInfo, err := dirEntry.Info()
		if err != nil {
			return nil, err
		}
		data, err := roomdata.ReadFile("roomdata/" + fileInfo.Name())
		if err != nil {
			return nil, err
		}

		var zone Zone
		if err := yaml.Unmarshal(data, &zone); err != nil {
			return nil, err
		}
		zones = append(zones, zone)
	}
	return zones, nil
}
