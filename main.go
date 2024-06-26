package main

import (
	log "log/slog"
	"net"
	"os"
	"ts-game/engine"
)

func main() {
	server := engine.New()
	log.Info("Listening for TCP connections...")
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	go func() {
		for {
			conn, err := listener.Accept()
			log.Info("Got connection")
			if err != nil {
				log.Error("Error accepting connection", "error", err.Error())
				continue
			}
			server.Connect(conn, conn, func() { conn.Close() })
		}
	}()
	server.Start()
}
