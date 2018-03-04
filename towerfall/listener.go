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
	inShot       = "arrow_shot"
	inShield     = "shield_state"
	inWings      = "wings_state"
	inOrbSlow    = "slow_orb_state"
	inOrbDark    = "dark_orb_state"
	inOrbLava    = "lava_orb_state"
	inOrbScroll  = "scroll_orb_state"
)

type KillMessage struct {
	Player int `json:"player"`
	Killer int `json:"killer"`
	Cause  int `json:"cause"`
}

type ArrowMessage struct {
	Player int    `json:"player"`
	Arrows Arrows `json:"arrows"`
}

type ShieldMessage struct {
	Player int  `json:"player"`
	State  bool `json:"state"`
}

type WingsMessage struct {
	Player int  `json:"player"`
	State  bool `json:"state"`
}

type SlowOrbMessage struct {
	State bool `json:"state"`
}

type DarkOrbMessage struct {
	State bool `json:"state"`
}

type ScrollOrbMessage struct {
	State bool `json:"state"`
}

type LavaOrbMessage struct {
	Player int  `json:"player"`
	State  bool `json:"state"`
}

// List of integers where one item is an arrow type as described in
// the arrow types above.
type Arrows []int

type StartRoundMessage struct {
	Arrows []Arrows `json:"arrows"`
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

	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}
	l.listener, err = net.ListenUDP("udp", addr)
	return &l, err
}

func (l *Listener) Serve() {
	log.Printf("Listening for messages on :%d...", l.port)

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

	t, err := l.DB.GetCurrentTournament()
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

		return t.Matches[t.Current].KillMessage(km)

	case inRoundStart:
		sr := StartRoundMessage{}
		err := mapstructure.Decode(in.Data, &sr)
		if err != nil {
			fmt.Println("Error: Could not decode mapstructure", err.Error())
		}
		return t.Matches[t.Current].StartRound(sr)

	case inRoundEnd:
		return t.Matches[t.Current].EndRound()

	case inMatchStart:
		return t.Matches[t.Current].Start(nil)

	case inMatchEnd:
		return t.Matches[t.Current].End(nil)

	case inPickup:
	case inShot:
		am := ArrowMessage{}
		err := mapstructure.Decode(in.Data, &am)
		if err != nil {
			fmt.Println("Error: Could not decode mapstructure", err.Error())
		}
		return t.Matches[t.Current].ArrowUpdate(am)

	case inShield:
		sm := ShieldMessage{}
		err := mapstructure.Decode(in.Data, &sm)
		if err != nil {
			fmt.Println("Error: Could not decode mapstructure", err.Error())
		}
		return t.Matches[t.Current].ShieldUpdate(sm)

	case inWings:
		wm := WingsMessage{}
		err := mapstructure.Decode(in.Data, &wm)
		if err != nil {
			fmt.Println("Error: Could not decode mapstructure", err.Error())
		}
		return t.Matches[t.Current].WingsUpdate(wm)

	default:
		log.Printf("Warning: Unknown message type '%s'", in.Type)
	}

	return nil
}
