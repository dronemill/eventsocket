package eventsocket

import (
	"fmt"

	log "github.com/dronemill/eventsocket/Godeps/_workspace/src/github.com/Sirupsen/logrus"
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
	log.Info("Starting hub execution")

	go h.recvRegistrations()
	go h.recvUnregistrations()

	for {
		select {
		case m := <-h.recvClientMessage:
			h.ingest(m)
		}
	}
}

func (h *hub) recvRegistrations() {
	log.Info("Ready to receive client registrations")
	for {
		h.registerClient(<-h.register)
	}
}

func (h *hub) recvUnregistrations() {
	log.Info("Ready to receive client unregistrations")
	for {
		h.unregisterConnection(<-h.unregister)
	}
}

func (h *hub) registerClient(client *Client) {
	log.WithField("clientID", client.Id).Info("Registring client")
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

	// if we dont know about this conneciton, then get out
	if _, ok := h.connections[ws]; !ok {
		log.Info("Aborting unregistration of unknown connection")
		return
	}

	id := h.connections[ws]
	delete(h.connections, ws)
	log.WithField("clientID", id).Info("Unregistring connection")

	// if we dont have any clientSubscriptions, then leave
	if _, ok := h.clientSubscriptions[id]; !ok {
		return
	}

	// remove all client subscriptions
	for s := range h.clientSubscriptions[id] {
		h.purgeSubscription(id, s)
	}

	delete(h.clientSubscriptions, id)
	log.WithField("clientID", id).Info("Successfully unregistered connection")
}

// ingest a message and route it to the propper destination
func (h *hub) ingest(cm *ClientMessage) {
	log.WithField("clientID", cm.ClientId).Info("Injesting message")

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
		log.WithField("clientID", cm.ClientId).
			WithField("type", cm.Message.MessageType).
			Error("Unknown message type")
		// FIXME: dont just silently drop message ..
	}
}

// broadcast a message to all active clients
func (h *hub) handleBroadcast(cm *ClientMessage) {
	log.WithField("clientID", cm.ClientId).Info("Handling broadcast message")

	for ws, clientID := range h.connections {
		log.WithField("clientID", clientID).Debug("Transmitting broadcast message")
		h.send(ws, &cm.Message)
	}

	return
}

// handle a standard message
func (h *hub) handleStandard(cm *ClientMessage) {
	log.WithField("clientID", cm.ClientId).Info("Handling standard message")

	// ensure there is atleast one suscriber of this event
	subscribers, ok := h.subscriptions[cm.Message.Event]
	if !ok {
		log.WithField("clientID", cm.ClientId).
			WithField("event", cm.Message.Event).
			Warn("No subscribers for event")
		return
	}

	for s := range subscribers {
		log.WithField("clientID", s).Debug("Transmitting standard message")
		h.send(clients[s].ws, &cm.Message)
	}

	return
}

// suscribe a client to events
func (h *hub) handleSuscribe(cm *ClientMessage) {
	log.WithField("clientID", cm.ClientId).Info("Handling subscription request")

	// sanity check
	events, ok := cm.Message.Payload["Events"]
	if !ok {
		log.WithField("clientID", cm.ClientId).Error("No events provided")
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
	log.WithField("clientID", cm.ClientId).Info("Handling unsuscribe request")

	// sanity check
	events, ok := cm.Message.Payload["Events"]
	if !ok {
		log.WithField("clientID", cm.ClientId).Error("No events provided")
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
	log.WithField("fromClientID", cm.ClientId).
		WithField("toClientID", cm.Message.RequestClientId).
		Info("Handling request")

	// store the clientID sending the message as the client to reply to
	cm.Message.ReplyClientId = cm.ClientId

	if clients[cm.Message.RequestClientId] == nil {
		log.WithField("fromClientID", cm.ClientId).
			WithField("toClientID", cm.Message.RequestClientId).
			Error("Client does not exist")

		h.requestError(cm,
			"RequestClientID does not exist",
			ErrorRequestClientNoExist,
		)
		return
	}

	err := h.send(clients[cm.Message.RequestClientId].ws, &cm.Message)
	if err != nil {
		log.WithField("fromClientID", cm.ClientId).
			WithField("toClientID", cm.Message.RequestClientId).
			Error("Client is not connected")

		h.requestError(cm,
			"RequestClientID is not connected",
			ErrorRequestClientNotConnected,
		)
	}
}

// handle a reply. forward the message along to the requestee
func (h *hub) handleReply(cm *ClientMessage) {
	log.WithField("fromClientID", cm.ClientId).
		WithField("toClientID", cm.Message.ReplyClientId).
		Info("Handling reply")

	// FIXME need to have some form of error handling here..
	h.send(clients[cm.Message.ReplyClientId].ws, &cm.Message)
}

// store the subscription for the given clientId and event
func (h *hub) storeSubscription(id, e string) {
	log.WithField("clientID", id).WithField("event", e).
		Info("Storing subscription")

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
	log.WithField("clientID", id).WithField("event", e).
		Info("Purgin subscription")

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
		log.WithField("event", e).
			Debug("No suscribers remain")

		delete(h.subscriptions, e)
	}

}

func (h *hub) send(ws *wsConnection, m *Message) error {
	var err error
	err = nil
	defer func() {
		if r := recover(); r != nil {
			h.unregister <- ws
			err = fmt.Errorf("Send failed: %s", r)
		}
	}()

	if ws.closed {
		h.unregister <- ws
		err = fmt.Errorf("WS is closed")
	} else {
		ws.send <- *m
	}
	return err
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

	// we dont need to check for errors here, because if the request fails,
	// then it doesnt matter, as the client needing to know that it failed
	// is the one who cant be sent to..
	h.send(clients[cm.Message.ReplyClientId].ws, &m)
}
