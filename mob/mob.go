package mob

import (
	"sync"
	"ts-game/dice"

	_ "gopkg.in/yaml.v3"
)

// The template that an individual mob is spawned from
type MobInfo struct {
	Name            string    `yaml:"name"`
	IdleDescription string    `yaml:"idleDescription"`
	TargetingNames  []string  `yaml:"targetingNames"`
	DamageDice      dice.Dice `yaml:"damageDice"`
	MaxHealth       int       `yaml:"health"`
	XpValue         int       `yaml:"xpValue"`
}

// Creates a Mob from a MobInfo instance
func (m *MobInfo) Spawn() *Mob {
	return &Mob{
		Name:            m.Name,
		IdleDescription: m.IdleDescription,
		TargetingNames:  m.TargetingNames,
		damageDice:      m.DamageDice,
		currHealth:      m.MaxHealth,
		maxHealth:       m.MaxHealth,
		xpValue:         m.XpValue,
	}
}

// An individual mob spawn
type Mob struct {
	spawnInfo       *SpawnInfo
	Name            string
	IdleDescription string
	TargetingNames  []string
	damageDice      dice.Dice
	currHealth      int
	maxHealth       int
	xpValue         int
	Dead            bool
	sync.Mutex
}

func (m *Mob) TakeDamage(dam int) {
	m.Lock()
	defer m.Unlock()
	m.currHealth = m.currHealth - dam
	if m.currHealth <= 0 {
		m.Dead = true
	}
}

func (m *Mob) Damage() (dealtDamage int) {
	return m.damageDice.Roll()
}

func (m *Mob) XpValue() int {
	return m.xpValue
}

func (m *Mob) HandleDeath() {
	if m.spawnInfo != nil {
		m.spawnInfo.HandleDeath()
	}
}

func (m *Mob) Tick() {}
