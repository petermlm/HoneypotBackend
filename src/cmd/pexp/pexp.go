package main

import (
	"honeypot/queue"
	"honeypot/timelines"
	"log"
	"time"
)

func main() {
	p, err := queue.NewPublisher("exp")
	if err != nil {
		log.Fatal(err)
	}
	p.Publish(&timelines.ConnAttemp{
		Time:       time.Now(),
		Port:       "1234",
		IP:         "192.168.0.1",
		ClientPort: "15432",
	})
}
