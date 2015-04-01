package eventsocket

// hub maintains the set of active connections, and broadcasts messages
// for given events to where they need to go
type hub struct {
	// Registered connections.
	connections map[*wsConnection]bool

	// Inbound messages from the connections.
	broadcast chan []byte

	// Register requests from the connections.
	register chan *wsConnection

	// Unregister requests from connections.
	unregister chan *wsConnection
}

var h = hub{
	broadcast:   make(chan []byte),
	register:    make(chan *wsConnection),
	unregister:  make(chan *wsConnection),
	connections: make(map[*wsConnection]bool),
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
		case m := <-h.broadcast:
			for c := range h.connections {
				select {
				case c.send <- m:
				default:
					close(c.send)
					delete(h.connections, c)
				}
			}
		}
	}
}
