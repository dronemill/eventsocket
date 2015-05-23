package eventsocket

import (
	"errors"

	log "github.com/dronemill/eventsocket/Godeps/_workspace/src/github.com/Sirupsen/logrus"
)

type Server struct {
	Config struct {
		listenAddr string
	}
	HttpServer *httpServer
}

// return a new, registered instance of the eventsocket server
func NewServer(listenAddr string) (server *Server, err error) {
	// ensure that we have a listenAddr
	if listenAddr == "" {
		err = errors.New("Empty listenAddr")
		log.Error("Empty listenAddr")
		return
	}

	// instantiate new server
	server = new(Server)
	server.Config.listenAddr = listenAddr

	registerServer(server)

	return
}

func registerServer(server *Server) {
	log.WithField("listenAddr", server.Config.listenAddr).Info("Registering new server")
	server.HttpServer = &httpServer{}

	server.HttpServer.route()
}

func (server *Server) Start() error {
	log.WithField("listenAddr", server.Config.listenAddr).Info("Starting server")
	go h.run()

	if err := server.HttpServer.listen(server.Config.listenAddr); err != nil {
		return err
	}

	return server.HttpServer.serve()
}

func (server *Server) Stop() error {
	log.WithField("listenAddr", server.Config.listenAddr).Info("Stopping server")
	return (*server.HttpServer.listener).Close()
}

// maximum message size allowed from peer
func (server *Server) SetDefaultMaxMessageSize(limit int64) {
	log.WithField("listenAddr", server.Config.listenAddr).
		WithField("size", limit).
		Info("Setting default MaxMessageSize")
	defaultMaxMessageSize = limit
}
