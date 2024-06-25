package player

import (
	"fmt"
	"io"
	"math"
	"math/rand/v2"
	"sync"
	"ts-game/mob"
)

// Returns a random integer in the closed range [min, max]
func randRange(min, max int) int {
	return rand.IntN(max+1-min) + min
}

const PROMPT string = "%d/%d HP %d/%d XP >>> "

type location interface {
	HandleLook() string
	HandleKill(*Player, string)
}

type Player struct {
	io.Reader
	io.Writer
	location location
	sync.Mutex
	exitCallback func()
	Name         string
	minDamage    int
	maxDamage    int
	currHealth   int
	maxHealth    int
	currXp       int
	xpTolevel    int
	level        int
	inCombat     bool
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
		level:        1,
		xpTolevel:    xpToLevel(1),
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
	fmt.Fprintf(p, PROMPT, p.currHealth, p.maxHealth, p.currXp, p.xpTolevel)
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

func (p *Player) Damage() int {
	return randRange(p.minDamage, p.maxDamage)
}

func (p *Player) TakeDamage(damage int) {
	p.Lock()
	defer p.Unlock()
	p.currHealth -= damage
}

func (p *Player) Location() location {
	return p.location
}

func (p *Player) SetLocation(l location) {
	p.location = l
}

func (p *Player) GainXp(xp int) {
	p.currXp += xp
	p.Send("You gain %d experience!", xp)
	if p.currXp >= p.xpTolevel {
		p.levelUp()
	}
}

func xpToLevel(level int) int {
	return int(100 * math.Pow(float64(level), 1.5))
}

func (p *Player) levelUp() {
	p.xpTolevel = xpToLevel(p.level)
	p.Send("PLACEHOLDER: You leveld up")
}
