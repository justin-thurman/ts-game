package main

import (
	"os"
	"ts-game/engine"
)

func main() {
	engine.Hello()
	engine.Run(os.Stdin, os.Stdout)
}
