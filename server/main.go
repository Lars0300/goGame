package main

import (
	"chatChannel/logic"
	"chatChannel/protocol"
	"fmt"
	"log"
	"net"
)

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "8080"
	SERVER_TYPE = "tcp"
)

func main() {
	fmt.Println("Server running...")

	server, err := net.Listen(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
	if err != nil {
		log.Fatalf("Error listening: %v", err)
		return
	}

	defer server.Close()

	fmt.Println("Listening on " + SERVER_HOST + ":" + SERVER_PORT)
	fmt.Println("Waiting for client...")
	go handleConsoleInput()
	for {
		con, err := server.Accept()

		if err != nil {
			log.Fatalf("Error accepting: %v", err)
			return
		}

		fmt.Println("Client connected from: " + con.RemoteAddr().String())

		go handleConnection(con)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)

	var player *logic.Player = logic.CreatePlayer(conn)
	go player.WriteLoop()
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("client " + conn.RemoteAddr().String() + " disconnected.")
			} else {
				fmt.Println("Error reading:", err.Error())
				log.Fatalf("Error reading: %v", err)
			}
			return
		}

		handleInput(buffer[:n], player, conn.RemoteAddr().String())
	}
}

func handleInput(data []byte, player *logic.Player, connID string) {
	decodedMessage, err := protocol.UnmarshalMessage(data)
	if err != nil {
		log.Printf("Error decoding message: %v", err)
		player.Outgoing <- protocol.BuildErrorMessage("Invalid json data", err.Error())
		return
	}
	switch decodedMessage.Type {
	case "join":
		handleJoin(data, player, connID)
	case "start":
		handleStart(data, player, connID)
	case "pass":
		handlePass(data, player, connID)
	case "chat":
		handleChat(data, player, connID)
	case "ping":
		handlePing(data, player, connID)
	default:
		player.Outgoing <- protocol.BuildErrorMessage("Invalid Message Request", "Type not accepted")
	}
}
