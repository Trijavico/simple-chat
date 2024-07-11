package main

import (
	"bytes"
	"log"
	"sync"
	"text/template"
)

type Hub struct {
	sync.RWMutex

	clients    map[*Client]bool
	broadcast  chan *Message
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan *Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (hub *Hub) Listen() {
	for {
		select {
		case message, ok := <-hub.broadcast:
			if !ok {
				log.Println("messge channel errors")
			}

			hub.Lock()
			for client := range hub.clients {
				select {
				case client.send <- getTemplateMessage(message):
				default:
					close(client.send)
					delete(hub.clients, client)
				}
			}
			hub.Unlock()

		case client, ok := <-hub.unregister:
			if !ok {
				log.Println("messge channel errors")
			}

			hub.Lock()
			if _, ok := hub.clients[client]; ok {
				close(client.send)
				delete(hub.clients, client)
			}
			hub.Unlock()
			log.Println("client unregister")

		case client, ok := <-hub.register:
			if !ok {
				log.Println("messge channel errors")
			}

			hub.Lock()
			hub.clients[client] = true
			hub.Unlock()
			log.Println("client register")
		}
	}
}

func getTemplateMessage(msg *Message) []byte {
	templ, err := template.ParseFiles("templates/message.html")
	if err != nil {
		log.Fatalf("error template: %v\n", err)
	}

	var buf bytes.Buffer
	err = templ.Execute(&buf, msg)
	if err != nil {
		log.Fatalf("error encoding template: %v\n", err)
	}

	return buf.Bytes()
}
