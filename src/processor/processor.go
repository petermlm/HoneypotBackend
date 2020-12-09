package processor

import (
	"context"
	"honeypot/queue"
	"honeypot/settings"
	"honeypot/timelines"
)

type env struct {
	c  queue.Consumer
	tl timelines.TimelinesWriter
}

func Start() error {
	c, err := queue.NewConsumer(settings.RabbitmqTaskProcessConnAttemp)
	if err != nil {
		return err
	}

	p, err := queue.NewPublisher(settings.RabbitmqTaskStoreConnAttemp)
	if err != nil {
		return err
	}

	// e := env{
	// 	c:  c,
	// 	tl: tl,
	// }

	ch, _ := c.Consume(context.Background())
	for m := range ch {
		p.Publish(m)
	}

	return nil
}
