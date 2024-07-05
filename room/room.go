package room

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"sync"
	"ts-game/mob"
	"ts-game/player"

	orderedmap "github.com/wk8/go-ordered-map/v2"
	_ "gopkg.in/yaml.v3"
)

type Room struct {
	zone            *Zone
	players         map[*player.Player][]*mob.Mob
	mobs            *orderedmap.OrderedMap[*mob.Mob, []*player.Player]
	Exits           map[direction]int `yaml:"exits"`
	description     string
	DescriptionBase string `yaml:"description"`
	Name            string `yaml:"name"`
	Id              int    `yaml:"id"`
	sync.Mutex
}

func (r *Room) initialize() {
	r.mobs = orderedmap.New[*mob.Mob, []*player.Player]()
	r.players = make(map[*player.Player][]*mob.Mob)
	r.updateDescription()
}

func (r *Room) HandleLook() string {
	return r.description
}

func (r *Room) updateDescription() {
	nameAndDescription := fmt.Sprintf("%s\n%s", r.Name, r.DescriptionBase)
	descList := []string{nameAndDescription}
	for pair := r.mobs.Oldest(); pair != nil; pair = pair.Next() {
		m := pair.Key
		players := pair.Value
		var s string
		if len(players) == 0 {
			s = m.IdleDescription
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
	mobName = strings.ToLower(mobName)
	r.Lock()
	defer r.Unlock()
	if !r.playerIsInRoom(p) {
		slog.Error("Player not in room when HandleKill command ran", "player", p, "room", r)
		return
	}
	var target *mob.Mob
outerLoop:
	for pair := r.mobs.Oldest(); pair != nil; pair = pair.Next() {
		tar := pair.Key
		if tar == nil {
			continue
		}
		for _, tarName := range tar.TargetingNames {
			if strings.HasPrefix(tarName, mobName) {
				target = tar
				break outerLoop
			}
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
			r.handlePlayerKilledMob(p, target)
		}
		p.HasActedThisRound = false
	}
	// Handle mob rounds
	for pair := r.mobs.Oldest(); pair != nil; pair = pair.Next() {
		m := pair.Key
		players := pair.Value
		mobIsInCombat := len(players) > 0
		if !mobIsInCombat {
			m.Tick()
			continue
		}
		target := players[0] // TODO: will need an aggro system later
		damage := m.Damage()
		target.TakeDamage(damage)
		target.BufferMsg("%s dealt %d damage to you!", m.Name, damage)
		if target.CurrHealth <= 0 {
			target.Death()
			target.BufferMsg("You died!")
			r.movePlayer(target, target.RecallRoomId)
		}
	}
}

func (r *Room) playerIsInRoom(p *player.Player) bool {
	_, found := r.players[p]
	return found
}

func (r *Room) PlayerIsInCombat(p *player.Player) bool {
	r.Lock()
	defer r.Unlock()
	mobs, found := r.players[p]
	if !found {
		slog.Error("Player not found in room when checking if player is in combat", "player", p.Name, "roomId", r.Id)
		return false
	}
	return len(mobs) > 0
}

func (r *Room) mobIsInRoom(m *mob.Mob) bool {
	_, found := r.mobs.Get(m)
	return found
}

func (r *Room) startCombat(p *player.Player, m *mob.Mob) {
	if !p.HasActedThisRound {
		p.HasActedThisRound = true
		damage := p.Damage()
		m.TakeDamage(damage)
		p.BufferMsg("You deal %d damage to %s!", damage, m.Name)
		if m.Dead {
			r.handlePlayerKilledMob(p, m)
			return
		}
	}
	r.players[p] = append(r.players[p], m)
	playerSlice, found := r.mobs.Get(m)
	if !found {
		slog.Error("Mob not found in orderedmap", "mob", m)
	}
	r.mobs.Set(m, append(playerSlice, p))
}

func (r *Room) handlePlayerKilledMob(p *player.Player, m *mob.Mob) {
	r.zone.handleMobDeath()
	r.removeMob(m)
	p.BufferMsg("You killed %s!", m.Name)
	p.GainXp(m.XpValue())
	m.HandleDeath()
}

func (r *Room) removeMob(m *mob.Mob) {
	r.mobs.Delete(m)
	for p, mobs := range r.players {
		for i, fightingMob := range mobs {
			if fightingMob == m {
				r.players[p] = slices.Delete(mobs, i, i+1)
				break
			}
		}
	}
}

func (r *Room) addMob(m *mob.Mob) {
	r.Lock()
	defer r.Unlock()
	defer r.updateDescription()
	r.mobs.Set(m, []*player.Player{})
}

func (r *Room) RemovePlayer(p *player.Player) {
	r.Lock()
	defer r.Unlock()
	slog.Debug("Room player count before removal", "roomId", r.Id, "playerCount", len(r.players))
	delete(r.players, p)
	for pair := r.mobs.Oldest(); pair != nil; pair = pair.Next() {
		m := pair.Key
		players := pair.Value
		for i, fightingPlayer := range players {
			if fightingPlayer == p {
				r.mobs.Set(m, slices.Delete(players, i, i+1))
				break
			}
		}
	}
	slog.Debug("Room player count after removal", "roomId", r.Id, "playerCount", len(r.players))
}

func (r *Room) AddPlayer(p *player.Player) {
	r.Lock()
	defer r.Unlock()
	defer r.updateDescription()
	r.players[p] = []*mob.Mob{}
	p.SetRoomId(r.Id)
	p.Send(r.HandleLook())
}
