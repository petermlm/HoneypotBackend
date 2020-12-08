package main

import (
	"context"
	"honeypot/listener"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	str := []string{"8000"}
	listener.Start(ctx, str)
}
