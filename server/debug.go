package main

import (
	"bufio"
	"chatChannel/logic"
	"fmt"
	"log"
	"os"
	"strings"
)
var commandMap = map[string]string{
	"!help": "Opens help menu",
	"!games" : "Show open games",
}
func handleConsoleInput(){
	scanner := bufio.NewScanner(os.Stdin)
	log.Println("Debuger ready. Type !help to see commands")
	for scanner.Scan(){
		text := scanner.Text()
		if text[0] != '!'{
			fmt.Printf("Command has to start with \"!\". Type !help to see all commands")
			continue;
		}

		switch strings.Trim(text, "!"){
		case "help":
			for command, description := range commandMap{
				fmt.Printf("%s \t %s\n", command, description)
			}
		case "games":
			for game, _ := range logic.OpenGames{
				fmt.Println(game.ToString())
			}
		}
	}
}