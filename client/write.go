package main

import (
	"chatChannel/protocol"
	"fmt"
	"time"
)
var (
	colorMap = map[protocol.UpdateStatus]string{
		protocol.ChatAlive: "\033[32m",
		protocol.ChatDead: "\033[37m",
		protocol.ChatBomb: "\033[37m",
		protocol.JoinGame: "\033[36m",
		protocol.StartGame: "\033[35m", 
		protocol.EndGame: "\033[35m",
		protocol.CreateGame: "\033[33m",
		protocol.KillPlayer: "\033[",
		protocol.Pass: "\033[31m",
		protocol.VoteStart: "\033[31m",

	}
)
func writePong(pong *protocol.Pong){
	fmt.Printf("Game:> Pong recieved after %d ms\n", time.Now().UnixMilli() - pong.PingTimestamp)
}

func writeGameUpdate(update *protocol.GameUpdate){
	var timestamp int64 = update.Time
	var time time.Time = time.Unix(timestamp, 0).In(time.Local)

	var clockString string = time.Format("15:04")

	fmt.Printf("%s %s%s\033[0m:> %s\n", clockString, colorMap[update.UpdateType], update.From, update.Msg)
	
}