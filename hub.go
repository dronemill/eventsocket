package eventsocket

import "fmt"

// hub maintains the set of active connections, and broadcasts messages
// for given events to where they need to go
type hub struct {
	// Registered connections.
	connections map[*wsConnection]bool

	// Inbound messages from the connections.
	recvClientMessage chan *ClientMessage

	// Handle message subscriptions
	subscribe chan *hubSubscription

	// Handle message subscriptions
	unsubscribe chan *hubSubscription

	// Register requests from the connections.
	register chan *wsConnection

	// Unregister requests from connections.
	unregister chan *wsConnection
}

type hubSubscription struct {
	event  string
	result chan error
}

var h = hub{
	recvClientMessage: make(chan *ClientMessage),
	register:          make(chan *wsConnection),
	unregister:        make(chan *wsConnection),
	subscribe:         make(chan *hubSubscription),
	unsubscribe:       make(chan *hubSubscription),
	connections:       make(map[*wsConnection]bool),
}

func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
			h.connections[c] = true
		case c := <-h.unregister:
			if _, ok := h.connections[c]; ok {
				delete(h.connections, c)
				close(c.send)
			}
		case m := <-h.recvClientMessage:
			h.ingest(m)
		}
	}
}

func (h *hub) ingest(cm *ClientMessage) {
	switch cm.Message.MessageType {
	case MESSAGE_TYPE_BROADCAST:
		h.handleBroadcast(cm)
	default:
		fmt.Printf("ERROR: unhandled MessageType:%v [ClientId:%s]\n", cm.Message.MessageType, cm.ClientId)
	}
}

func (h *hub) handleBroadcast(cm *ClientMessage) error {
	for c := range h.connections {
		select {
		case c.send <- cm.Message:
		default:
			close(c.send)
			delete(h.connections, c)
		}
	}

	return nil
}
