package room

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"sync"
	"ts-game/mob"
	"ts-game/player"

	_ "gopkg.in/yaml.v3"
)

type direction string

const (
	north direction = "north"
	south direction = "south"
	east  direction = "east"
	west  direction = "west"
	up    direction = "up"
	down  direction = "down"
)

type Room struct {
	players         map[*player.Player][]*mob.Mob
	mobs            map[*mob.Mob][]*player.Player // similar map of mobs to players
	Exits           map[direction]int             `yaml:"exits"`
	description     string
	DescriptionBase string `yaml:"description"`
	Name            string `yaml:"name"`
	Id              int    `yaml:"id"`
	sync.Mutex
}

func New(name, description string) *Room {
	room := &Room{
		Name:            name,
		DescriptionBase: description,
		mobs:            make(map[*mob.Mob][]*player.Player),
		players:         make(map[*player.Player][]*mob.Mob),
	}
	room.updateDescription()
	return room
}

func (r *Room) initialize() {
	r.mobs = make(map[*mob.Mob][]*player.Player)
	r.players = make(map[*player.Player][]*mob.Mob)
	r.updateDescription()
}

func (r *Room) GetId() int {
	return r.Id
}

func (r *Room) HandleLook() string {
	return r.description
}

func (r *Room) updateDescription() {
	nameAndDescription := fmt.Sprintf("%s\n%s", r.Name, r.DescriptionBase)
	descList := []string{nameAndDescription}
	for m, players := range r.mobs {
		var s string
		if len(players) == 0 {
			s = fmt.Sprintf("%s is standing here.", m.Name)
		} else {
			s = fmt.Sprintf("%s is fighting for its life!", m.Name)
		}
		descList = append(descList, s)
	}
	if len(r.Exits) == 0 {
		descList = append(descList, "Exits: None")
	} else {
		exitStrs := []string{"Exits:"}
		for exit := range r.Exits {
			exitStrs = append(exitStrs, string(exit))
		}
		descList = append(descList, strings.Join(exitStrs, " "))
	}
	r.description = strings.Join(descList, "\n")
}

func (r *Room) HandleKill(p *player.Player, mobName string) {
	r.Lock()
	defer r.Unlock()
	if !r.playerIsInRoom(p) {
		slog.Error("Player not in room when HandleKill command ran", "player", p, "room", r)
		return
	}
	var target *mob.Mob
	for tar := range r.mobs {
		if tar != nil && strings.HasPrefix(tar.Name, mobName) {
			target = tar
			break
		}
	}
	if target == nil {
		p.Send("No one named %s here!", mobName)
		return
	}
	if slices.Contains(r.players[p], target) {
		p.Send("You're doing your best!")
		return
	}
	p.BufferMsg("You begin to fight %s!", target.Name)
	defer p.SendBufferedMsgs()
	defer r.updateDescription()
	r.startCombat(p, target)
}

func (r *Room) Tick() {
	r.Lock()
	defer r.Unlock()
	defer r.updateDescription()
	// Handle player rounds
	for p, mobs := range r.players {
		defer p.SendBufferedMsgs()
		playerIsInCombat := len(mobs) > 0
		if !playerIsInCombat {
			p.Tick()
			continue
		}
		if p.HasActedThisRound {
			p.HasActedThisRound = false
			continue
		}
		target := mobs[0] // TODO: player will need control over this later; and AoE damage
		damage := p.Damage()
		target.TakeDamage(damage)
		p.BufferMsg("You deal %d damage to %s!", damage, target.Name)
		if target.Dead {
			r.removeMob(target)
			p.BufferMsg("You killed %s!", target.Name)
			p.GainXp(target.XpValue())
		}
		p.HasActedThisRound = false
	}
	// Handle mob rounds
	for m, players := range r.mobs {
		mobIsInCombat := len(players) > 0
		if !mobIsInCombat {
			m.Tick()
			continue
		}
		target := players[0] // TODO: will need an aggro system later
		damage := m.Damage()
		target.TakeDamage(damage)
		target.BufferMsg("%s dealt %d damage to you!", m.Name, damage)
		// TODO: handle player death
	}
}

func (r *Room) playerIsInRoom(p *player.Player) bool {
	_, found := r.players[p]
	return found
}

func (r *Room) mobIsInRoom(m *mob.Mob) bool {
	_, found := r.mobs[m]
	return found
}

func (r *Room) startCombat(p *player.Player, m *mob.Mob) {
	if !p.HasActedThisRound {
		p.HasActedThisRound = true
		damage := p.Damage()
		m.TakeDamage(damage)
		p.BufferMsg("You deal %d damage to %s!", damage, m.Name)
		if m.Dead {
			r.removeMob(m)
			p.BufferMsg("You killed %s!", m.Name)
			p.GainXp(m.XpValue())
			return
		}
	}
	r.players[p] = append(r.players[p], m)
	r.mobs[m] = append(r.mobs[m], p)
}

func (r *Room) removeMob(m *mob.Mob) {
	delete(r.mobs, m)
	for p, mobs := range r.players {
		for i, fightingMob := range mobs {
			if fightingMob == m {
				r.players[p] = slices.Delete(mobs, i, i+1)
				break
			}
		}
	}
}

func (r *Room) AddMob(m *mob.Mob) {
	r.Lock()
	defer r.Unlock()
	defer r.updateDescription()
	r.mobs[m] = []*player.Player{}
}

func (r *Room) RemovePlayer(p *player.Player) {
	r.Lock()
	defer r.Unlock()
	delete(r.players, p)
	for m, players := range r.mobs {
		for i, fightingPlayer := range players {
			if fightingPlayer == p {
				r.mobs[m] = slices.Delete(players, i, i+1)
				break
			}
		}
	}
}

func (r *Room) AddPlayer(p *player.Player) {
	r.Lock()
	defer r.Unlock()
	defer r.updateDescription()
	r.players[p] = []*mob.Mob{}
	p.SetLocation(r)
}
