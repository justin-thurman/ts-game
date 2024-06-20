package main

import (
	"fmt"
	"log"
	"os"
	"ts-game/engine"
)

func main() {
	server := engine.New()
	if len(os.Args) == 1 {
		fmt.Println("No mode provided, defaulting to local")
		server.Connect(os.Stdin, os.Stdout)
	} else if os.Args[1] == "local" {
		server.Connect(os.Stdin, os.Stdout)
	} else if os.Args[1] == "telnet" {
		fmt.Println("TODO: Handle telnet")
	} else {
		log.Fatalf("Unknown mode: %v", os.Args[1])
	}
	server.Start()
}
