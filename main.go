package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"net"
	"os"
	"ts-game/engine"
	"ts-game/persistence"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Info("No .env file found")
	}

	ctx := context.TODO()
	pgPool, err := persistence.New(ctx)
	if err != nil {
		log.Fatal(err)
	}
	err = pgPool.Ping(ctx)
	if err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

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
