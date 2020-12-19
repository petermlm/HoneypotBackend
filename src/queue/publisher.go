package queue

import (
	"honeypot/timelines"

	"github.com/streadway/amqp"
)

type Publisher interface {
	Destroy()
	Name() string
	Publish(*timelines.ConnAttemp) error
}

type publisher struct {
	*queue
}

func NewPublisher(name string) (Publisher, error) {
	q, err := newQueue(name)
	if err != nil {
		return nil, err
	}
	return &publisher{q}, nil
}

func (p *publisher) Destroy() {
	p.queue.Destroy()
}

func (p *publisher) Name() string {
	return p.queue.name()
}

func (p *publisher) Publish(connAttemp *timelines.ConnAttemp) error {
	js, err := connAttemp.ToJSON()
	if err != nil {
		return err
	}
	return p.queue.ch.Publish(
		"",       // exchange
		p.Name(), // routing key
		false,    // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         js,
		})
}
