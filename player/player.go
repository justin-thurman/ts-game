package player

import (
	"fmt"
	"io"
	"math/rand/v2"
	"ts-game/mob"
)

// Returns a random integer in the closed range [min, max]
func randRange(min, max int) int {
	return rand.IntN(max+1-min) + min
}

const PROMPT string = ">>> "

type Player struct {
	fighting   *mob.Mob
	killedMobs chan *mob.Mob
	io.Reader
	io.Writer
	exitCallback func()
	Name         string
	minDamage    int
	maxDamage    int
	currHealth   int
	maxHealth    int
}

func New(name string, r io.Reader, w io.Writer, exitCallback func(), killedMobs chan *mob.Mob) *Player {
	return &Player{Name: name, Reader: r, Writer: w, exitCallback: exitCallback, minDamage: 1, maxDamage: 8, currHealth: 30, maxHealth: 30, killedMobs: killedMobs}
}

func (p *Player) Quit() {
	p.save()
	p.Send("Goodbye, %s!\n", p.Name)
	p.exitCallback()
}

func (p *Player) save() {
	p.Send("If we had persistence, we'd be saving your character now.")
}

func (p *Player) Send(msg string, a ...any) {
	fmt.Fprintf(p, msg, a...)
	fmt.Fprintln(p, "")
	fmt.Fprint(p, PROMPT)
}

func (p *Player) Tick() {
	p.Send("Ticking...")
	if p.fighting == nil {
		// Then we're regening
		if p.currHealth < p.maxHealth {
			p.maxHealth = p.maxHealth + 1
		}
	} else {
		// Then we're fighting
		playerDamage := randRange(p.minDamage, p.maxDamage)
		p.fighting.TakeDamage(playerDamage)
		p.Send("You did %d damage!", playerDamage)
		if p.fighting.Dead {
			p.Send("You killed %s!", p.fighting.Name)
			p.killedMobs <- p.fighting
			p.fighting = nil
		} else {
			mobDamage := p.fighting.GetDamage()
			p.currHealth -= mobDamage
			p.Send("%s does %d damage to you! You have %d health remaining.", p.fighting.Name, mobDamage, p.currHealth)
			if p.currHealth <= 0 {
				p.Send("You died!")
				p.fighting = nil
			}
		}
	}
}

func (p *Player) BeginCombat(m *mob.Mob) {
	p.fighting = m
	p.Send("You begin to fight %s!", m.Name)
}
