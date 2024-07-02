package mob

import (
	"math/rand/v2"
	"sync"

	_ "gopkg.in/yaml.v3"
)

// Returns a random integer in the closed range [min, max]
func randRange(min, max int) int {
	return rand.IntN(max+1-min) + min
}

// The template that an individual mob is spawned from
type MobInfo struct {
	Name            string   `yaml:"name"`
	IdleDescription string   `yaml:"idleDescription"`
	TargetingNames  []string `yaml:"targetingNames"`
	MinDamage       int      `yaml:"minDamage"`
	MaxDamage       int      `yaml:"maxDamage"`
	MaxHealth       int      `yaml:"health"`
	XpValue         int      `yaml:"xpValue"`
}

// Creates a Mob from a MobInfo instance
func (m *MobInfo) Spawn() Mob {
	return Mob{
		Name:            m.Name,
		IdleDescription: m.IdleDescription,
		TargetingNames:  m.TargetingNames,
		minDamage:       m.MinDamage,
		maxDamage:       m.MaxDamage,
		currHealth:      m.MaxHealth,
		maxHealth:       m.MaxHealth,
		xpValue:         m.XpValue,
	}
}

// An individual mob spawn
type Mob struct {
	Name            string
	IdleDescription string
	TargetingNames  []string
	minDamage       int
	maxDamage       int
	currHealth      int
	maxHealth       int
	xpValue         int
	Dead            bool
	sync.Mutex
}

func New(name string) *Mob {
	return &Mob{Name: name, minDamage: 1, maxDamage: 3, currHealth: 10, maxHealth: 10, xpValue: 10}
}

func (m *Mob) TakeDamage(dam int) {
	m.Lock()
	defer m.Unlock()
	m.currHealth = m.currHealth - dam
	if m.currHealth <= 0 {
		m.Dead = true
	}
}

func (m *Mob) Damage() int {
	return randRange(m.minDamage, m.maxDamage)
}

func (m *Mob) XpValue() int {
	return m.xpValue
}

func (m *Mob) Tick() {}
