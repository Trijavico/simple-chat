package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	hub := NewHub()
	mux := http.NewServeMux()

	go hub.Listen()

	mux.Handle("/", http.FileServer(http.Dir("templates")))

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		hub.ServeWebSocket(w, r)
	})

	fmt.Println("Listening on PORT: 3000")
	log.Fatal(http.ListenAndServe(":3000", mux))
}
