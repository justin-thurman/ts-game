package room

import (
	"embed"
	"log"

	"gopkg.in/yaml.v3"
)

//go:embed roomdata
var roomdata embed.FS

var Rooms []*Room

func Load() error {
	log.Println(roomdata)
	data, err := roomdata.ReadFile("roomdata/newbietown.yaml")
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(data, Rooms); err != nil {
		return err
	}
	for _, r := range Rooms {
		log.Printf("Room name: %s", r.name)
	}
	return nil
}
