package storer

import (
	"context"
	"honeypot/queue"
	"honeypot/settings"
	"honeypot/timelines"
)

func Start() error {
	c, err := queue.NewConsumer(settings.RabbitmqTaskStoreConnAttemp)
	if err != nil {
		return err
	}
	defer c.Destroy()

	tl := timelines.InitTimelines()
	defer tl.Close()

	ch, _ := c.Consume(context.Background())
	for m := range ch {
		tl.InsertConnAttemp(m)
	}

	return nil
}
