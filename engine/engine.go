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
	"ts-game/room"
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
		// TODO: validate name, go back to loop if invalid
		player = playerModule.New(name, r, w, exitCallback)
		break
	}
	s.playerLock.Lock()
	s.players = append(s.players, player)
	s.playerLock.Unlock()
	log.Info("User connected", "user", player.Name, "clientCount", len(s.players))
	player.Send("Welcome to my very professional game, %s!", player.Name)
	var playerRoom *room.Room
	var err error
	playerRoom, err = room.FindRoomById(player.RoomId)
	if err != nil {
		log.Error("Player saved to invalid room", "player", player.Name, "roomId", player.RoomId)
		playerRoom, err = room.FindRoomById(1)
		if err != nil {
			log.Error("Room 1 not found during player login", "player", player.Name)
			player.Send("Internal server error")
			exitCallback()
			return
		}
	}
	playerRoom.AddPlayer(player)
	player.Send(player.Location().HandleLook())
	go s.listenForCommands(player)
}

func (s *server) Start() error {
	err := room.Load()
	if err != nil {
		log.Info(err.Error())
		return err
	}
	starterMob := mob.New("ant")
	townCenter := room.Rooms[0]
	townCenter.AddMob(starterMob)
	for {
		for _, r := range room.Rooms {
			go r.Tick()
		}
		time.Sleep(time.Second * 4)
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
		case strings.HasPrefix("quit", cmd):
			// Logic handled below in order to also save users on disconnect
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
		case strings.HasPrefix("look", cmd):
			p.Send(p.Location().HandleLook())
		case strings.HasPrefix("kill", cmd):
			if strings.TrimSpace(cmdArgs) == "" {
				p.Send("Who do you want to kill?")
				break
			}
			p.Location().HandleKill(p, cmdArgs)
		default:
			p.Send("Unknown command: %s\n", cmd)
		}
	}
	err := scanner.Err()
	if err != nil {
		log.Info("Read error", "error", err.Error())
	}
	s.playerLock.Lock()
	s.players = slices.DeleteFunc(s.players, func(player *playerModule.Player) bool { return player == p })
	s.playerLock.Unlock()
	p.Quit()
	log.Info("User exit", "user", p.Name, "clientCount", len(s.players))
}
