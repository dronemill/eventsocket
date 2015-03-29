package eventsocket

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"
)

func Test_Server(t *testing.T) {
	server, err := NewServer("127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	if err := server.HttpServer.listen(server.Config.listenAddr); err != nil {
		t.Fatal(err)
	}
	defer (*server.HttpServer.listener).Close()

	go server.HttpServer.serve()

	resp, err := http.Post(fmt.Sprintf("http://%s/v1/clients", (*server.HttpServer.listener).Addr().String()), "application/json", bytes.NewBufferString(""))

	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected statys %v. Got: %v", http.StatusOK, resp.StatusCode)
	}
}
