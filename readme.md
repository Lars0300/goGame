# Chat Channel Game 
## Organization
### How to run it
There are two Files for the game:
- `bin/server.exe`needs to be run first and initializes the server. It outputs debugging for the server.
- `bin/client.exe`has the game interface. Each player needs to run their own `client.exe`
### Where is it hosted
The server is hosted on `localhost/8080`
## How To Play
### Game Rules
Multiple people can join a game, where one 'hot potato' is given to a random player. Players can give their potato to the other players. The exploding time is set to random. If the potato does explode the player holding it is out until the end of the game
The winner is the last person standing.
### Commands
Type commands into the console to play. A comprehensive command list can be seen by typing `!help` or down below 
## Developer Notes
### Changing Dialog
If you want to change the messages being send, edit the dialog/dialog.json file.
If you want to add new messages, just edit the dialogStructure struct. The programm automatically updates the json file in `dialog`. 
Note that information can be lost if you update the json reference inside the struct.