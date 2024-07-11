package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type UserReq struct {
	Name string `json:"name"`
}

func main() {
	hub := NewHub()
	mux := http.NewServeMux()

	go hub.Listen()

	mux.Handle("/", http.FileServer(http.Dir("templates")))

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Origin") != "http://"+r.Host {
			http.Error(w, "Origin not allowed", http.StatusForbidden)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "Handshake error", http.StatusInternalServerError)
			return
		}

		var usrName UserReq
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()

		if err != nil {
			http.Error(w, "Server error", http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(body, &usrName)
		if err != nil {
			usrName.Name = "Default Name"
		}

		client := &Client{
			ID:   uuid.NewString(),
			Name: usrName.Name,
			conn: conn,
			hub:  hub,
			send: make(chan []byte),
		}

		hub.register <- client

		go client.readPump()
		go client.responsePump()
	})

	fmt.Println("Listening on PORT: 3000")
	log.Fatal(http.ListenAndServe(":3000", mux))
}
