package websocket

type MessageType string

const (
	MessageTypePresence MessageType = "presence"
	MessageTypeChanges  MessageType = "BOARD_UPDATE"
	MessageTypeError    MessageType = "error"
	MessageTypeJoin     MessageType = "join"
	MessageTypeLeave    MessageType = "leave"
)

type Message struct {
	Type      MessageType `json:"type"`
	RoomID    string      `json:"roomId"`
	ClientID  string      `json:"clientId"`
	Username  string      `json:"username"`
	Payload   interface{} `json:"payload,omitempty"`
	Timestamp int64       `json:"timestamp"`
}

type Changes struct {
	Shapes        []interface{} `json:"shapes"`
	IsInitialSync bool          `json:"isInitialSync,omitempty"`
}

type HistoryPayload struct {
	History       []Message `json:"history"`
	IsInitialSync bool      `json:"isInitialSync"`
}
