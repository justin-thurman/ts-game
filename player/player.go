package player

type Player struct {
	Name string
}

func New(name string) *Player {
	return &Player{Name: name}
}
