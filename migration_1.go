package main

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"log"
	"time"
)

// MigrateOriginalColorPreferredColor changes Player{} to have a preferred
// color rather than an original color
func MigrateOriginalColorPreferredColor(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		type previousPlayer struct {
			Person        *Person `json:"person"`
			Color         string  `json:"color"`
			OriginalColor string  `json:"original_color"`
			Shots         int     `json:"shots"`
			Sweeps        int     `json:"sweeps"`
			Kills         int     `json:"kills"`
			Self          int     `json:"self"`
			Explosions    int     `json:"explosions"`
			Matches       int     `json:"matches"`
			TotalScore    int     `json:"score"`
			Match         *Match  `json:"-"`
		}

		type previousMatch struct {
			Players    []previousPlayer `json:"players"`
			Judges     []Judge          `json:"judges"`
			Kind       string           `json:"kind"`
			Index      int              `json:"index"`
			Length     int              `json:"length"`
			Pause      time.Duration    `json:"pause"`
			Scheduled  time.Time        `json:"scheduled"`
			Started    time.Time        `json:"started"`
			Ended      time.Time        `json:"ended"`
			Tournament *Tournament      `json:"-"`
		}

		type previousTournament struct {
			Name        string           `json:"name"`
			ID          string           `json:"id"`
			Players     []previousPlayer `json:"players"`
			Winners     []previousPlayer `json:"winners"`
			Runnerups   []string         `json:"runnerups"`
			Judges      []Judge          `json:"judges"`
			Tryouts     []*previousMatch `json:"tryouts"`
			Semis       []*previousMatch `json:"semis"`
			Final       *previousMatch   `json:"final"`
			Current     CurrentMatch     `json:"current"`
			Opened      time.Time        `json:"opened"`
			Scheduled   time.Time        `json:"scheduled"`
			Started     time.Time        `json:"started"`
			Ended       time.Time        `json:"ended"`
			db          *Database
			server      *Server
			length      int
			finalLength int
		}

		// Load the tournaments
		var ts []*previousTournament
		err := db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket(TournamentKey)
			err := b.ForEach(func(k []byte, v []byte) error {
				t := &previousTournament{}
				err := json.Unmarshal(v, t)
				if err != nil {
					log.Print(err)
					return err
				}
				if err != nil {
					return err
				}

				ts = append(ts, t)
				return nil
			})
			return err
		})

		if err != nil {
			log.Fatal(err)
		}

		type currentPlayer struct {
			Person         *Person `json:"person"`
			Color          string  `json:"color"`
			PreferredColor string  `json:"preferred_color"`
			Shots          int     `json:"shots"`
			Sweeps         int     `json:"sweeps"`
			Kills          int     `json:"kills"`
			Self           int     `json:"self"`
			Explosions     int     `json:"explosions"`
			Matches        int     `json:"matches"`
			TotalScore     int     `json:"score"`
			Match          *Match  `json:"-"`
		}

		type currentMatch struct {
			Players    []currentPlayer `json:"players"`
			Judges     []Judge         `json:"judges"`
			Kind       string          `json:"kind"`
			Index      int             `json:"index"`
			Length     int             `json:"length"`
			Pause      time.Duration   `json:"pause"`
			Scheduled  time.Time       `json:"scheduled"`
			Started    time.Time       `json:"started"`
			Ended      time.Time       `json:"ended"`
			Tournament *Tournament     `json:"-"`
		}

		type currentTournament struct {
			Name        string          `json:"name"`
			ID          string          `json:"id"`
			Players     []currentPlayer `json:"players"`
			Winners     []currentPlayer `json:"winners"`
			Runnerups   []*Person       `json:"runnerups"`
			Judges      []Judge         `json:"judges"`
			Tryouts     []*currentMatch `json:"tryouts"`
			Semis       []*currentMatch `json:"semis"`
			Final       *currentMatch   `json:"final"`
			Current     CurrentMatch    `json:"current"`
			Opened      time.Time       `json:"opened"`
			Scheduled   time.Time       `json:"scheduled"`
			Started     time.Time       `json:"started"`
			Ended       time.Time       `json:"ended"`
			db          *Database
			server      *Server
			length      int
			finalLength int
		}

		var out []*currentTournament
		for _, pt := range ts {
			t := currentTournament{
				Name:      pt.Name,
				ID:        pt.ID,
				Judges:    pt.Judges,
				Current:   pt.Current,
				Opened:    pt.Opened,
				Scheduled: pt.Scheduled,
				Started:   pt.Started,
				Ended:     pt.Ended,
			}

			// Update the player objects
			t.Players = make([]currentPlayer, len(pt.Players))
			for x, pp := range pt.Players {
				t.Players[x] = currentPlayer{
					Person:         pp.Person,
					Color:          pp.Color,
					PreferredColor: pp.OriginalColor,
					Shots:          pp.Shots,
					Sweeps:         pp.Sweeps,
					Kills:          pp.Kills,
					Self:           pp.Self,
					Explosions:     pp.Explosions,
					Matches:        pp.Matches,
					TotalScore:     pp.TotalScore,
				}
			}

			// Update the tryouts
			t.Tryouts = make([]*currentMatch, len(pt.Tryouts))
			for x, pm := range pt.Tryouts {
				t.Tryouts[x] = &currentMatch{
					Judges:     pm.Judges,
					Kind:       pm.Kind,
					Index:      pm.Index,
					Length:     pm.Length,
					Pause:      pm.Pause,
					Scheduled:  pm.Scheduled,
					Started:    pm.Started,
					Ended:      pm.Ended,
					Tournament: pm.Tournament,
				}

				// For each match, also update each of the player objects
				t.Tryouts[x].Players = make([]currentPlayer, len(pt.Tryouts[x].Players))
				for y, pp := range pt.Tryouts[x].Players {
					t.Tryouts[x].Players[y] = currentPlayer{
						Person:         pp.Person,
						Color:          pp.Color,
						PreferredColor: pp.OriginalColor,
						Shots:          pp.Shots,
						Sweeps:         pp.Sweeps,
						Kills:          pp.Kills,
						Self:           pp.Self,
						Explosions:     pp.Explosions,
						Matches:        pp.Matches,
						TotalScore:     pp.TotalScore,
					}
				}
			}

			// Update the semis
			t.Semis = make([]*currentMatch, len(pt.Semis))
			for x, pm := range pt.Semis {
				t.Semis[x] = &currentMatch{
					Judges:     pm.Judges,
					Kind:       pm.Kind,
					Index:      pm.Index,
					Length:     pm.Length,
					Pause:      pm.Pause,
					Scheduled:  pm.Scheduled,
					Started:    pm.Started,
					Ended:      pm.Ended,
					Tournament: pm.Tournament,
				}

				// For each match, also update each of the player objects
				t.Semis[x].Players = make([]currentPlayer, len(pt.Semis[x].Players))
				for y, pp := range pt.Semis[x].Players {
					t.Semis[x].Players[y] = currentPlayer{
						Person:         pp.Person,
						Color:          pp.Color,
						PreferredColor: pp.OriginalColor,
						Shots:          pp.Shots,
						Sweeps:         pp.Sweeps,
						Kills:          pp.Kills,
						Self:           pp.Self,
						Explosions:     pp.Explosions,
						Matches:        pp.Matches,
						TotalScore:     pp.TotalScore,
					}
				}
			}

			// Update the final
			t.Final = &currentMatch{
				Judges:     pt.Final.Judges,
				Kind:       pt.Final.Kind,
				Index:      pt.Final.Index,
				Length:     pt.Final.Length,
				Pause:      pt.Final.Pause,
				Scheduled:  pt.Final.Scheduled,
				Started:    pt.Final.Started,
				Ended:      pt.Final.Ended,
				Tournament: pt.Final.Tournament,
			}

			// For each match, also update each of the player objects
			for y, pp := range pt.Final.Players {
				t.Final.Players[y] = currentPlayer{
					Person:         pp.Person,
					Color:          pp.Color,
					PreferredColor: pp.OriginalColor,
					Shots:          pp.Shots,
					Sweeps:         pp.Sweeps,
					Kills:          pp.Kills,
					Self:           pp.Self,
					Explosions:     pp.Explosions,
					Matches:        pp.Matches,
					TotalScore:     pp.TotalScore,
				}
			}

			out = append(out, &t)
		}

		// Loop the tournament results and save them into the db
		for _, t := range out {
			j, err := json.Marshal(t)
			if err != nil {
				log.Fatal(err)
			}

			err = tx.Bucket(TournamentKey).Put([]byte(t.ID), j)
			if err != nil {
				log.Fatal(err)
			}
		}

		return setVersion(tx, 2)
	})
}
