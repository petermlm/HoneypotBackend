package queue

import (
	"honeypot/settings"
	"log"
	"time"

	"github.com/streadway/amqp"
)

type queue struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	q    amqp.Queue
}

func newQueue(name string) (*queue, error) {
	var conn *amqp.Connection
	var err error

	log.Println("Connecting with RabbitMQ...")
	for i := 0; i < settings.ConnectionRetriesTotal; i++ {
		conn, err = amqp.Dial(makeConnString())
		if err == nil {
			break
		}
		log.Printf("...attemp %d\n", i)
		time.Sleep(settings.ConnectionRetriesWait)
	}
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

	log.Println("Connection established")
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
