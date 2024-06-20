package player

import (
	"fmt"
	"io"
)

type Player struct {
	io.Reader
	io.Writer
	exitCallback func()
	Name         string
}

func New(name string, r io.Reader, w io.Writer, exitCallback func()) *Player {
	return &Player{Name: name, Reader: r, Writer: w, exitCallback: exitCallback}
}

func (p *Player) Quit() {
	p.save()
	fmt.Fprintf(p, "Goodbye, %s!\n", p.Name)
	p.exitCallback()
}

func (p *Player) save() {
	fmt.Fprintln(p, "If we had persistence, we'd be saving your character now.")
}
