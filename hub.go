package eventsocket

import "fmt"

// hub maintains the set of active connections, and broadcasts messages
// for given events to where they need to go
type hub struct {
	// Inbound messages from the connections.
	recvClientMessage chan *ClientMessage
}

// hubSubscription is a list of the clients who have subscribed to a given event
type hubSubscription map[string]bool

// hubSubscriptions is a list of the events that have atleast one subscriber
var hubSubscriptions = make(map[string]hubSubscription)

var h = hub{
	recvClientMessage: make(chan *ClientMessage),
}

func (h *hub) run() {
	for {
		select {
		case m := <-h.recvClientMessage:
			h.ingest(m)
		}
	}
}

func (h *hub) ingest(cm *ClientMessage) {
	switch cm.Message.MessageType {
	case MESSAGE_TYPE_BROADCAST:
		h.handleBroadcast(cm)
	case MESSAGE_TYPE_STANDARD:
		h.handleStandard(cm)
	case MESSAGE_TYPE_SUSCRIBE:
		h.handleSuscribe(cm)
	case MESSAGE_TYPE_UNSUSCRIBE:
		h.handleUnsuscribe(cm)
	default:
		fmt.Printf("ERROR: unhandled MessageType:%v [ClientId:%s]\n", cm.Message.MessageType, cm.ClientId)
	}
}

// broadcast a message to all active clients
func (h *hub) handleBroadcast(cm *ClientMessage) {
	for _, c := range clients {
		select {
		case c.ws.send <- cm.Message:
		default:
			close(c.ws.send)
			// delete(h.connections, c)
			fmt.Println("Need to clean up this connection")
		}

	}

	return
}

// handle a standard message
func (h *hub) handleStandard(cm *ClientMessage) {
	// ensure there is atleast one suscriber of this event
	subscribers, ok := hubSubscriptions[cm.Message.Event]
	if !ok {
		return
	}

	fmt.Printf("%+v\n", subscribers)

	for s := range subscribers {
		clients[s].ws.send <- cm.Message
	}

	return
}

// suscribe a client to events
func (h *hub) handleSuscribe(cm *ClientMessage) {
	// sanity check
	events, ok := cm.Message.Payload["Events"]
	if !ok {
		return
	}

	for _, event := range events.([]interface{}) {
		e := event.(string)
		// make sure we have a subscription key
		if _, ok := hubSubscriptions[e]; !ok {
			hubSubscriptions[e] = make(hubSubscription)
		}

		hubSubscriptions[e][cm.ClientId] = true
	}

	return
}

// unsuscribe a client from events
func (h *hub) handleUnsuscribe(cm *ClientMessage) {
	// sanity check
	events, ok := cm.Message.Payload["Events"]
	if !ok {
		return
	}

	for _, event := range events.([]interface{}) {
		e := event.(string)
		// make sure we are actually suscribing to this event
		if _, ok := hubSubscriptions[e][cm.ClientId]; !ok {
			continue
		}

		// remove the suscription
		delete(hubSubscriptions[e], cm.ClientId)

		// if no one else if suscribing, then remove the suscription key
		if len(hubSubscriptions[e]) == 0 {
			delete(hubSubscriptions, e)
		}
	}

	return
}
