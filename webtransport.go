package main

import (
	"context"
	"log"
	"net/http"

	"github.com/adriancable/webtransport-go"
)

type WebTransportStream struct {
	stream webtransport.Stream
}

// Implement WebStream interface

func (ws WebTransportStream) ReadMessage() (string, error) {
	buf := make([]byte, 1024)
	n, err := ws.stream.Read(buf)

	if err != nil {
		return "", err
	}

	return string(buf[:n]), nil
}

func (ws WebTransportStream) WriteMessage(message string) error {
	_, err := ws.stream.Write([]byte(message))
	return err
}

func (ws WebTransportStream) Close() error {
	return ws.stream.Close()
}

// Set up a WebTransport server

func webtransportServer() {
	server := &webtransport.Server{
		ListenAddr: ":4433",
		TLSCert:    webtransport.CertFile{Path: "cert.pem"},
		TLSKey:     webtransport.CertFile{Path: "cert.key"},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte(indexpage)) })
	mux.HandleFunc("/chat", wtHandler)
	server.Handler = mux

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-done
		cancel()
	}()

	log.Printf("WebTransport server starting")
	if err := server.Run(ctx); err != context.Canceled {
		log.Print("WebTransport server error:", err)
		cancel()
	}
	log.Printf("WebTransport server has been shut down")
}

func wtHandler(w http.ResponseWriter, r *http.Request) {
	session := r.Body.(*webtransport.Session)
	session.AcceptSession()
	// session.RejectSession(400)

	log.Println("Accepted incoming WebTransport session")

	s, err := session.OpenStreamSync(session.Context())
	if err != nil {
		log.Println(err)
	}
	log.Printf("Listening on server-initiated bidi stream %v\n", s.StreamID())

	client := NewClient(chat, WebTransportStream{s})

	chat.register <- client

	go client.Reader()
	go client.Writer()
}
