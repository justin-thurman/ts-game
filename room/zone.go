package room

import (
	_ "gopkg.in/yaml.v3"
)

type Zone struct {
	Name  string  `yaml:"zone"`
	Rooms []*Room `yaml:"rooms"`
}
