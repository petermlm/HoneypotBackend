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
		IP:         "94.46.160.198",
		ClientPort: "15432",
	})
}
