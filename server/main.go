package main

import (
	"chatChannel/logic"
	"chatChannel/protocol"
	"encoding/json"
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
	fmt.Println("Server running....")

	server, err := net.Listen(SERVER_TYPE, SERVER_HOST + ":" + SERVER_PORT)
	if err != nil {
		log.Fatalf("Error listening: %v", err)
		return 
	}

	defer server.Close()

	fmt.Println("Listening on " + SERVER_HOST + ":" + SERVER_PORT)
	fmt.Println("Waiting for client...")

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

func handleConnection(conn net.Conn){
	defer conn.Close()

	buffer := make([]byte, 1024)

	var player *logic.Player = logic.CreatePlayer(conn)
	go player.WriteLoop()
	for {
		_, err := conn.Read(buffer)
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("client " + conn.RemoteAddr().String() + " disconnected.")
			} else {
				fmt.Println("Error reading:", err.Error())
				log.Fatalf("Error reading: %v", err)
			}
			return
		}

		handleInput(buffer, player, conn.RemoteAddr().Network())
	}
}

func handleInput(data []byte, player *logic.Player, connID string){
	decodedMessage, err := protocol.UnmarshalMessage(data)

	if err != nil{
		log.Printf("Error decoding message: %v", err)
		player.Outgoing <- protocol.BuildErrorMessage("Invalid json data", err.Error())
		return 
	}
	switch decodedMessage.Type{
	case "join":
		var joinRequest protocol.JoinRequest
		err := json.Unmarshal(data, &joinRequest)
		if err != nil{
			log.Printf("Error for Client %s while unmarshaling JoinRequest: %v", connID, err)
			player.Outgoing <- protocol.BuildErrorMessage("Invalid Join Request", "Malformed message")
			return
		}

		if joinRequest.Random {
			log.Printf("Client %s joins Random Game", connID)
			player.JoinRandom()
		} else {
			log.Printf("Client %s wants to join game", joinRequest.Hash)
			player.JoinHash(joinRequest.Hash)
		}
	case "start":
		player.VoteStart()
	case "pass":
		var passBomb protocol.PassBomb
		err := json.Unmarshal(data, &passBomb)
		if err != nil{
			log.Printf("Error for Client %s while unmarshaling JoinRequest: %v", connID, err)
			player.Outgoing <-protocol.BuildErrorMessage("Invalid Join Request", "Malformed message")
			return
		}
		if player != player.GetCurrentGame().GetCurrentHolder(){
			return
		}
		err = player.GetCurrentGame().Pass(passBomb.Recipient)
		if err != nil{
			player.Outgoing <- protocol.BuildErrorMessage("Invalid Recipient", err.Error())
		}
	case "chat":
		var chatMsg protocol.ClientChatMessage
		err := json.Unmarshal(data, &chatMsg) 
		if err != nil{
			log.Printf("Error for Client %s while unmarshaling JoinRequest: %v", connID, err)
			player.Outgoing <- protocol.BuildErrorMessage("Invalid Join Request", "Malformed message")
			return
		}
		player.Chat(chatMsg.Content)
	case "ping":
		var ping protocol.Ping
		err := json.Unmarshal(data, &ping)
		
		if err != nil{
			log.Printf("Error for Client %s while unmarshaling JoinRequest: %v", connID, err)
			player.Outgoing <- protocol.BuildErrorMessage("Invalid Join Request", "Malformed message")
			return
		}
		player.Outgoing <- protocol.BuildPong()
		
	default:
		player.Outgoing <- protocol.BuildErrorMessage("Invalid Message Request", "Type not accepted")
	}
}