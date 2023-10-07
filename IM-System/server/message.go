package server

const (
	GlobalMsg = iota
	PrivateMsg
	SearchAllOnlineClientMsg
	RenameMsg
)

type Message struct {
	MsgType  int    `json:"msgType"`
	Content  string `json:"content"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
}

func CopyMessage(msg *Message) *Message {
	newMsg := &Message{MsgType: msg.MsgType, Content: msg.Content, Sender: msg.Sender, Receiver: msg.Receiver}
	return newMsg
}
