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
	description      string
	name             string
	players          []*player.Player
	mobs             []*mob.Mob
	battles          []*battle
	mobLock          sync.Mutex
	battleStartQueue chan *battleStart
	battleEndQueue   chan *battleEnd
	mobDeaths        chan *mob.Mob
}

type battle struct {
	players []*player.Player
	mobs    []*mob.Mob
}

type battleStart struct {
	player    *player.Player
	mob       *mob.Mob
	aggressor Combatant
}

type Combatant string

const (
	PlayerCombatant Combatant = "player"
	MobCombatant    Combatant = "mob"
)

type battleEnd struct {
	battle *battle
}

// Listens for activity in any channels that could impact inhabitants. Includes battle and movement.
func (r *Room) ListenForInhabitantChanges() {
	select {
	case startingBattle := <-r.battleStartQueue:
		if !slices.Contains(r.mobs, startingBattle.mob) {
			if startingBattle.aggressor == PlayerCombatant {
				startingBattle.player.Send("No one named %s here!", startingBattle.mob.Name)
			}
			break
		}
		var playerBattle *battle
		var mobBattle *battle
		for _, b := range r.battles {
			if slices.Contains(b.mobs, startingBattle.mob) {
				mobBattle = b
			}
			if slices.Contains(b.players, startingBattle.player) {
				playerBattle = b
			}
		}
		// Five possibilities: 1. both in the same battle, 2. both in different battles, 3. player in battle, 4. mob in battle, 5. neither in battle
		if playerBattle == nil && mobBattle == nil {
			newBattle := &battle{}
			newBattle.players = append(newBattle.players, startingBattle.player)
			newBattle.mobs = append(newBattle.mobs, startingBattle.mob)
			r.battles = append(r.battles, newBattle)
		} else if playerBattle == mobBattle {
			// I don't think this should happen, but if it does, we don't need to do anything
			slog.Warn("Starting a battle, but mob and player already fighting each other in same battle")
		} else if mobBattle != nil {
			mobBattle.players = append(mobBattle.players, startingBattle.player)
		} else if playerBattle != nil {
			playerBattle.mobs = append(playerBattle.mobs, startingBattle.mob)
		} else {
			slog.Error("Unexpected battle state when starting new battle")
		}
	case endingBattle := <-r.battleEndQueue:
		if len(r.battles) == 1 {
			r.battles = make([]*battle, 0)
			break
		}
		for i, b := range r.battles {
			if b == endingBattle.battle {
				r.battles = slices.Delete(r.battles, i, i+1)
				break
			}
		}
	}
	// mobs entering and leaving
}

func New(name, description string) *Room {
	return &Room{
		name:        name,
		description: description,
	}
}

func (r *Room) HandleLook() string {
	return fmt.Sprintf("%s\n%s", r.name, r.description)
}

func (r *Room) HandleKill(p *player.Player, mobName string) {
	var target *mob.Mob
	for _, tar := range r.mobs {
		if strings.HasPrefix(tar.Name, mobName) {
			target = tar
			break
		}
	}
	if target == nil {
		p.Send("No one named %s here!", mobName)
		return
	}
	r.battleStartQueue <- &battleStart{player: p, mob: target, aggressor: PlayerCombatant}
}

func (r *Room) Tick() {
	for _, b := range r.battles {
		b.tick()
	}
}

func (b *battle) tick() {
	// process player round
	for _, p := range b.players {
		damage := p.GetDamage()
		target := b.mobs[0]
		if target == nil {
			slog.Error("Mob was nil in combat")
		}
		target.TakeDamage(damage)
		if target.Dead {
			// TODO:handle killed mobs
			if len(b.mobs) > 1 {
				b.mobs = slices.Delete(b.mobs, 1, 2)
			} else {
				// TODO: combat is over
			}
		}
	}
	for _, m := range b.mobs {
		damage := m.GetDamage()
		target := b.players[0]
		if target == nil {
			slog.Error("Player was nil in combat")
		}
		target.TakeDamage(damage)
	}
}
