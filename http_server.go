package eventsocket

import (
	"fmt"
	"net"
	"net/http"

	"github.com/LiftMe/glip/log"
	"github.com/gorilla/mux"
)

type httpServer struct {
	router   *mux.Router
	listener *net.Listener
}

// install the http server's router
func (h *httpServer) route() error {
	// log.Info("Initializing EventSocket Router")

	// instantiate a new controller
	C, err := newHttpController()
	if err != nil {
		// log.Error(fmt.Sprintf("Encountered error while instantiating new HttpController: %s", err.Error()))
		return err
	}

	// get a new router
	h.router = mux.NewRouter()
	s := h.router.PathPrefix("/v1").Subrouter()

	// handle routes
	s.HandleFunc("/clients/{id}/ws", C.Client.ServeWs).Methods("GET")
	// s.HandleFunc("/clients/{cid}", C.Client.Get).Methods("GET")
	s.HandleFunc("/clients", C.Client.Create).Methods("POST")
	s.HandleFunc("/dev/ping", C.Dev.Ping).Methods("GET")

	return nil
}

// listen
func (h *httpServer) listen(listenAddr string) error {
	l, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Error(fmt.Sprintf("%s: %s", "Listen Error", err.Error()))
		return err
	}

	h.listener = &l
	return nil
}

// server
func (h *httpServer) serve() error {
	http.Handle("/", h.router)

	err := http.Serve(*h.listener, nil)
	if err != nil {
		return err
	}

	return nil
}
