package main

import (
	"context"
	"honeypot/listener"
	"log"
	"os"
	"os/signal"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	defer func() {
		signal.Stop(c)
		cancel()
	}()

	go func() {
		select {
		case <-c:
			cancel()
		}
	}()

	str := []string{"8000"}
	if err := listener.Start(ctx, str); err != nil {
		log.Println(err)
	}
}
