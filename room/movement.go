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

// Parses a string, returning a movement direction and a boolean indicating whether the direction
// parsed successfully.
func ParseMovementDirection(s string) direction {
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

func (r *Room) HandleMovement(p *player.Player, direction string) {
	if r.PlayerIsInCombat(p) {
		p.Send("You're too busy fighting for your life!")
		return
	}
	dir := ParseMovementDirection(direction)
	if dir == invalid {
		p.Send("Go where?")
		return
	}
	destId, found := r.Exits[dir]
	if !found {
		p.Send("You can't go %s.", string(dir))
		return
	}
	r.movePlayer(p, destId)
}

func (r *Room) HandleRecall(p *player.Player, destinationId int) {
	if r.PlayerIsInCombat(p) {
		p.Send("You're too busy fighting for your life!")
		return
	}
	r.movePlayer(p, destinationId)
}

func (r *Room) movePlayer(p *player.Player, destId int) {
	dest, err := FindRoomById(destId)
	if err != nil {
		slog.Error("Error finding room during player movement", "player", p.Name, "startingRoom", r.Id, "destinationRoom", destId)
		p.Send("Internal server error finding room")
		return
	}
	r.RemovePlayer(p)
	dest.AddPlayer(p)
}
