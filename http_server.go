package eventsocket

import (
	"fmt"
	"net"
	"net/http"

	"github.com/LiftMe/glip/log"
	"github.com/gorilla/mux"
)

type httpServer struct {
	router *mux.Router
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
	// s.HandleFunc("/clients/{cid}/ws", C.Client.GetWs).Methods("GET")
	// s.HandleFunc("/clients/{cid}", C.Client.Get).Methods("GET")
	s.HandleFunc("/clients", C.Client.Create).Methods("POST")

	return nil
}

// handle and serve the api
func (h *httpServer) listen(listenAddr string) error {
	http.Handle("/", h.router)

	l, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Error(fmt.Sprintf("%s: %s", "Listen Error", err.Error()))
		return err
	}

	err = http.Serve(l, nil)
	if err != nil {
		log.Error(fmt.Sprintf("%s: %s", "Serve Error", err.Error()))
		return err
	}

	return nil
}
