package combat

import (
	"slices"
	"sync"
	"ts-game/mob"
)

type PlayerCombatInfo struct {
	mobs []*mob.Mob
	sync.Mutex
}

func (pc *PlayerCombatInfo) EngagingMob(m *mob.Mob) {
	pc.Lock()
	defer pc.Unlock()
	pc.mobs = append(pc.mobs, m)
}

func (pc *PlayerCombatInfo) KilledMob(m *mob.Mob) {
	pc.Lock()
	defer pc.Unlock()
	for i, target := range pc.mobs {
		if target == m {
			pc.mobs = slices.Delete(pc.mobs, i, i+1)
		}
	}
}
