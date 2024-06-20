package main

import (
	"fmt"
	log "log/slog"
	"net"
	"os"
	"ts-game/engine"
)

func main() {
	server := engine.New()
	if len(os.Args) == 1 || os.Args[1] == "telnet" {
		log.Info("Listening for TCP connections...")
		listener, err := net.Listen("tcp", ":8081")
		if err != nil {
			log.Error(err.Error())
			os.Exit(1)
		}
		go func() {
			for {
				conn, err := listener.Accept()
				if err != nil {
					log.Error("Error accepting connection", "error", err.Error())
					continue
				}
				server.Connect(conn, conn, func() { conn.Close() })
			}
		}()
	} else if os.Args[1] == "local" {
		server.Connect(os.Stdin, os.Stdout, func() { fmt.Println("Goodbye!") })
	} else {
		log.Error("Unknown mode", "mode", os.Args[1])
		os.Exit(1)
	}
	server.Start()
}
