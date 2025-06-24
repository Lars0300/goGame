package protocol



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
	UpdateType UpdateStatus `json:"update_type"`
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
	PingTimestamp int64 `json:"ping_timestamp"`
}

type UpdateStatus int

const (
	StartGame UpdateStatus = iota
	EndGame
	CreateGame
	JoinGame
	KillPlayer
	Pass
	VoteStart
	ChatAlive
	ChatBomb
	ChatDead

)