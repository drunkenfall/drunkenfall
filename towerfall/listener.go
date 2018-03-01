package towerfall

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
)

type Message struct {
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

// Kill reasons
const (
	ReasonArrow = iota
	ReasonJumpedOn
	ReasonExplosion
)

// Message types
const (
	inKill = "kill"
)

type KillMessage struct {
	Killer int `json:"killer"`
	Corpse int `json:"corpse"`
	Action int `json:"action"`
}

type Listener struct {
	DB       *Database
	listener net.Listener
	port     int
}

// NewListener sets up a new listener
func NewListener(db *Database, port int) (*Listener, error) {
	var err error
	l := Listener{
		DB:   db,
		port: port,
	}

	l.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
	return &l, err
}

func (l *Listener) Serve() {
	log.Printf("Listening for messages on %d...", l.port)

	// run loop forever (or until ctrl-c)
	for {
		conn, err := l.listener.Accept()
		if err != nil {
			log.Printf("Error when accepting connection: %s", err.Error())
		}

		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Printf("Error when receiving message: %s", err.Error())
		}

		log.Println("Incoming message:", string(message))
		go l.handle(message)
	}
}

func (l *Listener) handle(msg string) {
	in := Message{
		Timestamp: time.Now(),
	}
	err := json.Unmarshal([]byte(msg), &in)

	if err != nil {
		log.Fatal(err)
	}

	switch in.Type {
	case inKill:
		l.KillMessage(in.Data.(KillMessage))
	default:
		log.Print("Warning: Unknown message type '%s'", in.Type)
	}
}

func (l *Listener) KillMessage(km KillMessage) error {
	return nil
}
