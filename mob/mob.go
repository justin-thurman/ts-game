package mob

import (
	"log/slog"
	"sync"
	"ts-game/dice"

	_ "gopkg.in/yaml.v3"
)

// The template that an individual mob is spawned from
type MobInfo struct {
	Name            string    `yaml:"name"`
	IdleDescription string    `yaml:"idleDescription"`
	TargetingNames  []string  `yaml:"targetingNames"`
	Level           int       `yaml:"level"`
	DamageDice      dice.Dice `yaml:"damageDice"`
	HitDice         dice.Dice `yaml:"hitDice"`
	XpValue         int       `yaml:"xpValue"`
}

func (m *MobInfo) getHealth() int {
	// Using a methodology similar to D&D. Here's the process:
	// For level 1, take the max value of the mob's hit dice.
	// For subsequent levels, rolls the dice 3 times and take the average value.
	// Add constitution modifier to each roll.
	// Health is the sum of all rolls.
	conModifier := 0 // TODO: will have to calculate this later
	health := m.HitDice.Max() + conModifier
	for i := 1; i < m.Level; i++ {
		health += m.HitDice.AverageN(3) + conModifier
	}
	return health
}

// Creates a Mob from a MobInfo instance
func (m *MobInfo) Spawn() *Mob {
	health := m.getHealth()
	slog.Debug("spawn", "mob", m.Name, "health", health)
	return &Mob{
		Name:            m.Name,
		IdleDescription: m.IdleDescription,
		TargetingNames:  m.TargetingNames,
		damageDice:      m.DamageDice,
		currHealth:      health,
		maxHealth:       health,
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
