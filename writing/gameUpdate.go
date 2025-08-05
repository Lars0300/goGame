package writing

import (
	"fmt"
	"chatChannel/protocol"
	"time"
)

func WritePong(pong *protocol.Pong){
	now := time.Now()
	msg := fmt.Sprintf("recieved after %d ms", now.UnixMilli() - pong.PingTimestamp)
	writeOut(protocol.PongMsg, "Pong", msg, now.Unix(), time.Duration(0))
	
}

func WriteGameUpdate(update *protocol.GameUpdate){

	writeOut(update.UpdateType, update.From, update.Msg, update.Time, time.Duration(0))

}

func WriteChangeConfirm(name string){
	writeOut(protocol.NameChange, "Game", fmt.Sprintf(GlobalDialog.Game.GameHost.ChangeName, name), time.Now().Unix(), time.Duration(0))
}