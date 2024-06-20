package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"ts-game/engine"
)

func main() {
	server := engine.New()
	if len(os.Args) == 1 {
		fmt.Println("No mode provided, defaulting to local")
		server.Connect(os.Stdin, os.Stdout, func() { fmt.Println("Goodbye!") })
	} else if os.Args[1] == "local" {
		server.Connect(os.Stdin, os.Stdout, func() { fmt.Println("Goodbye!") })
	} else if os.Args[1] == "telnet" {
		log.Println("Listening for TCP connections...")
		listener, err := net.Listen("tcp", ":8081")
		if err != nil {
			log.Fatal(err.Error())
		}
		go func() {
			for {
				conn, err := listener.Accept()
				if err != nil {
					log.Printf("Error accepting connection: %v", err.Error())
					continue
				}
				server.Connect(conn, conn, func() { conn.Close() })
			}
		}()
	} else {
		log.Fatalf("Unknown mode: %v", os.Args[1])
	}
	server.Start()
}
