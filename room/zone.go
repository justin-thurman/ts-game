package room

import (
	"log/slog"
	"math"
	"math/rand"
	"sync"
	"ts-game/mob"

	_ "gopkg.in/yaml.v3"
)

type Zone struct {
	Name         string        `yaml:"zone"`
	Rooms        []*Room       `yaml:"rooms"`
	Spawns       []mob.MobInfo `yaml:"mobs"`
	MobCap       int           `yaml:"mobCap"`
	currMobCount int
	mobCountLock sync.RWMutex
}

func (z *Zone) Tick() {
	if z.Name == "anthill" {
		slog.Debug("Mobs", "mobCount", z.currMobCount)
	}
	if z.shouldSpawnMob() {
		z.spawnRandomMob()
	}
}

// Initialize zone by spawning random mobs up to the mob cap
func (z *Zone) initialize() {
	for z.currMobCount < z.MobCap {
		z.spawnRandomMob()
	}
}

// Implements an exponential spawning probability, where probability increases the further away the current
// count is from the max count
func (z *Zone) shouldSpawnMob() bool {
	z.mobCountLock.RLock()
	defer z.mobCountLock.RUnlock()
	if z.MobCap == 0 || len(z.Spawns) == 0 || len(z.Rooms) == 0 {
		return false
	}
	maxMobSpawnProbability := 0.25 // this results in approximately 16% spawn chance per tick if zone is empty
	spawnProbability := maxMobSpawnProbability * (1 - math.Exp(float64(-(float64(z.MobCap-z.currMobCount))/float64(z.MobCap))))
	slog.Debug("Spawn probability", "zone", z.Name, "probability", spawnProbability)
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
	z.mobCountLock.Lock()
	defer z.mobCountLock.Unlock()
	roomToSpawnIn.addMob(&mobInstance)
	z.currMobCount += 1
}

func (z *Zone) handleMobDeath() {
	z.mobCountLock.Lock()
	defer z.mobCountLock.Unlock()
	z.currMobCount -= 1
}
