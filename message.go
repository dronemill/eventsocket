package eventsocket

type ClientMessage struct {
	ClientId string
	Message  Message
}

type Message struct {
	MessageType MessageType            `json:messageType`
	Payload     map[string]interface{} `json:payload`
}

type MessageType int

const MESSAGE_TYPE_BROADCAST = 0
const MESSAGE_TYPE_STANDARD = 1
const MESSAGE_TYPE_REQUEST = 2
const MESSAGE_TYPE_REPLY = 3
const MESSAGE_TYPE_SUSCRIBE = 4
const MESSAGE_TYPE_UNSUSCRIBE = 5
