package main

import (
	"bufio"
	"chatChannel/protocol"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "8080"
	SERVER_TYPE = "tcp"
)

type Profile struct{
	GameActive bool 
	HasBomb bool 
	currentHash string
	name string 
}
func main() {
	conn, err := net.Dial(SERVER_TYPE, SERVER_HOST + ":" + SERVER_PORT)

	if err != nil {
		log.Printf("Error connecting: %v", err)
		fmt.Println("Error connecting:", err)
		return 
	}

	defer conn.Close()

	
	fmt.Print("Enter your name:")
	nameReader := bufio.NewReader(os.Stdin)
	name, _ := nameReader.ReadString('\n')
	name = strings.TrimSpace(name)

	profile := Profile{
		GameActive: false,
		HasBomb: false,
		currentHash: "",
		name: strings.TrimSpace(name),
	}
	go listenToServer(conn)

	input := bufio.NewScanner(os.Stdin)
	


	for input.Scan() {
		text := input.Text()
		if strings.TrimSpace(text) == ""{
			continue
		}
		var encodedMsg []byte 
		encodedMsg, err := handleInput(text, &profile)
		
		if err != nil{
			continue
		}
		

		_, err = conn.Write(append(encodedMsg, '\n'))
        if err != nil {
            log.Printf("Error sending message: %v", encodedMsg)
        }
	}
}

func listenToServer(conn net.Conn) {
	buffer := make([]byte, 1024)

	for{
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

		handleWrite(buffer[:n])
	}
}


func handleInput(text string, profile *Profile) ([]byte, error){
	var encodedMsg []byte
	var err error
	if text[0] == '!'{
			var key, message string
			var found bool
			key, message, found = strings.Cut(text," ") 
			
			switch strings.TrimPrefix(key, "!"){
			case "pass":
				if !profile.HasBomb{
					fmt.Println("Can't pass bomb if you don't hold it mate")
					return nil, fmt.Errorf("can't pass the bomb")
				}
				if !profile.GameActive{
					fmt.Println("Game's over mate")
					return nil, fmt.Errorf("games over")
				}
				if !found{
					message = ""
				}
				passMsg := protocol.PassBomb{
					Message: protocol.Message{Type: "pass_bomb"},
					Recipient: message,
				}
				encodedMsg, err = protocol.MarshallMessage(passMsg)

				if err != nil {
					log.Printf("Error encoding pass message: %v", passMsg)
					return nil, err
				}
			case "ping":
				pingMsg := protocol.Ping{
					Message: protocol.Message{Type: "ping"},
					Timestamp: time.Now().UnixMilli(),
				}
				encodedMsg, err = protocol.MarshallMessage(pingMsg)

				if err != nil {
					log.Printf("Error encoding ping message: %v", pingMsg)
					return nil, err
				}
			case "join":
				var random bool = false
				if !found{
					message = "random"
				}
				if message == "random"{
					random = true
				} else {
					profile.currentHash = message
				}
				joinMsg := protocol.JoinRequest{
					Message: protocol.Message{Type: "join"},
					Random: random,
					Hash: profile.currentHash,
				}

				encodedMsg, err = protocol.MarshallMessage(joinMsg)

				if err != nil {
					log.Printf("Error encoding join message: %v", joinMsg)
					return nil, err
				}
			case "start":
				startMsg := protocol.StartMessage{
					Message: protocol.Message{Type: "start"},
				}

				encodedMsg, err = protocol.MarshallMessage(startMsg)

				if err != nil {
					log.Printf("Error encoding start message: %v", startMsg)
					return nil, err
				}
			default: 
				log.Printf("Error using command : %v", text)
				return nil, fmt.Errorf("error using command : %v", text)
			}
		} else {
			chatMsg := protocol.ClientChatMessage{
				Message: protocol.Message{Type: "chat"},
				Content: text,
			}

			encodedMsg, err = protocol.MarshallMessage(chatMsg)

			if err != nil {
					log.Printf("Error encoding chat message: %v", chatMsg)
					return nil, err
				}
		}
		return encodedMsg, nil
}

func handleWrite(data []byte){
	decodedMessage, err := protocol.UnmarshalMessage(data)

	if err != nil{
		log.Printf("Error decoding message: %v", err)
		return 
	}
	switch decodedMessage.Type{
	case "pong":
		var pong protocol.Pong

		err = json.Unmarshal(data, &pong)

		if err != nil{
			log.Printf("Error unmarshalling message:  %s", err)
		}

		writePong(&pong)
	case "update":
		var update protocol.GameUpdate

		err = json.Unmarshal(data, &update)

		if err != nil{
			log.Printf("Error unmarshalling message:  %s", err)
		}
		writeGameUpdate(&update)
	}
		

}