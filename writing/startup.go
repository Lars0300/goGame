package writing

import (
	"bufio"
	"chatChannel/protocol"
	"fmt"
	"os"
	"strings"
	"time"
)
const typingSpeed time.Duration = 50
func WriteStartup() string{
	var d *Dialog = GlobalDialog
	writeOut(protocol.Intro, "Game", d.Start.Intro, time.Now().Unix(), time.Duration(typingSpeed * time.Millisecond))
	writeOut(protocol.Intro, "Game", d.Start.EnterName, time.Now().Unix(), time.Duration(typingSpeed * time.Millisecond))
	nameReader := bufio.NewReader(os.Stdin)
	name, _ := nameReader.ReadString('\n')
	name = strings.TrimSpace(name)


	writeOut(protocol.Intro, "Game", fmt.Sprintf(d.Start.Greeting, name), time.Now().Unix(), time.Duration(typingSpeed * time.Millisecond))
	writeOut(protocol.Intro, "Game", d.Start.Info, time.Now().Unix(), time.Duration(typingSpeed * time.Millisecond))
	writeOut(protocol.Intro, "Game", d.Start.Help, time.Now().Unix(), time.Duration(typingSpeed * time.Millisecond))
	return name
	
}