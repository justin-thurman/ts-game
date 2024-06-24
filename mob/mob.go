package mob

import (
	"math/rand/v2"
	"sync"
)

// Returns a random integer in the closed range [min, max]
func randRange(min, max int) int {
	return rand.IntN(max+1-min) + min
}

type Combatant interface {
	TakeDamage()
}

type Mob struct {
	Name       string
	minDamage  int
	maxDamage  int
	currHealth int
	maxHealth  int
	Dead       bool
	sync.Mutex
}

func New(name string) *Mob {
	return &Mob{Name: name, minDamage: 1, maxDamage: 3, currHealth: 10, maxHealth: 10}
}

func (m *Mob) TakeDamage(dam int) {
	m.Lock()
	defer m.Unlock()
	m.currHealth = m.currHealth - dam
	if m.currHealth <= 0 {
		m.Dead = true
	}
}

func (m *Mob) GetDamage() int {
	return randRange(m.minDamage, m.maxDamage)
}

func (m *Mob) Tick() {}
