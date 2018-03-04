package towerfall

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type Listener struct {
	DB       *Database
	listener net.Conn
	addr     *net.UDPAddr
	port     int
}

// NewListener sets up a new listener
func NewListener(db *Database, port int) (*Listener, error) {
	var err error
	l := Listener{
		DB:   db,
		port: port,
	}

	l.addr, err = net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}
	return &l, err
}

func (l *Listener) Serve() {
	var err error
	log.Printf("Listening for messages on :%d...", l.port)
	l.listener, err = net.ListenUDP("udp", l.addr)
	if err != nil {
		log.Print("Could not bind UDP")
		log.Fatal(err)
	}

	for {
		// conn, err := l.listener.Accept()
		// if err != nil {
		// 	log.Printf("Error when accepting connection: %s", err.Error())
		// }

		b := make([]byte, 1024)
		_, err := l.listener.Read(b)

		message := string(bytes.Trim(b, "\x00"))
		// message, err := bufio.NewReader(conn).ReadString('\n')
		// EOF is returned every time because the message sending code on
		// the TowerFall end will always close the connection as soon as
		// the message is sent.
		if err != nil && err != io.EOF {
			log.Printf("Error when receiving message: %s", err.Error())
		}

		// err = conn.Close()
		// if err != nil {
		// 	log.Printf("Couldn't close: %s", err.Error())
		// }

		t, err := l.DB.GetCurrentTournament()
		if err != nil {
			log.Printf("Could not get current tournament, skipping message: %s", err.Error())
			continue
		}

		go func() {
			err := l.handle(t, message)
			if err != nil {
				log.Printf("Handling failed: %s", err.Error())
			}
		}()
	}
}

func (l *Listener) handle(t *Tournament, body string) error {
	log.Println("Incoming message:", string(body))
	msg := Message{
		Timestamp: time.Now(),
	}

	err := json.Unmarshal([]byte(body), &msg)
	if err != nil {
		return err
	}

	//
	return t.Matches[t.Current].handleMessage(msg)
}
