package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
}

func main() {
    mux := http.NewServeMux()

    mux.Handle("/", http.FileServer(http.Dir("templates")))

    mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Origin") != "http://" + r.Host {
	    http.Error(w, "Origin not allowed", http.StatusForbidden)
	    return
	}

	_, err := upgrader.Upgrade(w, r, nil)
	if err != nil{
	    http.Error(w, "Handshake error", http.StatusBadRequest)
	    return
	}

	fmt.Println("client connected")
    })

    fmt.Println("Listening on PORT: 3000")
    log.Fatal(http.ListenAndServe(":3000", mux))
}
