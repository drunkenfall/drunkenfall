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
	outgoing amqp.Queue
	ch       *amqp.Channel
	msgs     <-chan amqp.Delivery
}

// NewListener sets up a new listener
func NewListener(conf *Config, db *Database) (*Listener, error) {
	var err error
	l := Listener{
		DB: db,
	}

	l.conn, err = amqp.Dial(conf.RabbitURL)
	if err != nil {
		log.Fatal(err)
	}

	l.ch, err = l.conn.Channel()
	failOnError(err, "Failed to open a channel")

	l.outgoing, err = l.ch.QueueDeclare(
		conf.RabbitOutgoingQueue, // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare the outgoing queue")

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

	l.Publish("test", 9)

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

// Publish sends a message on the defualt exchange to the currently
// configured queue
func (l *Listener) Publish(kind string, data interface{}) error {
	msg := Message{
		kind,
		data,
		time.Now(),
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = l.ch.Publish(
		"",              // exchange
		l.outgoing.Name, // routing key
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		})
	return nil
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
