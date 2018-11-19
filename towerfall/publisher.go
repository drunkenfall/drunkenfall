package towerfall

import (
	"encoding/json"
	"log"
	"time"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

type Publisher struct {
	conn     *amqp.Connection
	outgoing amqp.Queue
	ch       *amqp.Channel
}

// NewPublisher sets up a new listener
func NewPublisher(conf *Config) (*Publisher, error) {
	var err error
	p := Publisher{}

	p.conn, err = amqp.Dial(conf.RabbitURL)
	if err != nil {
		log.Fatal(err)
	}

	p.ch, err = p.conn.Channel()
	if err != nil {
		return &p, errors.New("Failed to open publishing channel")
	}

	p.outgoing, err = p.ch.QueueDeclare(
		conf.RabbitOutgoingQueue, // name
		false,                    // durable
		false,                    // delete when unused
		false,                    // exclusive
		false,                    // no-wait
		nil,                      // arguments
	)
	if err != nil {
		return &p, errors.New("Failed to declare the outgoing queue")
	}

	return &p, err
}

// Publish sends a message on the default exchange to the currently
// configured queue
func (p *Publisher) Publish(kind string, data interface{}) error {
	msg := Message{
		Type:      kind,
		Data:      data,
		Timestamp: time.Now(),
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = p.ch.Publish(
		"",              // exchange (default)
		p.outgoing.Name, // routing key
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		})
	return err
}
