package main

import (
	"log"
	"log/slog"
	"net"
	"os"
	"ts-game/engine"
)

func main() {
	server := engine.New()
	slog.Info("Listening for TCP connections...")
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	go func() {
		for {
			conn, err := listener.Accept()
			slog.Info("Got connection")
			if err != nil {
				slog.Error("Error accepting connection", "error", err.Error())
				continue
			}
			server.Connect(conn, conn, func() { conn.Close() })
		}
	}()
	err = server.Start()
	if err != nil {
		log.Fatal(err.Error())
	}
}
