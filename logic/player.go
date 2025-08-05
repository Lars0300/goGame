package logic

import (
	"chatChannel/protocol"
	"chatChannel/writing"
	"fmt"
	"log"
	"net"
)

type Player struct {
	playerID    string
	currentGame *Game
	Outgoing    chan []byte
	connection  net.Conn
	username    string
}

func CreatePlayer(conn net.Conn) *Player {
	var player *Player = &Player{
		playerID:    generatePlayerID(),
		currentGame: nil,
		connection:  conn,
		Outgoing:    make(chan []byte, 1024),
		username:    "",
	}

	return player
}

func (player *Player) GetCurrentGame() *Game {
	return player.currentGame
}

func (player *Player) GetPlayerID() string {
	return player.playerID
}

func (player *Player) JoinRandom() {
	for game := range OpenGames {
		err := player.JoinHash(game.hash)
		if err == nil{
			return
		}
	}
	player.createGame("")
}

func (player *Player) JoinHash(hash string) error {
	game, ok := HashToGame[hash]
	if !ok {
		log.Printf("Game for the hash %s not found, building game for it now\n", hash)
		player.currentGame = player.createGame(hash)
		return nil
	}
	if game.hasStarted {
		return fmt.Errorf("Game has already started")
	}
	for _, currPlayer := range(game.players){
		if player.username == currPlayer.username{
			player.Outgoing <- protocol.BuildGameUpdate(protocol.JoinFailed, "Game", fmt.Sprintf(writing.GlobalDialog.Game.GameHost.CantJoin, hash, hash));
			return fmt.Errorf("username has already been taken for player %s and game %s", player.playerID, game.hash)
		}
	}
	game.addPlayer(player)
	return nil
}

func (player *Player) createGame(suggestedHash string) *Game {
	log.Printf("Creating game for suggested Hash %s", suggestedHash)
	var hash string
	if _, inMap := HashToGame[suggestedHash]; inMap {
		log.Printf("Player %s couldn't join game with suggested hash: %s. Creates now new random game instead", player.playerID, suggestedHash)
	}
	if isValidHash(suggestedHash) {
		hash = suggestedHash
	} else {
		for {
			hash = generateHash()
			if _, inMap := HashToGame[hash]; !inMap {
				break
			}
		}
	}

	var game *Game = &Game{
		currentHolder:    player,
		players:          make([]*Player, 0),
		alivePlayers:     make(map[*Player]struct{}),
		broadcastChannel: make(chan []byte, 1024),
		startingPlayers:  map[*Player]struct{}{},
		hash:             hash,
		hasStarted:       false,
	}

	game.addPlayer(player)
	go game.broadcast()

	game.UpdateBroadcast(protocol.BuildGameUpdate(protocol.CreateGame, "Game", fmt.Sprintf("Game has been created. Hash: %s", hash)))
	HashToGame[hash] = game
	OpenGames[game] = struct{}{}
	return game
}

func (player *Player) VoteStart() {
	log.Printf("Player %s votes for start", player.username)
	player.currentGame.startingPlayers[player] = struct{}{}
	player.currentGame.checkForStart()
}

func (player *Player) WriteLoop() {
	log.Printf("Starting writeLoop for player %s", player.playerID)
	for msg := range player.Outgoing {
		player.connection.Write(append(msg, '\n'))
	}
	log.Printf("writeLoop closes for player %s", player.playerID)
}

func (player *Player) Chat(content string) {
	var chatType protocol.UpdateStatus
	var currentGame *Game = player.currentGame
	
	if currentGame == nil {
		player.Outgoing <- protocol.BuildGameUpdate(protocol.ChatDead, player.username, content)
		return
	}
	if _, inSet := currentGame.alivePlayers[player]; !inSet {
		chatType = protocol.ChatDead
	} else {
		if player.username == player.currentGame.currentHolder.username {
			chatType = protocol.ChatBomb
		} else {
			chatType = protocol.ChatAlive
		}
	}
	player.currentGame.broadcastChannel <- protocol.BuildGameUpdate(chatType, player.username, content)
}

func (player *Player) GetUsername() string {
	return player.username
}
func (player *Player) SetUsername(username string){
	player.username = username
}