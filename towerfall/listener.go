package towerfall

import (
	"encoding/json"
	"log"
	"time"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

type Listener struct {
	DB    *Database
	conn  *amqp.Connection
	queue amqp.Queue
	ch    *amqp.Channel
	msgs  <-chan amqp.Delivery
}

// NewListener sets up a new listener
func NewListener(db *Database) (*Listener, error) {
	var err error
	l := Listener{
		DB: db,
	}

	l.conn, err = amqp.Dial("amqp://rabbitmq:thiderman@drunkenfall.com:5672/")
	if err != nil {
		log.Fatal(err)
	}

	l.ch, err = l.conn.Channel()
	failOnError(err, "Failed to open a channel")

	l.queue, err = l.ch.QueueDeclare(
		"drunkenfall-events", // name
		false,                // durable
		false,                // delete when unused
		false,                // exclusive
		false,                // no-wait
		nil,                  // arguments
	)
	failOnError(err, "Failed to declare a queue")

	l.msgs, err = l.ch.Consume(
		l.queue.Name, // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	failOnError(err, "Failed to register a consumer")

	return &l, err
}

func (l *Listener) Serve() {
	log.Println("Listening for AMQP messages...")
	defer l.ch.Close()

	for d := range l.msgs {
		t, err := l.DB.GetCurrentTournament()
		if err != nil {
			log.Printf("Could not get current tournament, skipping message: %s", err.Error())
			continue
		}

		msg := d.Body
		go func() {
			err := l.handle(t, msg)
			if err != nil {
				log.Printf("Handling failed: %s", err.Error())
			}
		}()
	}
}

func (l *Listener) handle(t *Tournament, body []byte) error {
	log.Printf("Incoming message: %s", body)
	msg := Message{
		Timestamp: time.Now(),
	}

	err := json.Unmarshal(body, &msg)
	if err != nil {
		return err
	}

	//
	return t.Matches[t.Current].handleMessage(msg)
}
