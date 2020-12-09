package main

import (
	"context"
	"honeypot/queue"
	"log"
)

func main() {
	q, _ := queue.NewConsumer("exp")
	ch, _ := q.Consume(context.Background())
	for m := range ch {
		log.Println(m)
	}
}
