package main

import (
	"context"
	"honeypot/listener"
	"honeypot/settings"
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

	if err := listener.Start(ctx, settings.Ports[:]); err != nil {
		log.Println(err)
	}
}
