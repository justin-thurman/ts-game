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
	playerModule "ts-game/player"
)

type server struct {
	players    []*playerModule.Player
	playerLock sync.RWMutex
}

func New() *server {
	return &server{}
}

func (s *server) Connect(r io.Reader, w io.Writer, exitCallback func()) {
	fmt.Fprintln(w, "Welcome! What is your name?")
	scanner := bufio.NewScanner(r)
	var player *playerModule.Player
	for scanner.Scan() {
		name := scanner.Text()
		player = playerModule.New(name, r, w, exitCallback)
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
	round := 1
	for {
		for _, c := range s.players {
			fmt.Fprintf(c, "Beginning round: %d\n", round)
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
		case cmd == "exit" || cmd == "quit":
			s.playerLock.Lock()
			s.players = slices.DeleteFunc(s.players, func(player *playerModule.Player) bool { return player == p })
			s.playerLock.Unlock()
			p.Quit()
			log.Info("User exit", "user", p.Name, "clientCount", len(s.players))
			break mainLoop
		case strings.HasPrefix("gossip", cmd):
			p.Send("You gossiped: %s\n", cmdArgs)
		default:
			p.Send("Unknown command: %s\n", cmd)
		}
	}
	err := scanner.Err()
	if err != nil {
		fmt.Fprintf(p, "Read error: %v\n", err.Error())
	}
}
