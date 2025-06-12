package protocol

import (
	"encoding/json"
	"log"
	"time"
)

type Message struct {
	Type string `json:"type"`
}

// client to server 
type JoinRequest struct {
	Message
	Username string `json:"username"`
	Hash string `json:"hash"`
	Random bool `json:"random"`
}

type PassBomb struct {
	Message
	Recipient string `json:"Recipient"`
}

type ClientChatMessage struct {
	Message
	Content string `json:"Content"`
}

type Ping struct {
	Message
	Timestamp int64 `json:"timestamp"`
}

type StartMessage struct{
	Message
}

// Server to Client Message
type GameUpdate struct {
	Message
	Type string `json:"type"`
	From string `json:"from"`
	Msg string `json:"message"`
	Time int64 `json:"time"`
}

type ErrorMessage struct {
	Message
	Error string
	Details string
}
type Pong struct {
	Message
	Timestamp int64 `json:"timestamp"`
}

func MarshallMessage(msg interface{}) ([]byte, error) {
	return json.Marshal(msg)
}

func UnmarshalMessage(data []byte) (*Message, error){
	var m Message 
	err := json.Unmarshal(data, &m)
	if err != nil{
		return nil, err
	}
	return &m, nil
}

func BuildErrorMessage(err string, details string) ([]byte) {
	errorMessage := &ErrorMessage{
		Message: Message{Type: "error"},
		Error: err,
		Details: details,
	}
	errData, errr := json.Marshal(errorMessage)
	if errr != nil{
		log.Printf("Error building error message with err %s and details %s. Error: %s", err, details, errr)
		return []byte{}
	}
	return append(errData, '\n')
}
func BuildGameUpdate(updateType string, from string, msg string)([] byte){
	gameUpdate := &GameUpdate{
		Message: Message{Type: "update"},
		Type: updateType,
		From: from,
		Msg: msg,
	}

	data, err := json.Marshal(gameUpdate)
	if err != nil{
		log.Printf("Eror building gameUpdate with updateType %s, From %s, msg %s. Error: %s", updateType, from, msg, err)
		return []byte{}
	}

	return append(data, '\n')
}

func BuildPong()([]byte){
	pong := Pong{
			Message: Message{Type: "pong"},
			Timestamp: time.Now().Unix(),
	}

	data, err := json.Marshal(pong)
	if err != nil{
		log.Printf("Eror building pong. Error: %s", err)
		return []byte{}
	}

	return append(data, '\n')
}
