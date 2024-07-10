package main

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type UserReq struct {
	Name string `json:"name"`
}

type Message struct {
	ClientID   string
	ClientName string
	Text       string `json:"text"`
}

type Client struct {
	ID   string
	Name string
	conn *websocket.Conn

	send chan []byte
}

func (c *Client) readPump() {
	for {
		_, req, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			continue
		}

		msg := &Message{}

		err = json.Unmarshal(req, msg)
		if err != nil {
			log.Printf("error: %v\n", err)
			continue
		}

		msg_encoded, err := json.Marshal(msg)
		if err != nil {
			log.Printf("error: %v\n", err)
		}

		c.send <- msg_encoded

	}
}

func (c *Client) responsePump() {
	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			writer, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			writer.Write(msg)

			if err := writer.Close(); err != nil {
				return
			}
		}
	}
}
