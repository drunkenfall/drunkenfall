package main

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"log"
	"sort"
	"time"
)

type mig3prevPlayer struct {
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

type mig3prevMatch struct {
	Players    []mig3prevPlayer `json:"players"`
	Judges     []Judge          `json:"judges"`
	Kind       string           `json:"kind"`
	Index      int              `json:"index"`
	Length     int              `json:"length"`
	Pause      time.Duration    `json:"pause"`
	Scheduled  time.Time        `json:"scheduled"`
	Started    time.Time        `json:"started"`
	Ended      time.Time        `json:"ended"`
	Tournament *Tournament      `json:"-"`
	ScoreOrder []int            `json:"score_order"`
	Commits    []MatchCommit    `json:"commits"`
}

type mig3prevTournament struct {
	Name      string           `json:"name"`
	ID        string           `json:"id"`
	Players   []mig3prevPlayer `json:"players"`
	Winners   []mig3prevPlayer `json:"winners"`
	Runnerups []*Person        `json:"runnerups"`
	Judges    []Judge          `json:"judges"`
	Tryouts   []*mig3prevMatch `json:"tryouts"`
	Semis     []*mig3prevMatch `json:"semis"`
	Final     *mig3prevMatch   `json:"final"`
	Current   CurrentMatch     `json:"current"`
	Opened    time.Time        `json:"opened"`
	Scheduled time.Time        `json:"scheduled"`
	Started   time.Time        `json:"started"`
	Ended     time.Time        `json:"ended"`
}

// mig3ByKills is a sort.Interface that sorts mig3prevPlayers by their kills
type mig3ByKills []mig3prevPlayer

func (s mig3ByKills) Len() int {
	return len(s)

}
func (s mig3ByKills) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]

}
func (s mig3ByKills) Less(i, j int) bool {
	// Technically not Less, but we want biggest first...
	return s[i].Kills > s[j].Kills
}

// mig3SortByKills returns a list in order of the kills the mig3prevPlayers have
func mig3SortByKills(ps []mig3prevPlayer) []mig3prevPlayer {
	tmp := make([]mig3prevPlayer, len(ps))
	copy(tmp, ps)
	sort.Sort(mig3ByKills(tmp))
	return tmp
}

// MigrateMatchScoreOrder updates the ScoreOrder since the implementation had
// a bugfix. Previously it sorted by the entertainment score and not the
// amount of kills. This has been fixed and all the scores need to be updated
// to reflect this in previous tournaments.
func MigrateMatchScoreOrder(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		// Load the tournaments
		var ts []*mig3prevTournament
		err := db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket(TournamentKey)
			err := b.ForEach(func(k []byte, v []byte) error {
				t := &mig3prevTournament{}
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

		var out []*mig3prevTournament
		for _, pt := range ts {
			t := mig3prevTournament{
				Name:      pt.Name,
				ID:        pt.ID,
				Players:   pt.Players,
				Winners:   pt.Winners,
				Runnerups: pt.Runnerups,
				Judges:    pt.Judges,
				Tryouts:   pt.Tryouts,
				Semis:     pt.Semis,
				Final:     pt.Final,
				Current:   pt.Current,
				Opened:    pt.Opened,
				Scheduled: pt.Scheduled,
				Started:   pt.Started,
				Ended:     pt.Ended,
			}

			for x := range t.Tryouts {
				t.Tryouts[x].ScoreOrder = make([]int, 0)
				ps := mig3SortByKills(t.Tryouts[x].Players)
				for _, p := range ps {
					for i, o := range t.Tryouts[x].Players {
						if p.Person.ID == o.Person.ID {
							t.Tryouts[x].ScoreOrder = append(t.Tryouts[x].ScoreOrder, i)
							break
						}
					}
				}
			}

			for x := range t.Semis {
				t.Semis[x].ScoreOrder = make([]int, 0)
				ps := mig3SortByKills(t.Semis[x].Players)
				for _, p := range ps {
					for i, o := range t.Semis[x].Players {
						if p.Person.ID == o.Person.ID {
							t.Semis[x].ScoreOrder = append(t.Semis[x].ScoreOrder, i)
							break
						}
					}
				}
			}

			t.Final.ScoreOrder = make([]int, 0)
			ps := mig3SortByKills(t.Final.Players)
			for _, p := range ps {
				for i, o := range t.Final.Players {
					if p.Person.ID == o.Person.ID {
						t.Final.ScoreOrder = append(t.Final.ScoreOrder, i)
						break
					}
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

		return setVersion(tx, 4)
	})
}
