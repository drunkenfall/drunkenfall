package towerfall

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/mitchellh/mapstructure"
)

type Message struct {
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

// Kill reasons
const (
	rArrow = iota
	rExplosion
	rBrambles
	rJumpedOn
	rLava
	rShock
	rSpikeBall
	rFallingObject
	rSquish
	rCurse
	rMiasma
	rEnemy
	rChalice
)

// Arrows
const (
	aNormal = iota
	aBomb
	aSuperBomb
	aLaser
	aBramble
	aDrill
	aBolt
	aToy
	aFeather
	aTrigger
	aPrism
)

// Message types
const (
	inKill       = "kill"
	inRoundStart = "round_start"
	inRoundEnd   = "round_end"
	inMatchStart = "match_start"
	inMatchEnd   = "match_end"
	inPickup     = "arrows_collected"
	inFire       = "arrow_shot"
)

type KillMessage struct {
	Player int `json:"player"`
	Killer int `json:"killer"`
	Cause  int `json:"cause"`
}

type PickupMessage struct {
	Player int            `json:"player"`
	Arrows map[string]int `json:"arrows"`
}

type Listener struct {
	DB       *Database
	listener net.Conn
	port     int
}

// NewListener sets up a new listener
func NewListener(db *Database, port int) (*Listener, error) {
	var err error
	l := Listener{
		DB:   db,
		port: port,
	}

	// l.listener, err = net.Listen("udp", fmt.Sprintf(":%d", port))
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}
	l.listener, err = net.ListenUDP("udp", addr)
	return &l, err
}

func (l *Listener) Serve() {
	log.Printf("Listening for messages on :%d...", l.port)

	// run loop forever (or until ctrl-c)
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

		go func() {
			err := l.handle(message)
			if err != nil {
				log.Printf("Handling failed: %s", err.Error())
			}
		}()
	}
}

func (l *Listener) handle(msg string) error {
	log.Println("Incoming message:", string(msg))

	in := Message{
		Timestamp: time.Now(),
	}
	err := json.Unmarshal([]byte(msg), &in)

	if err != nil {
		return err
	}

	switch in.Type {
	case inKill:
		km := KillMessage{}
		err := mapstructure.Decode(in.Data, &km)
		if err != nil {
			fmt.Println("Error: Could not decode mapstructure", err.Error())
		}
		return l.KillMessage(km)

	case inRoundStart:
		return nil
		// return l.StartRound()

	case inRoundEnd:
		return l.EndRound()

	case inMatchStart:
		return l.StartMatch()

	case inMatchEnd:
		return l.EndMatch()

	case inPickup:
		pm := PickupMessage{}
		err := mapstructure.Decode(in.Data, &pm)
		if err != nil {
			fmt.Println("Error: Could not decode mapstructure", err.Error())
		}
		return l.ArrowPickup(pm)

	default:
		log.Printf("Warning: Unknown message type '%s'", in.Type)
	}

	return nil
}

func (l *Listener) KillMessage(km KillMessage) error {
	t, err := l.DB.GetCurrentTournament()
	if err != nil {
		return err
	}

	return t.Matches[t.Current].KillMessage(km)
}

func (l *Listener) StartMatch() error {
	t, err := l.DB.GetCurrentTournament()
	if err != nil {
		return err
	}

	return t.Matches[t.Current].Start(nil)
}

func (l *Listener) EndMatch() error {
	t, err := l.DB.GetCurrentTournament()
	if err != nil {
		return err
	}

	return t.Matches[t.Current].End(nil)
}

func (l *Listener) EndRound() error {
	t, err := l.DB.GetCurrentTournament()
	if err != nil {
		return err
	}

	return t.Matches[t.Current].EndRound()
}

func (l *Listener) ArrowPickup(pm PickupMessage) error {
	t, err := l.DB.GetCurrentTournament()
	if err != nil {
		return err
	}

	return t.Matches[t.Current].ArrowPickup(pm)
}
