package main

import (
	"flag"
	"log"
	"log/slog"
	"net"
	"os"
	"ts-game/engine"
)

func main() {
	verboseLoggingFlag := flag.Bool("v", false, "enables verbose logging output")
	flag.Parse()
	programLevel := new(slog.LevelVar)
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: programLevel})
	slog.SetDefault(slog.New(handler))
	if *verboseLoggingFlag {
		programLevel.Set(slog.LevelDebug)
	}
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
