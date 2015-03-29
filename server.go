package eventsocket

import "errors"

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
		// log.Fatal(err.Error())
		return
	}

	// instantiate new server
	server = new(Server)
	server.Config.listenAddr = listenAddr

	registerServer(server)

	return
}

func registerServer(server *Server) {
	server.HttpServer = &httpServer{}

	server.HttpServer.route()
}

func (server *Server) Start() error {
	if err := server.HttpServer.listen(server.Config.listenAddr); err != nil {
		return err
	}

	return server.HttpServer.serve()
}

func (server *Server) Stop() error {
	return (*server.HttpServer.listener).Close()
}
