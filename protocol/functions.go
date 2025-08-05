package protocol

import (
	"encoding/json"
	"log"
	"time"
)

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
func BuildGameUpdate(updateType UpdateStatus, from string, msg string)([]byte){
	gameUpdate := &GameUpdate{
		Message: Message{Type: "update"},
		UpdateType: updateType,
		From: from,
		Msg: msg,
		Time: time.Now().Unix(),
	}

	data, err := json.Marshal(gameUpdate)
	if err != nil{
		log.Printf("Eror building gameUpdate with updateType %v, From %s, msg %s. Error: %s", updateType, from, msg, err)
		return []byte{}
	}

	return append(data, '\n')
}

func BuildChangeConfirm(name string) ([]byte){
	confirm := ChangeConfirm{
			Message: Message{Type: "changeConfirm"},
			Name: name,
	}

	data, err := json.Marshal(confirm)
	if err != nil{
		log.Printf("Eror building pong. Error: %s", err)
		return []byte{}
	}

	return append(data, '\n')
}

func BuildPong(timestamp int64)([]byte){
	pong := Pong{
			Message: Message{Type: "pong"},
			PingTimestamp: timestamp,
	}

	data, err := json.Marshal(pong)
	if err != nil{
		log.Printf("Eror building pong. Error: %s", err)
		return []byte{}
	}

	return append(data, '\n')
}

