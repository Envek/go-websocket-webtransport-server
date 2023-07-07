package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type WebSocketStream struct {
	conn *websocket.Conn
}

// Implement WebStream interface

func (ws WebSocketStream) ReadMessage() (string, error) {
	_, message, err := ws.conn.ReadMessage()
	if err != nil {
		return "", err
	}
	return string(message), nil
}

func (ws WebSocketStream) WriteMessage(message string) error {
	return ws.conn.WriteMessage(websocket.TextMessage, []byte(message))
}

func (ws WebSocketStream) Close() error {
	if err := ws.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
		return err
	}

	return ws.conn.Close()
}

// Set up a WebSocket server

func websocketServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte(indexpage)) })
	http.HandleFunc("/chat", wsHandler)
	srv := &http.Server{Addr: ":8090"}

	idleConnsClosed := make(chan struct{})
	go func() {
		<-done
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	log.Printf("WebSocket server starting")
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
	<-idleConnsClosed
	log.Printf("HTTP server has been shut down")
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client := NewClient(chat, WebSocketStream{conn})

	chat.register <- client

	go client.Reader()
	go client.Writer()
}
