package room

import (
	"math"
	"math/rand"
	"ts-game/mob"

	_ "gopkg.in/yaml.v3"
)

type Zone struct {
	Name         string        `yaml:"zone"`
	Rooms        []*Room       `yaml:"rooms"`
	Spawns       []mob.MobInfo `yaml:"mobs"`
	MobCap       int           `yaml:"mobCap"`
	currMobCount int
}

func (z *Zone) Tick() {
	if z.shouldSpawnMob() {
		z.spawnRandomMob()
	}
}

// Implements an exponential spawning probability, where probability increases the further away the current
// count is from the max count
func (z *Zone) shouldSpawnMob() bool {
	if z.MobCap == 0 || len(z.Spawns) == 0 || len(z.Rooms) == 0 {
		return false
	}
	maxMobSpawnProbability := 0.25 // this results in approximately 16% spawn chance per tick if zone is empty
	spawnProbability := maxMobSpawnProbability * (1 - math.Exp(float64(-(z.MobCap-z.currMobCount)/z.MobCap)))
	return rand.Float64() < spawnProbability
}

func (z *Zone) pickMobToSpawn() mob.MobInfo {
	return z.Spawns[rand.Intn(len(z.Spawns))]
}

func (z *Zone) pickRoomToSpawnIn() *Room {
	return z.Rooms[rand.Intn(len(z.Rooms))]
}

func (z *Zone) spawnRandomMob() {
	roomToSpawnIn := z.pickRoomToSpawnIn()
	mobToSpawn := z.pickMobToSpawn()
	mobInstance := mobToSpawn.Spawn()
	roomToSpawnIn.AddMob(&mobInstance)
}
