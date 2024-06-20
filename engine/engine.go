package engine

import (
	"bufio"
	"fmt"
	"io"
	log "log/slog"
	"slices"
	"sync"
	"time"
	playerModule "ts-game/player"
)

type server struct {
	clients    []*client
	clientLock sync.RWMutex
}

type client struct {
	io.Reader
	io.Writer
	exitCallback func()
	player       *playerModule.Player
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
		player = playerModule.New(name)
		break
	}
	c := client{r, w, exitCallback, player}
	s.clients = append(s.clients, &c)
	log.Info("User connected", "user", player.Name, "clientCount", len(s.clients))
	go s.listenForCommands(&c)
}

func (s *server) Start() error {
	round := 1
	for {
		for _, c := range s.clients {
			fmt.Fprintf(c, "Beginning round: %d\n", round)
		}
		round++
		time.Sleep(time.Second * 6)
	}
}

func (s *server) listenForCommands(c *client) {
	scanner := bufio.NewScanner(c)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "exit" || line == "quit" {
			c.exitCallback()
			s.clientLock.Lock()
			s.clients = slices.DeleteFunc(s.clients, func(item *client) bool { return item == c })
			s.clientLock.Unlock()
			log.Info("User exit", "user", c.player.Name, "clientCount", len(s.clients))
			break
		}
		fmt.Fprintf(c, "You entered: %s\n", line)
	}
	err := scanner.Err()
	if err != nil {
		fmt.Fprintf(c, "Read error: %v\n", err.Error())
	}
}
