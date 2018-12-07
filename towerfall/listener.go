package towerfall

import (
	"encoding/json"
	"log"
	"time"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

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
	if err != nil {
		return nil, errors.New("Failed to open a channel")
	}

	l.incoming, err = l.ch.QueueDeclare(
		conf.RabbitIncomingQueue, // name
		false,                    // durable
		false,                    // delete when unused
		false,                    // exclusive
		false,                    // no-wait
		nil,                      // arguments
	)

	if err != nil {
		return nil, errors.New("Failed to declare the incoming queue")
	}

	l.msgs, err = l.ch.Consume(
		conf.RabbitIncomingQueue, // queue
		"",                       // consumer
		true,                     // auto-ack
		false,                    // exclusive
		false,                    // no-local
		false,                    // no-wait
		nil,                      // args
	)

	if err != nil {
		return nil, errors.New("Failed to register the consumer")
	}

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

	// Set connection to true so that if the game restarts we can assume we can send messages
	t.connect(true)

	// If it wasn't, then it's about a match
	m, err := t.CurrentMatch()
	if err != nil {
		l.log.Info("Couldn't find current match", zap.Error(err))
	}

	err = m.handleMessage(msg)
	if err != nil {
		l.log.Info("Match handle failed", zap.Error(err))
		return err
	}

	return nil
}
