package engine

import (
	"bufio"
	"fmt"
	"io"
	log "log/slog"
	"slices"
	"strings"
	"sync"
	"time"
	"ts-game/mob"
	playerModule "ts-game/player"
)

type server struct {
	players     []*playerModule.Player
	killedMobs  chan *mob.Mob
	spawnedMobs chan *mob.Mob
	mobs        []*mob.Mob
	playerLock  sync.RWMutex
}

func New() *server {
	return &server{killedMobs: make(chan *mob.Mob, 100), spawnedMobs: make(chan *mob.Mob, 100)}
}

func (s *server) Connect(r io.Reader, w io.Writer, exitCallback func()) {
	fmt.Fprintln(w, "Welcome! What is your name?")
	scanner := bufio.NewScanner(r)
	var player *playerModule.Player
	for scanner.Scan() {
		name := scanner.Text()
		player = playerModule.New(name, r, w, exitCallback, s.killedMobs)
		break
	}
	s.playerLock.Lock()
	s.players = append(s.players, player)
	s.playerLock.Unlock()
	log.Info("User connected", "user", player.Name, "clientCount", len(s.players))
	player.Send("Welcome to my very professional game, %s!\n", player.Name)
	go s.listenForCommands(player)
}

func (s *server) Start() error {
	starterMob := mob.New("ant")
	s.spawnedMobs <- starterMob
	go s.listenForMobs()
	round := 1
	for {
		for _, p := range s.players {
			go p.Tick()
		}
		round++
		time.Sleep(time.Second * 6)
	}
}

func (s *server) listenForCommands(p *playerModule.Player) {
	scanner := bufio.NewScanner(p)
mainLoop:
	for scanner.Scan() {
		cmd, cmdArgs, _ := strings.Cut(scanner.Text(), " ")
		switch {
		case cmd == "":
			p.Send("")
		case cmd == "exit" || cmd == "quit":
			s.playerLock.Lock()
			s.players = slices.DeleteFunc(s.players, func(player *playerModule.Player) bool { return player == p })
			s.playerLock.Unlock()
			p.Quit()
			log.Info("User exit", "user", p.Name, "clientCount", len(s.players))
			break mainLoop
		case strings.HasPrefix("gossip", cmd):
			if strings.TrimSpace(cmdArgs) == "" {
				p.Send("What do you want to gossip?")
				break
			}
			p.Send("You gossip, \"%s\"", cmdArgs)
			for _, player := range s.players {
				if player == p {
					continue
				}
				player.Send("%s gossips, \"%s\"", p.Name, cmdArgs)
			}
		case strings.HasPrefix("kill", cmd):
			if strings.TrimSpace(cmdArgs) == "" {
				p.Send("Who do you want to kill?")
				break
			}
			var target *mob.Mob
			for _, tar := range s.mobs {
				if strings.HasPrefix(tar.Name, cmdArgs) {
					target = tar
					break
				}
			}
			if target == nil {
				p.Send("No one named %s here!", cmdArgs)
				break
			}
			p.BeginCombat(target)
		default:
			p.Send("Unknown command: %s\n", cmd)
		}
	}
	err := scanner.Err()
	if err != nil {
		fmt.Fprintf(p, "Read error: %v\n", err.Error())
	}
}

func (s *server) listenForMobs() {
	for {
		select {
		case deadMob := <-s.killedMobs:
			idx := slices.Index(s.mobs, deadMob)
			if idx == -1 {
				log.Error("Dead mob not found", "mob", deadMob)
			}
			s.mobs = slices.Delete(s.mobs, idx, idx+1)
		case spawningMob := <-s.spawnedMobs:
			s.mobs = append(s.mobs, spawningMob)
		}
	}
}
