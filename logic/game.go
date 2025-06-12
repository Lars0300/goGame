package logic

import (
	"chatChannel/protocol"
	"fmt"
	"log"
	"math/rand"
	"time"
)

type Game struct {
	currentHolder    *Player
	players          []*Player
	alivePlayers     map[*Player]struct{}
	broadcastChannel chan []byte
	startingPlayers map[*Player]struct{}
}


func (g *Game) UpdateBroadcast(gameUpdate *protocol.GameUpdate) {
	var data []byte
	data, err := protocol.MarshallMessage(gameUpdate)
	if err != nil {
		return
	}
	g.broadcastChannel <- data
}

func (g *Game) broadcast() {
	for msg := range g.broadcastChannel {
		for _, player := range g.players {
			player.Outgoing <- msg
		}
	}
}

func (g *Game) checkForStart(){
	if len(g.startingPlayers) > len(g.players){
		g.startGame()
	}
}

func(g *Game) startGame(){
	delete(openGames, g)
	go g.gameLoop()
}

func (g *Game) gameLoop(){
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		delta := r.Int63n(int64(maxDuration - minDuration + 1))
		randomDuration := minDuration + time.Duration(delta)

		timer := time.NewTimer(randomDuration)
		<- timer.C
		g.killPlayer()

		if g.isOver(){
			g.endGame()
			return
		}
	}
}

func (g *Game) killPlayer(){
	delete(g.alivePlayers, g.currentHolder)
	var broadcastMessage string = fmt.Sprintf("Player %s got eliminated", g.currentHolder.username)
	g.broadcastChannel <- protocol.BuildGameUpdate("kill", "game", broadcastMessage)
	for player := range g.alivePlayers{
		g.currentHolder = player
		return
	}
}

func (g *Game) endGame(){
	for _, player := range g.players{
		g.alivePlayers[player] = struct{}{}
	}
	var broadcastMessage string = fmt.Sprintf("Game Over! \n The Winner is %s", g.currentHolder.username)
	g.broadcastChannel <- protocol.BuildGameUpdate("end", "game", broadcastMessage)
	g.startingPlayers = make(map[*Player]struct{})
	openGames[g] = struct{}{}
}
func (g *Game) Pass(toUsername string) error {
	log.Printf("Passing now from %s to %s", g.currentHolder.username, toUsername)
	toPlayer, err := g.getPlayerForUsername(toUsername)

	if err != nil {
		return err
	}
	broadcastMessage := fmt.Sprintf("%s holds the bomb now", toUsername)
	g.broadcastChannel <- protocol.BuildGameUpdate("pass", "game", broadcastMessage)
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
		if player.username == username {
			return false
		}
	}
	return true
}

func(g *Game) GetCurrentHolder() *Player{
	return g.currentHolder
}

func (g *Game) getPlayerForUsername(username string) (*Player, error) {
	for player := range g.alivePlayers {
		if player.username == username {
			return player, nil
		}
	}
	return nil, fmt.Errorf("ussername %s not found or not alive", username)
}
