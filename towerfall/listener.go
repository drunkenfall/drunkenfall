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
	DB       *Database
	conn     *amqp.Connection
	incoming amqp.Queue
	ch       *amqp.Channel
	msgs     <-chan amqp.Delivery
	log      *zap.Logger
}

// NewListener sets up a new listener
func NewListener(conf *Config, db *Database) (*Listener, error) {
	var err error
	l := Listener{
		DB:  db,
		log: conf.log,
	}

	l.conn, err = amqp.Dial(conf.RabbitURL)
	if err != nil {
		log.Fatal(err)
	}

	l.ch, err = l.conn.Channel()
	failOnError(err, "Failed to open a channel")

	l.incoming, err = l.ch.QueueDeclare(
		conf.RabbitIncomingQueue, // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare the incoming queue")

	l.msgs, err = l.ch.Consume(
		conf.RabbitIncomingQueue, // queue
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	failOnError(err, "Failed to register the consumer")

	return &l, err
}

func (l *Listener) Serve() {
	log.Println("Listening for AMQP messages...")
	defer l.ch.Close()

	for d := range l.msgs {
		t, err := l.DB.GetCurrentTournament(l.DB.Server)
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

	// Check if this is a meta communication from the game.
	switch msg.Type {
	case gConnect:
		// If this is a connect event, it means that the game wants to
		// know what's up with the current match.
		t.connect(true)
		if err := t.PublishNext(); err != nil && err != ErrPublishDisconnected {
			return err
		}
		return nil

	case gDisconnect:
		t.connect(false)
		return nil
	}

	// If it wasn't, then it's about a match
	return t.Matches[t.Current].handleMessage(msg)
}
