package eventsocket

import (
	"fmt"
	"os"
)

// hub maintains the set of active connections, and broadcasts messages
// for given events to where they need to go
type hub struct {
	// Inbound messages from the connections.
	recvClientMessage chan *ClientMessage

	// Register a client
	register chan *Client

	// Unregister requests from connections.
	unregister chan *wsConnection

	// Registered connection map from connection to clientId
	connections map[*wsConnection]string

	// subscriptions is a list of the events that have atleast one subscriber
	subscriptions hubSubscriptions

	// clientSubscriptions is a map of clientId to a slice of events
	// that the client has subscriped to
	clientSubscriptions map[string]map[string]bool
}

// hubSubscription is a list of the clients who have subscribed to a given event
type hubSubscription map[string]bool

// hubSubscriptions is a list of the events that have atleast one subscriber
type hubSubscriptions map[string]hubSubscription

var h = hub{
	recvClientMessage:   make(chan *ClientMessage),
	register:            make(chan *Client),
	unregister:          make(chan *wsConnection),
	connections:         make(map[*wsConnection]string),
	subscriptions:       make(map[string]hubSubscription),
	clientSubscriptions: make(map[string]map[string]bool),
}

func (h *hub) run() {
	for {
		select {
		case cl := <-h.register:
			h.registerClient(cl)
		case c := <-h.unregister:
			h.unregisterConnection(c)
		case m := <-h.recvClientMessage:
			h.ingest(m)
		}
	}
}

func (h *hub) registerClient(client *Client) {
	h.connections[client.ws] = client.Id
	h.clientSubscriptions[client.Id] = make(map[string]bool)

	for s := range client.subscriptions {
		h.storeSubscription(client.Id, s)
	}
}

// unsuscribe from all events for this connection..? or set to false
// also remove the connection->client reference, and close the send chan
func (h *hub) unregisterConnection(ws *wsConnection) {
	defer func() {
		ws.Close()
	}()

	// if we dont know about htis conneciton, then get out
	if _, ok := h.connections[ws]; !ok {
		return
	}

	id := h.connections[ws]
	delete(h.connections, ws)

	// if we dont have any clientSubscriptions, then leave
	if _, ok := h.clientSubscriptions[id]; !ok {
		return
	}

	// remove all client subscriptions
	for s := range h.clientSubscriptions[id] {
		h.purgeSubscription(id, s)
	}

	delete(h.clientSubscriptions, id)
}

// ingest a message and route it to the propper destination
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
	case MESSAGE_TYPE_REQUEST:
		h.handleRequest(cm)
	case MESSAGE_TYPE_REPLY:
		h.handleReply(cm)
	default:
		fmt.Printf("ERROR: unhandled MessageType:%v [ClientId:%s]\n", cm.Message.MessageType, cm.ClientId)
		os.Exit(1)
	}
}

// broadcast a message to all active clients
func (h *hub) handleBroadcast(cm *ClientMessage) {
	for ws, _ := range h.connections {
		select {
		case ws.send <- cm.Message:
		default:
			h.unregister <- ws
		}

	}

	return
}

// handle a standard message
func (h *hub) handleStandard(cm *ClientMessage) {
	// ensure there is atleast one suscriber of this event
	subscribers, ok := h.subscriptions[cm.Message.Event]
	if !ok {
		return
	}

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
		h.storeSubscription(cm.ClientId, e)
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
		h.purgeSubscription(cm.ClientId, e)

		// remove the event from the client id subscriptions map
		if _, ok := clients[cm.ClientId].subscriptions[e]; ok {
			delete(clients[cm.ClientId].subscriptions, e)
		}
	}

	return
}

// handle a request message. this is done by storing the requestee's
// ClientId as the ReplyTo, and then forwarding the message along
// to the receiving client
func (h *hub) handleRequest(cm *ClientMessage) {
	cm.Message.ReplyClientId = cm.ClientId

	if clients[cm.Message.RequestClientId] == nil {
		h.requestError(cm,
			"RequestClientID does not exist or is not connected",
			ErrorRequestClientNoExist,
		)
		return
	}

	clients[cm.Message.RequestClientId].ws.send <- cm.Message
}

// handle a reply. forward the message along to the requestee
func (h *hub) handleReply(cm *ClientMessage) {
	clients[cm.Message.ReplyClientId].ws.send <- cm.Message
}

// store the subscription for the given clientId and event
func (h *hub) storeSubscription(id, e string) {
	// make sure we have a subscription key
	if _, ok := h.subscriptions[e]; !ok {
		h.subscriptions[e] = make(hubSubscription)
	}

	// store the clientid in the subscription map
	h.subscriptions[e][id] = true

	h.clientSubscriptions[id][e] = true

	// store the event in the client id subscriptions map
	clients[id].subscriptions[e] = true
}

// purge the subscription
func (h *hub) purgeSubscription(id, e string) {
	// remove the clientSubscriptions map, if it exists
	if _, ok := h.clientSubscriptions[id][e]; ok {
		delete(h.clientSubscriptions[id], e)
	}

	// make sure we are actually suscribing to this event
	if _, ok := h.subscriptions[e][id]; !ok {
		return
	}

	// remove the suscription
	delete(h.subscriptions[e], id)

	// if no one else if suscribing, then remove the suscription key
	if len(h.subscriptions[e]) == 0 {
		delete(h.subscriptions, e)
	}

}

func (h *hub) requestError(cm *ClientMessage, message, etype string) {
	e := map[string]interface{}{
		"Message": message,
		"Type":    etype,
	}

	m := Message{
		MessageType: MESSAGE_TYPE_REPLY,
		RequestId:   cm.Message.RequestId,
		Error:       e,
	}

	clients[cm.Message.ReplyClientId].ws.send <- m
}
