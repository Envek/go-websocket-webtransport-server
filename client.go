package main

import (
	"fmt"
	"log"
	"strings"
	"time"
)

type Client struct {
	chat     *Chat
	stream   WebStream
	username string
	send     chan Message
}

// Common interface for websocket.Conn and webtransport.Stream
type WebStream interface {
	ReadMessage() (string, error)
	WriteMessage(string) error
	Close() error
}

func NewClient(chat *Chat, stream WebStream) *Client {
	return &Client{
		chat:   chat,
		stream: stream,
		send:   make(chan Message, 100),
	}
}

func (c *Client) Writer() {
	defer c.stream.Close()

	for message := range c.send {
		formatted := fmt.Sprintf("%s [%s]: %s\n", message.Author, message.Timestamp.Format("2006-01-02 15:04:05"), message.Text)
		if err := c.stream.WriteMessage(formatted); err != nil {
			log.Println("Error sending message: ", err)
			return
		}
	}

	// The hub closed the channel, say bye-bye and close.
	c.stream.WriteMessage("Bye-bye!")
}

func (c *Client) Reader() {
	for {
		message, err := c.stream.ReadMessage()
		if err != nil {
			break
		}
		messageText := strings.TrimSuffix(message, "\n")

		if strings.HasPrefix(message, "Bye") {
			chat.unregister <- c
			return
		}

		if c.username == "" {
			c.username = strings.TrimSpace(messageText)
			c.send <- Message{
				Author:    "Chat",
				Text:      fmt.Sprintf("Welcome, %s!", c.username),
				Timestamp: time.Now(),
			}
			log.Printf("%T user set up: %s", c.stream, c.username)
			continue
		}

		c.chat.broadcast <- Message{
			Author:    c.username,
			Text:      messageText,
			Timestamp: time.Now(),
		}
	}
}
