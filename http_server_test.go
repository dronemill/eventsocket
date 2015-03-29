package eventsocket

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func testHttpServer(t *testing.T) (*httpServer, *httptest.ResponseRecorder) {
	server := &httpServer{}

	err := server.route()
	if err != nil {
		t.Fatalf("Received error while installing router %s", err.Error())
	}

	w := httptest.NewRecorder()

	return server, w
}

func Test_HttpServer_Router(t *testing.T) {
	server, w := testHttpServer(t)

	req, _ := http.NewRequest("GET", "/v1/dev/ping", nil)
	server.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Dev ping page didn't return %v. Got: %v", http.StatusOK, w.Code)
	}

	if s := w.Body.String(); s != "pong" {
		t.Fatalf("Expected \"pong\" got \"%s\"", s)
	}

	return
}
