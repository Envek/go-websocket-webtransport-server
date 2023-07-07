package main

import (
	"log"
	"time"
)

type Message struct {
	Author    string
	Text      string
	Timestamp time.Time
}

type Chat struct {
	clients    map[*Client]bool
	broadcast  chan Message
	register   chan *Client
	unregister chan *Client
}

func NewChat() *Chat {
	return &Chat{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (chat *Chat) run() {
	// chatHistory := make(map[time.Time]Message)

	for {
		select {
		case client := <-chat.register:
			chat.clients[client] = true
			log.Printf("New %T client registered!", client.stream)
			chat.send(client, Message{"Chat", "Hey stranger, what is your name?", time.Now()})

			// for _, message := range chatHistory {
			// 	chat.send(client, message)
			// }
		case client := <-chat.unregister:
			log.Printf("%T client leaves: %s", client.stream, client.username)
			if _, ok := chat.clients[client]; ok {
				chat.send(client, Message{"Chat", "Bye-bye!", time.Now()})
				delete(chat.clients, client)
				close(client.send)
			}
		case message := <-chat.broadcast:
			log.Println("Message from ", message.Author, ":", message.Text)
			for client := range chat.clients {
				chat.send(client, message)
			}
		case <-done:
			log.Println("Chat is closing, saying bye for everyone")
			for client := range chat.clients {
				chat.send(client, Message{"Chat", "Bye-bye!", time.Now()})
				delete(chat.clients, client)
				close(client.send)
			}
			return
		}
	}
}

func (chat *Chat) send(client *Client, message Message) {
	select {
	case client.send <- message:
	default:
		close(client.send)
		delete(chat.clients, client)
	}
}
