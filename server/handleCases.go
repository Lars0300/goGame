package main

import (
	"chatChannel/logic"
	"chatChannel/protocol"
	"encoding/json"
	"fmt"
	"log"
)
func handleJoin(data []byte, player *logic.Player, connID string) {
	var joinRequest protocol.JoinRequest
	err := json.Unmarshal(data, &joinRequest)
	if err != nil {
		log.Printf("Error for Client %s while unmarshaling JoinRequest: %v", connID, err)
		player.Outgoing <- protocol.BuildErrorMessage("Invalid Join Request", "Malformed message")
		return
	}
	player.SetUsername(joinRequest.Username)
	beforeGame := player.GetCurrentGame()
	if joinRequest.Random {
		log.Printf("Client %s joins Random Game", connID)
		player.JoinRandom()
	} else {
		log.Printf("Client %s wants to join game %s", connID, joinRequest.Hash)
		err = player.JoinHash(joinRequest.Hash)
	}
	
	if err != nil {
		log.Printf("Error for Client %s while unmarshaling JoinRequest: %v", connID, err)
		player.Outgoing <- protocol.BuildErrorMessage("Coldn't join game", err.Error())
		return
	}
	fmt.Printf("beforeGame")
	if beforeGame != nil{
		beforeGame.RemovePlayer(player)
	}
	
	var currentGame *logic.Game = player.GetCurrentGame()
	currentGame.UpdateBroadcast(protocol.BuildGameUpdate(protocol.JoinGame, "Game", fmt.Sprintf("Player %s joins the game", player.GetUsername())))

}

func handleStart(data []byte, player *logic.Player, connID string) {
	log.Printf("Player %s votes to start the game %s", player.GetPlayerID(), player.GetCurrentGame().GetHash())
	var currentGame *logic.Game = player.GetCurrentGame()
	if currentGame.HasStarted(){
		log.Println("Game has already started, returning")
		return
	}
	if currentGame.AlreadyVotedToStart(player){
		log.Println("Player has already voted, returning")
		return
	}
	if currentGame.OnlyOnePlayer(){
		log.Println("Can't play alone mate")
		return
	}
	currentGame.UpdateBroadcast(protocol.BuildGameUpdate(protocol.VoteStart, "Game", fmt.Sprintf("Player %s votes to start the game", player.GetUsername())))
	player.VoteStart()
}

func handlePass(data []byte, player *logic.Player, connID string) {
	var passBomb protocol.PassBomb
	var currentGame *logic.Game = player.GetCurrentGame()
	err := json.Unmarshal(data, &passBomb)
	if err != nil {
		log.Printf("Error for Client %s while unmarshaling JoinRequest: %v", connID, err)
		player.Outgoing <- protocol.BuildErrorMessage("Invalid Join Request", "Malformed message")
		return
	}
	if player != currentGame.GetCurrentHolder() {
		return
	}
	var recipient *logic.Player
	recipient, err = currentGame.GetPlayerForUsername(passBomb.Recipient)

	if err != nil{
		player.Outgoing <- protocol.BuildErrorMessage("Target Player doesn't exist", err.Error())
	}
	err = currentGame.Pass(recipient)
	if err != nil {
		player.Outgoing <- protocol.BuildErrorMessage("Target Player Not Alive", err.Error())
	}

	currentGame.UpdateBroadcast(protocol.BuildGameUpdate(protocol.Pass, "Game", fmt.Sprintf("%s passed the bomb to %s", player.GetUsername(), recipient.GetUsername())))
	recipient.Outgoing <- protocol.BuildGameUpdate(protocol.Pass, "Game", "You have now the bomb, type !pass <username> to pass it to someone else")
}

func handleChat(data []byte, player *logic.Player, connID string) {
	var chatMsg protocol.ClientChatMessage
	err := json.Unmarshal(data, &chatMsg)
	if err != nil {
		log.Printf("Error for Client %s while unmarshaling JoinRequest: %v", connID, err)
		player.Outgoing <- protocol.BuildErrorMessage("Invalid Join Request", "Malformed message")
		return
	}
	player.Chat(chatMsg.Content)
}

func handlePing(data []byte, player *logic.Player, connID string) {
	var ping protocol.Ping
	err := json.Unmarshal(data, &ping)

	if err != nil {
		log.Printf("Error for Client %s while unmarshaling JoinRequest: %v", connID, err)
		player.Outgoing <- protocol.BuildErrorMessage("Invalid Join Request", "Malformed message")
		return
	}
	player.Outgoing <- protocol.BuildPong(ping.Timestamp)
}

func handleChange(data []byte, player *logic.Player, connID string){
	var change protocol.ChangeName
	

	err := json.Unmarshal(data, &change)

	if err != nil {
		log.Printf("Error for Client %s while unmarshaling JoinRequest: %v", connID, err)
		player.Outgoing <- protocol.BuildErrorMessage("Invalid Change Request", "Malformed message")
		return
	}
	player.SetUsername(change.Name)
	fmt.Println(change)
	player.Outgoing <- protocol.BuildGameUpdate(protocol.NameChange, "Game", fmt.Sprintf("Your name has changed to %s.", change.Name))
	player.Outgoing <- protocol.BuildChangeConfirm(change.Name)
}