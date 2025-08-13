package writing

import (
	"chatChannel/protocol"
	"fmt"
	"strings"
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
		protocol.KillPlayer: "\033[35m",
		protocol.Pass: "\033[31m",
		protocol.VoteStart: "\033[31m",
		protocol.Intro: "\033[31m",
		protocol.PongMsg: "\033[32m",
		protocol.JoinFailed: "\033[31m",
		protocol.Help: "\033[33m",
	}
)

func writeOut(updateType protocol.UpdateStatus, from, msg string, timestamp int64, delay time.Duration){
	var timeimg time.Time = time.Unix(timestamp, 0).In(time.Local)

	var clockString string = timeimg.Format("15:04")
	
	var lineArray []string = strings.Split(msg, "\n")
	for _, line := range(lineArray){
		if line == ""{
			continue
		}
		var msgArr []string = strings.Split(line, " ")
		fmt.Printf("%s %s%s\033[0m:> ", clockString, colorMap[updateType], from)
		for _, word := range(msgArr){
			fmt.Print(word)
			time.Sleep(delay)
			fmt.Print(" ")
		}
		fmt.Print("\n")
	}
}