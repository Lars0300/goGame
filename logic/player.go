package logic

import (
	"chatChannel/protocol"
	"net"
)

type Player struct {
	username    string
	currentGame *Game
	Outgoing    chan []byte
	connection  net.Conn
}

func CreatePlayer(conn net.Conn) *Player {
	var player *Player = &Player{
		username:    generateUsername(),
		currentGame: nil,
		connection:  conn,
		Outgoing:    make(chan []byte, 1024),
	}
	return player
}

func (player *Player) GetCurrentGame() *Game {
	return player.currentGame
}
func (player *Player) ChangeUsername(username string) {
	player.username = username
}

func (player *Player) JoinRandom() {
	for game := range openGames {
		game.addPlayer(player)
		return
	}
	player.createGame()
}

func (player *Player) JoinHash(hash string) {
	game, ok := hashToGame[hash]
	if !ok {
		player.currentGame = player.createGame()
		return
	}
	game.players = append(game.players, player)
}

func (player *Player) createGame() *Game {

	
	var hash string
	for {
		hash = generateHash()
		if _, inMap := hashToGame[hash]; !inMap {
			break
		}
	}

	

	var game *Game = &Game{
		currentHolder:    player,
		players:          make([]*Player, 0),
		alivePlayers:     make(map[*Player]struct{}),
		broadcastChannel: make(chan []byte, 1024),
		startingPlayers: map[*Player]struct{}{},
	}

	game.addPlayer(player)
	go game.broadcast()
	hashToGame[hash] = game
	openGames[game] = struct{}{}
	return game
}

func (player *Player) VoteStart(){
	player.currentGame.startingPlayers[player] = struct{}{}
	player.currentGame.checkForStart()
}

func (player *Player) WriteLoop() {
	for msg := range player.Outgoing {
		player.connection.Write(append(msg, '\n'))
	}
}

func (player *Player) Chat(content string){
	player.currentGame.broadcastChannel <- protocol.BuildGameUpdate("chat", player.username, content)
}
