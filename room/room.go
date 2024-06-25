package room

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"sync"
	"ts-game/mob"
	"ts-game/player"
)

type Room struct {
	description string
	name        string
	players     map[*player.Player][]*mob.Mob
	mobs        map[*mob.Mob][]*player.Player // similar map of mobs to players
	sync.Mutex
}

func New(name, description string) *Room {
	return &Room{
		name:        name,
		description: description,
		mobs:        make(map[*mob.Mob][]*player.Player),
		players:     make(map[*player.Player][]*mob.Mob),
	}
}

func (r *Room) HandleLook() string {
	r.Lock()
	defer r.Unlock()
	// FIX: This lock is awful. Find another way
	nameAndDescription := fmt.Sprintf("%s\n%s", r.name, r.description)
	descList := []string{nameAndDescription}
	for m := range r.mobs {
		descList = append(descList, fmt.Sprintf("%s is standing here.", m.Name))
	}
	return strings.Join(descList, "\n")
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
	r.startCombat(p, target)
	p.Send("You begin to fight %s!", target.Name)
}

func (r *Room) Tick() {
	slog.Info("Room ticking")
	r.Lock()
	defer r.Unlock()
	// Handle player rounds
	for p, mobs := range r.players {
		playerIsInCombat := len(mobs) > 0
		if !playerIsInCombat {
			p.Tick()
			continue
		}
		target := mobs[0] // TODO: player will need control over this later; and AoE damage
		damage := p.Damage()
		target.TakeDamage(damage)
		p.Send("You deal %d damage to %s!", damage, target.Name)
		if target.Dead {
			r.removeMob(target)
			p.Send("You killed %s!", target.Name)
			p.GainXp(target.XpValue())
		}
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
		target.Send("%s dealt %d damage to you!", m.Name, damage)
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
	r.mobs[m] = []*player.Player{}
}

func (r *Room) AddPlayer(p *player.Player) {
	r.Lock()
	defer r.Unlock()
	r.players[p] = []*mob.Mob{}
	p.SetLocation(r)
}
