package queue

import (
	"context"
	"honeypot/timelines"
)

type Consumer interface {
	Destroy()
	Name() string
	Consume(context.Context) (chan *timelines.ConnAttemp, error)
}

type consumer struct {
	*queue
}

func NewConsumer(name string) (Consumer, error) {
	q, err := newQueue(name)
	if err != nil {
		return nil, err
	}
	return &consumer{q}, nil
}

func (c *consumer) Destroy() {
	c.queue.Destroy()
}

func (c *consumer) Name() string {
	return c.queue.name()
}

func (c *consumer) Consume(ctx context.Context) (chan *timelines.ConnAttemp, error) {
	msgs, err := c.queue.ch.Consume(
		c.Name(), // queue
		"",       // consumer
		false,    // auto-ack
		false,    // exclusive
		false,    // no-local
		false,    // no-wait
		nil,      // args
	)
	if err != nil {
		return nil, err
	}

	ch := make(chan *timelines.ConnAttemp)
	go func() {
		for {
			select {
			case msg := <-msgs:
				connAttemp, err := timelines.ConnAttempFromJson(msg.Body)
				if err != nil {
					// TODO: Handle error
					continue
				}
				ch <- connAttemp
			case <-ctx.Done():
				return
			}
		}
	}()
	return ch, nil
}
