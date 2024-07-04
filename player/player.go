package player

import (
	"fmt"
	"io"
	"math"
	"math/rand/v2"
	"strings"
	"sync"
)

// Returns a random integer in the closed range [min, max]
func randRange(min, max int) int {
	return rand.IntN(max+1-min) + min
}

const PROMPT string = "%d/%d HP %d/%d XP >>> "

type location interface {
	HandleKill(*Player, string)
	GetId() int
	HandleRecall(*Player, int)
	RemovePlayer(*Player)
}

type Player struct {
	io.Reader
	io.Writer
	location          location
	msgBuffer         strings.Builder
	exitCallback      func()
	Name              string
	minDamage         int
	maxDamage         int
	CurrHealth        int
	MaxHealth         int
	currXp            int
	xpTolevel         int
	level             int
	RoomId            int
	RecallRoomId      int
	HasActedThisRound bool
	sync.Mutex
	// TODO: Will eventually need a queued command for skills and spells, to go off on next combat round
}

func New(name string, r io.Reader, w io.Writer, exitCallback func()) *Player {
	return &Player{
		Name:         name,
		Reader:       r,
		Writer:       w,
		exitCallback: exitCallback,
		minDamage:    3,
		maxDamage:    8,
		CurrHealth:   10,
		MaxHealth:    30,
		level:        1,
		RoomId:       1,
		RecallRoomId: 1,
		xpTolevel:    xpToLevel(1),
	}
}

func (p *Player) Quit() {
	p.save()
	if p.location != nil {
		p.location.RemovePlayer(p)
	}
	p.Send("Goodbye, %s!\n", p.Name)
	p.exitCallback()
}

func (p *Player) Recall() {
	if p.RoomId == p.RecallRoomId {
		p.Send("You are already at your recall point.")
		return
	}
	p.location.HandleRecall(p, p.RecallRoomId)
}

func (p *Player) Death() {
	p.Lock()
	defer p.Unlock()
	p.currXp = p.currXp / 2
	p.CurrHealth = p.MaxHealth
}

func (p *Player) save() {
	p.Send("If we had persistence, we'd be saving your character now.")
}

func (p *Player) prompt() string {
	return fmt.Sprintf(PROMPT, p.CurrHealth, p.MaxHealth, p.currXp, p.xpTolevel)
}

func (p *Player) Send(msg string, a ...any) {
	fmt.Fprintf(p, msg, a...)
	fmt.Fprintf(p, "\n"+p.prompt())
}

func (p *Player) BufferMsg(msg string, a ...any) {
	p.msgBuffer.WriteString(fmt.Sprintf(msg+"\n", a...))
}

func (p *Player) SendBufferedMsgs() {
	if p.msgBuffer.Len() == 0 {
		return
	}
	fmt.Fprint(p, "\n"+p.msgBuffer.String()+p.prompt())
	p.msgBuffer.Reset()
}

func (p *Player) Tick() {
}

func (p *Player) Damage() int {
	return randRange(p.minDamage, p.maxDamage)
}

func (p *Player) TakeDamage(damage int) {
	p.Lock()
	defer p.Unlock()
	p.CurrHealth -= damage
}

func (p *Player) Location() location {
	return p.location
}

func (p *Player) SetLocation(l location) {
	p.location = l
	p.RoomId = l.GetId()
}

func (p *Player) GainXp(xp int) {
	p.currXp += xp
	p.BufferMsg("You gain %d experience!", xp)
	if p.currXp >= p.xpTolevel {
		p.levelUp()
	}
}

func xpToLevel(level int) int {
	return int(100 * math.Pow(float64(level), 1.5))
}

func (p *Player) levelUp() {
	p.xpTolevel = xpToLevel(p.level)
	p.BufferMsg("PLACEHOLDER: You leveld up")
}
