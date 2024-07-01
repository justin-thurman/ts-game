package room

import (
	"log/slog"
	"strings"
	"ts-game/player"
)

type direction string

const (
	north   direction = "north"
	south   direction = "south"
	east    direction = "east"
	west    direction = "west"
	up      direction = "up"
	down    direction = "down"
	invalid direction = "invalid"
)

func parseMovementDirection(s string) direction {
	s = strings.TrimSpace(s)
	switch {
	case strings.HasPrefix(string(north), s):
		return north
	case strings.HasPrefix(string(south), s):
		return south
	case strings.HasPrefix(string(east), s):
		return east
	case strings.HasPrefix(string(west), s):
		return west
	case strings.HasPrefix(string(up), s):
		return up
	case strings.HasPrefix(string(down), s):
		return down
	default:
		return invalid
	}
}

func (r *Room) HandleMovement(player *player.Player, direction string) {
	dir := parseMovementDirection(direction)
	if dir == invalid {
		player.Send("Go where?")
		return
	}
	destId, found := r.Exits[dir]
	if !found {
		player.Send("You can't go %s.", string(dir))
		return
	}
	dest, err := FindRoomById(destId)
	if err != nil {
		slog.Error("Error finding room during player movement", "player", player.Name, "startingRoom", r.Id, "destinationRoom", destId)
		player.Send("Internal server error finding room")
	}
	r.RemovePlayer(player)
	dest.AddPlayer(player)
	player.Send(dest.HandleLook())
}
