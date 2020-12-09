package queue

import (
	"github.com/streadway/amqp"
)

type queue struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	q    amqp.Queue
}

func newQueue(name string) (*queue, error) {
	conn, err := amqp.Dial(makeConnString())
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		// No defer, we things to be open at return
		conn.Close()
		return nil, err
	}

	q, err := ch.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		// No defer, we things to be open at return
		conn.Close()
		ch.Close()
		return nil, err
	}

	ret := &queue{
		conn: conn,
		ch:   ch,
		q:    q,
	}
	return ret, nil
}

func (q *queue) Destroy() {
	q.conn.Close()
	q.ch.Close()
}

func (q *queue) name() string {
	return q.q.Name
}
