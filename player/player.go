package player

import (
	"fmt"
	"io"
	"math"
	"strings"
	"sync"
	"ts-game/classes"
	"ts-game/dice"
	"ts-game/items"
	"ts-game/stats"
)

const PROMPT string = "%d/%d HP %d/%d XP >>> "

type class interface {
	StartingStats() *stats.Stats
	HitDice() *dice.Dice
	String() string
}

type Player struct {
	exitCallback func()
	io.Reader
	io.Writer
	class             class
	stats             *stats.Stats
	Equip             *items.EquipInfo
	Inventory         *items.Inventory
	Name              string
	msgBuffer         strings.Builder
	hitDice           dice.Dice
	CurrHealth        int
	MaxHealth         int
	currXp            int
	xpTolevel         int
	level             int
	RoomId            int
	RecallRoomId      int
	HasActedThisRound bool
	mu                sync.Mutex
	// TODO: Will eventually need a queued command for skills and spells, to go off on next combat round
}

func New(name string, r io.Reader, w io.Writer, exitCallback func()) *Player {
	class := &classes.Warrior{}
	startingStats := class.StartingStats()
	// TODO: Not sure how to handle health. I think I want players health to scale faster than mobs
	hitDice := class.HitDice()
	startingHealth := hitDice.Max() + startingStats.ConModifier + 15 // 15 as extra base for now

	return &Player{
		Name:         name,
		Reader:       r,
		Writer:       w,
		exitCallback: exitCallback,
		class:        class,
		Equip:        class.StartingEquipment(),
		Inventory:    class.StartingInventory(),
		stats:        startingStats,
		hitDice:      *hitDice,
		CurrHealth:   startingHealth,
		MaxHealth:    startingHealth,
		level:        1,
		RoomId:       1,
		RecallRoomId: 1,
		xpTolevel:    xpToLevel(1),
	}
}

func (p *Player) Quit() {
	p.save()
	p.Send("Goodbye, %s!\n", p.Name)
	p.exitCallback()
}

func (p *Player) Death() {
	p.mu.Lock()
	defer p.mu.Unlock()
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

func (p *Player) Tick(inCombat bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if !inCombat {
		p.CurrHealth += 5 // TODO: health regen
		if p.CurrHealth >= p.MaxHealth {
			p.CurrHealth = p.MaxHealth
		}
	}
}

func (p *Player) Damage() int {
	return p.Equip.Damage() + p.stats.StrModifier
}

func (p *Player) TakeDamage(damage int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.CurrHealth -= damage
}

func (p *Player) SetRoomId(id int) {
	p.RoomId = id
}

func (p *Player) GainXp(xp int) {
	p.mu.Lock()
	defer p.mu.Unlock()
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
	// Increase level by 1
	p.level += 1
	// health gain
	healthGain := p.hitDice.AverageN(3) + p.stats.ConModifier
	p.MaxHealth += healthGain
	p.CurrHealth = p.MaxHealth
	// reset currXp, with carry over, and set next xpTolevel
	p.currXp -= p.xpTolevel
	p.xpTolevel = xpToLevel(p.level)
	p.BufferMsg("You gained a level! You gained %d health!", healthGain)
}

// Score returns the string representing the player's character sheet, for use in the score command.
func (p *Player) Score() string {
	scoreString := `Name: %s
  Class: %s
  Level: %d
  %s
  Health: %d/%d
  XP: %d/%d`
	return fmt.Sprintf(scoreString, p.Name, p.class.String(), p.level, p.stats.String(), p.CurrHealth, p.MaxHealth, p.currXp, p.xpTolevel)
}
