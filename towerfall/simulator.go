package towerfall

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"github.com/streadway/amqp"
)

var defaultSleep = time.Millisecond * 650

type Simulator struct {
	Tournament *Tournament
	enabled    bool
	DB         *Database
	conn       *amqp.Connection
	ch         *amqp.Channel
	q          amqp.Queue
	sleep      time.Duration
}

const dogebutt = "/static/img/dogebutt.png"

func NewSimulator(s *Server) (*Simulator, error) {
	sim := &Simulator{
		sleep: defaultSleep,
		DB:    s.DB,
	}
	return sim, nil
}

func (s *Simulator) Connect() error {
	// var err error
	// s.conn, err = amqp.Dial()
	// if err != nil {
	// 	return err
	// }

	// s.ch, err = s.conn.Channel()
	// if err != nil {
	// 	return err
	// }

	// s.q, err = s.ch.QueueDeclare(
	// 	s.conf.RabbitIncomingQueue, // name
	// 	false, // durable
	// 	false, // delete when unused
	// 	false, // exclusive
	// 	false, // no-wait
	// 	nil,   // arguments
	// )
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (s *Simulator) Start(tid string) {
	if s.Tournament == nil {
		s.Tournament = s.DB.tournamentRef[tid]
	}

	log.Printf("Simulation starting for %s...", tid)
	s.enabled = true
	go s.Serve()
}

func (s *Simulator) Stop() {
	log.Print("Simulation stopped...")
	s.enabled = false
}

// Serve starts the loop and starts doing things
func (s *Simulator) Serve() {
	var sleep time.Duration
	for {
		if !s.enabled {
			log.Print("Simulator stopped; halting execution.")
			return
		}

		sleep = s.Action()
		time.Sleep(sleep)
	}
}

// Action is the main execution point of the simulator
func (s *Simulator) Action() time.Duration {
	m := s.Tournament.Matches[s.Tournament.Current]

	// If the match isn't started, the first thing we need to do is to
	// start it!
	if !m.IsStarted() {
		s.send(inMatchStart)
		s.send(inRoundStart)
		return s.sleep * 3
	}

	// If there is only one player left alive then we should end the
	// current round. Match end will be handled by the app automatically.
	if len(s.alivePlayers()) <= 1 {
		s.send(inRoundEnd)
		time.Sleep(s.sleep * 4)
		s.send(inRoundStart, StartRoundMessage{[]Arrows{
			[]int{aNormal, aNormal, aNormal},
			[]int{aNormal, aNormal, aNormal},
			[]int{aNormal, aNormal, aNormal},
			[]int{aNormal, aNormal, aNormal},
		}})
		return s.sleep
	}

	p := s.randomPlayer()
	// 5% of the times - environment suicide
	// 15% of the times - kill a player (yourself included)
	// Otherwise - how do i shot arw
	y := rand.Intn(100)
	if y <= 5 {
		s.envkill(p)
	} else if y <= 50 {
		s.kill(p, s.otherPlayer(p))
	} else {
		s.arrow(p)
	}

	return s.sleep
}

// alivePlayers returns a list of the indexes of alive players
func (s *Simulator) alivePlayers() (ret []int) {
	for i, p := range s.Tournament.Matches[s.Tournament.Current].Players {
		if p.State.Alive {
			ret = append(ret, i)
		}
	}
	return ret
}

// randomPlayer returns the index of a random player that is still alive.
func (s *Simulator) randomPlayer() int {
	ps := s.alivePlayers()
	rand.Seed(time.Now().UnixNano())
	return ps[rand.Intn(len(ps))]
}

// randomPlayer returns the index of a random player that is still alive.
func (s *Simulator) otherPlayer(p int) (ret int) {
	ps := s.alivePlayers()
	rand.Seed(time.Now().UnixNano())
	for {
		ret = ps[rand.Intn(len(ps))]
		if ret != p {
			return ret
		}
	}
}

// envkill simulates being killed by the environment
func (s *Simulator) envkill(p int) {
	reasons := []int{rSpikeBall, rFallingObject, rSquish, rMiasma}
	reason := reasons[rand.Intn(len(reasons))]
	s.send(inKill, KillMessage{p, p, reason})
}

// kill simulates one player killing another
func (s *Simulator) kill(k, p int) {
	reasons := []int{
		rArrow,
		rExplosion,
		rBrambles,
		rJumpedOn,
		rLava,
		rShock,
		rFallingObject,
		rSquish,
	}

	if k == p {
		// If this is a suicide, we need to add curse, and we add a lot of
		// them to make that more common.
		for x := 0; x <= 8; x++ {
			reasons = append(reasons, rCurse)
		}
	}

	reason := reasons[rand.Intn(len(reasons))]
	s.send(inKill, KillMessage{k, p, reason})
}

// arrow simulates either shooting an arrow or picking one up
func (s *Simulator) arrow(p int) {
	arrows := []int{
		aNormal,
		aBomb,
		aSuperBomb,
		aLaser,
		aBramble,
		aDrill,
		aBolt,
		aFeather,
		aTrigger,
		aPrism,
	}

	y := rand.Intn(100)
	a := s.Tournament.Matches[s.Tournament.Current].Players[p].State.Arrows

	// 30% of the time - pick up arrows
	if y <= 30 {
		arrow := arrows[rand.Intn(len(arrows))]
		a = append(Arrows{arrow}, a...)
		// From chests you get two arrows at once, so half of the time we
		// add two arrows rather than one.
		if y <= 15 {
			a = append(Arrows{arrow}, a...)
		}

		s.send(inPickup, ArrowMessage{p, a})

	} else {
		// If we're not picking up arrows, then we shoot one, if possible!
		if len(a) > 0 {
			a = a[1:len(a)]
			s.send(inShot, ArrowMessage{p, a})
		} else {
			// Curse
			s.send(inKill, KillMessage{p, p, rCurse})
		}
	}

}

// send sends a message to the listener
func (s *Simulator) send(t string, m ...interface{}) error {
	msg := Message{Type: t}
	if len(m) == 1 {
		msg.Data = m[0]
	}
	body, err := json.Marshal(msg)

	if err != nil {
		return err
	}

	err = s.ch.Publish(
		"",       // exchange
		s.q.Name, // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		})

	return err
}
