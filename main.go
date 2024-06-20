package main

import (
	"os"
	"ts-game/engine"
)

func main() {
	server := engine.New()
	server.Connect(os.Stdin, os.Stdout)
	server.Start()
}
