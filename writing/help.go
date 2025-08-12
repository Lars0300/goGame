package writing

import (
	"chatChannel/protocol"
	"time"
)

func WriteHelp() {
	var d *Dialog = GlobalDialog
	writeOut(protocol.Help, "Game", d.Help.HelpMenu, time.Now().Unix(), time.Duration(typingSpeed*time.Millisecond))
}