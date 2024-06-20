package player

import (
	"fmt"
	"io"
)

const PROMPT string = ">>> "

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
	p.Send("Goodbye, %s!\n", p.Name)
	p.exitCallback()
}

func (p *Player) save() {
	p.Send("If we had persistence, we'd be saving your character now.")
}

func (p *Player) Send(msg string, a ...any) {
	fmt.Fprintf(p, msg, a...)
	fmt.Fprintln(p, "")
	fmt.Fprint(p, PROMPT)
}
