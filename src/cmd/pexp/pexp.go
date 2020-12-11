package main

import (
	"honeypot/queue"
	"honeypot/settings"
	"honeypot/timelines"
	"log"
	"time"
)

func main() {
	p, err := queue.NewPublisher(settings.RabbitmqTaskProcessConnAttemp)
	if err != nil {
		log.Fatal(err)
	}
	defer p.Destroy()
	p.Publish(&timelines.ConnAttemp{
		Time:       time.Now(),
		Port:       "1234",
		IP:         "178.5.1.161",
		ClientPort: "15432",
	})
}
