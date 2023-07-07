package main

import (
	_ "embed"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var chat = NewChat()
var done = make(chan struct{})

//go:embed index.html
var indexpage string

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		chat.run()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		webtransportServer()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		websocketServer()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		sig := <-signals
		log.Printf("Got signal: %v, shutting down…", sig)
		close(done)
	}()

	wg.Wait()
}
