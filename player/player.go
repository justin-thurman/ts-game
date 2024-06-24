package player

import (
	"fmt"
	"io"
	"math/rand/v2"
	"sync"
	"ts-game/mob"
)

// Returns a random integer in the closed range [min, max]
func randRange(min, max int) int {
	return rand.IntN(max+1-min) + min
}

const PROMPT string = ">>> "

type location interface {
	HandleLook() string
	HandleKill(*Player, string)
}

type Player struct {
	io.Reader
	io.Writer
	exitCallback func()
	Name         string
	minDamage    int
	maxDamage    int
	currHealth   int
	maxHealth    int
	inCombat     bool
	Location     location
	sync.Mutex
}

func New(name string, r io.Reader, w io.Writer, exitCallback func()) *Player {
	return &Player{
		Name:         name,
		Reader:       r,
		Writer:       w,
		exitCallback: exitCallback,
		minDamage:    1,
		maxDamage:    8,
		currHealth:   30,
		maxHealth:    30,
		inCombat:     false,
	}
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
}

func (p *Player) EnterCombat(m *mob.Mob) {
	p.Send("You begin to fight %s!", m.Name)
	p.inCombat = true
}

func (p *Player) LeaveCombat() {
	p.inCombat = false
}

func (p *Player) GetDamage() int {
	return randRange(p.minDamage, p.maxDamage)
}

func (p *Player) TakeDamage(damage int) {
	p.Lock()
	defer p.Unlock()
	p.currHealth -= damage
}
