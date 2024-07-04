package mob

import (
	"log/slog"
	"math/rand/v2"
	"sync"
	"time"

	_ "gopkg.in/yaml.v3"
)

type SpawnInfo struct {
	MobInfo     MobInfo `yaml:"mobInfo"`
	Timer       int     `yaml:"timer"` // in seconds
	RoomId      int     `yaml:"roomId"`
	SpawnChance float64 `yaml:"spawnChance"`
	canSpawn    bool
	mu          sync.Mutex
}

// Spawns a special mob from a SpawnInfo instance
func (si *SpawnInfo) Spawn() *Mob {
	si.mu.Lock()
	defer si.mu.Unlock()
	m := si.MobInfo.Spawn()
	m.spawnInfo = si
	si.canSpawn = false
	return m
}

// Whether this mob should spawn
func (si *SpawnInfo) ShouldSpawn() bool {
	si.mu.Lock()
	defer si.mu.Unlock()
	slog.Debug("Can spawn", "canSpawn", si.canSpawn)
	if !si.canSpawn {
		return false
	}
	return rand.Float64() < si.SpawnChance
}

// Handles the special mobs death, starting a timer for when it can respawn
func (si *SpawnInfo) HandleDeath() {
	si.mu.Lock()
	defer si.mu.Unlock()
	time.AfterFunc(time.Duration(si.Timer)*time.Second, si.reset)
}

func (si *SpawnInfo) Initialize() {
	si.canSpawn = true
}

func (si *SpawnInfo) reset() {
	si.mu.Lock()
	defer si.mu.Unlock()
	si.canSpawn = true
}
