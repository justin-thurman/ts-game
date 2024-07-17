package engine

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	log "log/slog"
	"slices"
	"strings"
	"sync"
	"time"
	"ts-game/auth"
	"ts-game/db/queries"
	"ts-game/help"
	"ts-game/items"
	playerModule "ts-game/player"
	"ts-game/room"
)

type server struct {
	queryEngine *queries.Queries
	players     []*playerModule.Player
	playerLock  sync.RWMutex
}

func New(queryEngine *queries.Queries) *server {
	return &server{queryEngine: queryEngine}
}

func (s *server) Connect(r io.Reader, w io.Writer, exitCallback func()) {
	ctx := context.Background()
	scanner := bufio.NewScanner(r)
	fmt.Fprintln(w, "Welcome! Enter your account name to login to an existing account or create a new one.")
	for scanner.Scan() {
		accountName := scanner.Text()
		accountExists, err := auth.AccountExists(ctx, s.queryEngine, accountName)
		if err != nil {
			fmt.Fprintln(w, "Error searching for account. Please try again.")
			continue
		}
		if accountExists {
			fmt.Fprintf(w, "Logging into account %s. Enter password.\n", accountName)
			scanner.Scan()
			password := scanner.Text()
			accountId, err := auth.Login(ctx, s.queryEngine, accountName, password)
			if err != nil {
				if err.Error() == "incorrect password" {
					fmt.Fprintln(w, "Incorrect password.")
				} else {
					fmt.Fprintln(w, "Error logging in. Please try again.")
					slog.Error("Error during login", "err", err, "accountName", accountName)
				}
				fmt.Fprintln(w, "Welcome! Enter your account name to login to an existing account or create a new one.")
				continue
			}
			slog.Debug("Login to account", "accountId", accountId)
			fmt.Fprintf(w, "Welcome back, %s!", accountName)
		} else {
			fmt.Fprintf(w, "Creating account with name %s. Continue? 'yes' or 'no'\n", accountName)
			scanner.Scan()
			answer := scanner.Text()
			if answer == "yes" {
				password := ""
				password2 := "not matching"
				for password != password2 {
					fmt.Fprintln(w, "Please enter your password.")
					scanner.Scan()
					password = scanner.Text()
					fmt.Fprintln(w, "Please enter your password again.")
					scanner.Scan()
					password2 = scanner.Text()
					if password != password2 {
						fmt.Fprintln(w, "Passwords do not match")
					}
				}
				accountId, err := auth.CreateAccount(ctx, s.queryEngine, accountName, password)
				if err != nil {
					fmt.Fprintln(w, "Error creating account. Please try again.")
					slog.Error("Error during account creation", "err", err)
				}
				slog.Debug("Account created", "accountId", accountId)
				fmt.Fprintln(w, "Account created successfully!")
			} else {
				fmt.Fprintln(w, "Welcome! Enter your account name to login to an existing account or create a new one.")
				continue
			}
		}
	}
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
	go s.listenForCommands(player)
}

func (s *server) Start() error {
	err := items.Load()
	if err != nil {
		log.Info(err.Error())
		return err
	}
	err = room.Load()
	if err != nil {
		log.Info(err.Error())
		return err
	}
	for {
		for _, r := range room.Rooms {
			go r.Tick()
		}
		for _, z := range room.Zones {
			go z.Tick()
		}
		time.Sleep(time.Second * 4)
	}
}

func (s *server) listenForCommands(p *playerModule.Player) {
	scanner := bufio.NewScanner(p)
mainLoop:
	for scanner.Scan() {
		playerRoom, err := room.FindRoomById(p.RoomId)
		cmd, cmdArgs, _ := strings.Cut(scanner.Text(), " ")
		switch {
		case cmd == "":
			p.Send("")
		case strings.HasPrefix("north", cmd):
			playerRoom.HandleMovement(p, cmd)
		case strings.HasPrefix("south", cmd):
			playerRoom.HandleMovement(p, cmd)
		case strings.HasPrefix("east", cmd):
			playerRoom.HandleMovement(p, cmd)
		case strings.HasPrefix("west", cmd):
			playerRoom.HandleMovement(p, cmd)
		case strings.HasPrefix("up", cmd):
			playerRoom.HandleMovement(p, cmd)
		case strings.HasPrefix("down", cmd):
			playerRoom.HandleMovement(p, cmd)
		case strings.HasPrefix("help", cmd):
			p.Send(help.HELP)
		case strings.HasPrefix("recall", cmd):
			playerRoom.HandleRecall(p, p.RecallRoomId)
		case strings.HasPrefix("equipment", cmd):
			p.Send(p.Equip.String())
		case strings.HasPrefix("inventory", cmd):
			p.Send(p.Inventory.String())
		case strings.HasPrefix("wear", cmd):
			p.Send(p.Inventory.Wear(cmdArgs, p.Equip))
			p.UpdateStats()
		case strings.HasPrefix("wield", cmd):
			p.Send(p.Inventory.Wield(cmdArgs, p.Equip))
			p.UpdateStats()
		case strings.HasPrefix("remove", cmd):
			p.Send(p.Inventory.Remove(cmdArgs, p.Equip))
			p.UpdateStats()
		case strings.HasPrefix("get", cmd):
			p.Send(playerRoom.HandleGet(p, cmdArgs))
		case strings.HasPrefix("drop", cmd):
			p.Send(playerRoom.HandleDrop(p, cmdArgs))
		case strings.HasPrefix("quit", cmd):
			if err != nil {
				log.Error("Player room not found during quit", "player", p.Name, "roomId", p.RoomId)
				p.Send("Internal server error during quit. Try again in a few moments.")
				break
			}
			if playerRoom.PlayerIsInCombat(p) {
				p.Send("You can't quit now. You're fighting for your life!")
				break
			}
			// Logic handled below in order to also save users on disconnect
			break mainLoop
		case strings.HasPrefix("score", cmd):
			p.Send(p.Score())
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
			p.Send(playerRoom.HandleLook())
		case strings.HasPrefix("kill", cmd):
			if strings.TrimSpace(cmdArgs) == "" {
				p.Send("Who do you want to kill?")
				break
			}
			playerRoom.HandleKill(p, cmdArgs)
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
	playerRoom, err := room.FindRoomById(p.RoomId)
	if err != nil {
		log.Error("Error finding player room on quit", "player", p.Name, "roomId", p.RoomId)
	}
	if playerRoom != nil {
		playerRoom.RemovePlayer(p)
	}
	p.Quit()
	log.Info("User exit", "user", p.Name, "clientCount", len(s.players))
}
