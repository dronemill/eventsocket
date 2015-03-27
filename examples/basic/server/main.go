package main

import (
	"flag"
	"log"

	"github.com/dronemill/eventsocket"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()

	_, err := eventsocket.NewServer(*addr)
	if err != nil {
		log.Fatal(err)
	}
}
