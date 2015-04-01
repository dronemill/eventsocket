package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"
)

var addr = flag.String("addr", ":8081", "http service address")
var srvPort = flag.String("srvPort", "8080", "http service address")

var homeTempl = template.Must(template.ParseFiles("./index.html"))

func serveHome(w http.ResponseWriter, r *http.Request) {
	// sanity check
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}

	// sanity check the method
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	// send html
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	homeTempl.Execute(w, fmt.Sprintf("%s:%s", strings.Split(r.Host, ":")[0], *srvPort))
}

func main() {
	flag.Parse()

	http.HandleFunc("/", serveHome)

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("Failed to ListenAndServe: ", err)
	}
}
