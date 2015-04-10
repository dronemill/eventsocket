package eventsocket

type ClientMessage struct {
	ClientId string
	Message  Message
}

type Message struct {
	MessageType MessageType `json:MessageType`
	Event       string      `json:Event,omitempty`
	RequestId   string      `json:RequestId,omitempty`
	// ReplyTo         string                 `json:ReplyTo,omitempty`
	ReplyClientId   string                 `json:ReplyClientId,omitempty`
	RequestClientId string                 `json:RequestClientId,omitempty`
	Payload         map[string]interface{} `json:Payload`
}

type MessageType int

const MESSAGE_TYPE_BROADCAST = 1
const MESSAGE_TYPE_STANDARD = 2
const MESSAGE_TYPE_REQUEST = 3
const MESSAGE_TYPE_REPLY = 4
const MESSAGE_TYPE_SUSCRIBE = 5
const MESSAGE_TYPE_UNSUSCRIBE = 6
