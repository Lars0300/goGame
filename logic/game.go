package logic

import (
	"chatChannel/protocol"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"
)

type Game struct {
	currentHolder    *Player
	players          []*Player
	alivePlayers     map[*Player]struct{}
	broadcastChannel chan []byte
	startingPlayers  map[*Player]struct{}
	hash             string
	hasStarted bool
}

func (g *Game) UpdateBroadcast(data []byte) {
	g.broadcastChannel <- data
}

func (g *Game) AlreadyVotedToStart(player *Player) (bool){
	_, inMap := g.startingPlayers[player]
	return inMap
}

func (g*Game) RemovePlayer(player *Player){
	delete(g.alivePlayers, player)
	delete(g.startingPlayers, player)
	for i, gamePlayer := range(g.players){
		if gamePlayer == player{
			g.players = append(g.players[:i], g.players[i+1:]...)
			break
		}
	}
	if len(g.players) == 0{
		delete(HashToGame, g.hash)
		delete(OpenGames, g)
		return
	}
	
	if len(g.players) == 1{
		g.endGame()

	}
	if g.currentHolder == player{
		g.currentHolder = nil
		for key, _ := range(g.alivePlayers){
			g.currentHolder = key
		}

	}
}
func (g *Game) broadcast() {
	log.Printf("Starting Broadcast for Game %s", g.hash)
	for msg := range g.broadcastChannel {
		for _, player := range g.players {
			player.Outgoing <- msg
		}
	}
	log.Printf("Broadcast for Game: %s has stopped", g.hash)
}

func (g *Game) checkForStart() {
	log.Printf("Checking start for Game %s", g.hash)
	if len(g.startingPlayers) > (len(g.players) / 2) {
		g.startGame()
	}
}

func (g *Game) startGame() {
	log.Printf("Starting Game %s", g.hash)
	delete(OpenGames, g)
	g.hasStarted = true
	g.UpdateBroadcast(protocol.BuildGameUpdate(protocol.StartGame, "Game", fmt.Sprintf("!!Game is starting!!\n Player %s is holding the bomb.", g.currentHolder.username)))
	go g.gameLoop()
}

func (g *Game) gameLoop() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		delta := r.Int63n(int64(maxDuration - minDuration + 1))
		randomDuration := minDuration + time.Duration(delta)

		timer := time.NewTimer(randomDuration)
		<-timer.C
		g.killPlayer()

		if g.isOver() {
			g.endGame()
			return
		}
	}
}

func (g *Game) killPlayer() {
	delete(g.alivePlayers, g.currentHolder)
	var broadcastMessage string = fmt.Sprintf("Player %s got eliminated", g.currentHolder.username)
	g.broadcastChannel <- protocol.BuildGameUpdate(protocol.KillPlayer, "Game", broadcastMessage)
	for player := range g.alivePlayers {
		g.currentHolder = player
		return
	}
}

func (g *Game) GetHash() string{
	return g.hash
}

func (g *Game) endGame() {
	for _, player := range g.players {
		g.alivePlayers[player] = struct{}{}
	}
	var broadcastMessage string = fmt.Sprintf("Game Over! \nThe Winner is %s", g.currentHolder.username)
	g.broadcastChannel <- protocol.BuildGameUpdate(protocol.EndGame, "Game", broadcastMessage)
	g.startingPlayers = make(map[*Player]struct{})
	OpenGames[g] = struct{}{}
	g.hasStarted = false
}
func (g *Game) Pass(toPlayer *Player) error {
	log.Printf("Passing now from %s to %s", g.currentHolder.username, toPlayer.username)
	g.currentHolder = toPlayer
	return nil
}

func (g *Game) isOver() bool {
	return len(g.alivePlayers) == 1
}

func (game *Game) addPlayer(player *Player) {
	game.alivePlayers[player] = struct{}{}
	game.players = append(game.players, player)
	player.currentGame = game
}

func (g *Game) CheckUsername(username string) bool {
	for _, player := range g.players {
		if player.playerID == username {
			return false
		}
	}
	return true
}

func (g *Game) GetCurrentHolder() *Player {
	return g.currentHolder
}

func (g *Game) GetPlayerForUsername(username string) (*Player, error) {
	for player := range g.alivePlayers {
		if player.username == username {
			return player, nil
		}
	}
	return nil, fmt.Errorf("username %s not found or not alive", username)
}

func (g *Game) ToString() string {
	var playersList []string
	for _, p := range g.players {
		playersList = append(playersList, p.playerID)
	}

	var aliveList []string
	for p := range g.alivePlayers {
		aliveList = append(aliveList, p.playerID)
	}

	var startList []string
	for p := range g.startingPlayers {
		startList = append(startList, p.playerID)
	}

	return fmt.Sprintf(
		"Game{\n hash: %s,\n  hasStarted: %v,\n  currentHolder: %s,\n  players: [%s],\n  alivePlayers: [%s],\n  startingPlayers: [%s]\n}",
		g.hash,
		g.hasStarted,
		g.currentHolder.playerID,
		strings.Join(playersList, ", "),
		strings.Join(aliveList, ", "),
		strings.Join(startList, ", "),
	)
}

func (g *Game) HasStarted() (bool){
	return g.hasStarted
}

func (g *Game) OnlyOnePlayer() (bool){
	return len(g.alivePlayers) == 1
}