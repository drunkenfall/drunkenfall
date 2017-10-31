package towerfall

import (
	"sort"
	"time"
)

type CompleteSnapshot map[string]*Snapshot

type Snapshot struct {
	Person *Person         `json:"person"`
	Total  *PlayerSnapshot `json:"total"`
	Rank   int             `json:"rank"`
	// The Tournaments map contains the same data structure, but
	// separated per tournament
	Tournaments map[string]*PlayerSnapshot `json:"tournaments"`
}

// A PlayerSnapshot is a sum of all total stats about a player
type PlayerSnapshot struct {
	Shots    int           `json:"shots"`
	Sweeps   int           `json:"sweeps"`
	Kills    int           `json:"kills"`
	Self     int           `json:"self"`
	Matches  int           `json:"matches"`
	Rounds   int           `json:"rounds"`
	Score    int           `json:"score"`
	Playtime time.Duration `json:"playtime"`
	Wins     int           `json:"wins"`
}

// NewSnapshot returns a full snapshot
func NewSnapshot(s *Server) CompleteSnapshot {
	ss := make(map[string]*Snapshot)
	s.DB.LoadPeople()
	for _, p := range s.DB.People {
		ss[p.ID] = &Snapshot{
			Person:      p,
			Tournaments: make(map[string]*PlayerSnapshot),
		}
	}

	// Calculate the per-tournament data points
	for _, t := range s.DB.Tournaments {
		tid := t.ID
		for _, m := range t.Matches {
			for _, p := range m.Players {
				pid := p.Person.ID
				if _, ok := ss[pid].Tournaments[tid]; !ok {
					ss[pid].Tournaments[tid] = &PlayerSnapshot{}
				}

				ss[pid].Tournaments[tid].Matches++
				ss[pid].Tournaments[tid].Rounds += len(m.Rounds)
				ss[pid].Tournaments[tid].Shots += p.Shots
				ss[pid].Tournaments[tid].Sweeps += p.Sweeps
				ss[pid].Tournaments[tid].Kills += p.Kills
				ss[pid].Tournaments[tid].Self += p.Self
				ss[pid].Tournaments[tid].Score += p.Score()
				ss[pid].Tournaments[tid].Playtime += m.Duration()
			}

		}

		// Only do the winner calculations if someone actually won
		if !t.Ended.IsZero() && len(t.Winners) > 0 {
			winner := t.Winners[0].Person.ID
			ss[winner].Tournaments[tid].Wins++
		}
	}

	// Summarize the per-tournament data points into the totals
	for pid := range ss {
		if ss[pid].Total == nil {
			ss[pid].Total = &PlayerSnapshot{}
		}

		for tid := range ss[pid].Tournaments {
			t := ss[pid].Tournaments[tid]

			ss[pid].Total.Matches += t.Matches
			ss[pid].Total.Rounds += t.Rounds
			ss[pid].Total.Shots += t.Shots
			ss[pid].Total.Sweeps += t.Sweeps
			ss[pid].Total.Kills += t.Kills
			ss[pid].Total.Self += t.Self
			ss[pid].Total.Score += t.Score
			ss[pid].Total.Playtime += t.Playtime
		}
	}

	// Add the Rank attribute to all the players
	ranked := make([]*Snapshot, len(ss))
	x := 0
	for _, p := range ss {
		ranked[x] = p
		x++
	}
	for x, p := range SortByRank(ranked) {
		ss[p.Person.ID].Rank = x + 1 // The +1 fixes zero-index.
	}

	return CompleteSnapshot(ss)
}

// ByRank is a sort.Interface that sorts players by tournament wins
// and then by total score.
type ByRank []*Snapshot

func (b ByRank) Len() int {
	return len(b)
}

func (b ByRank) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b ByRank) Less(i, j int) bool {
	// If a person is not actively playing, they are definitely less.
	if b[i].Person.Disabled != b[j].Person.Disabled {
		return b[j].Person.Disabled
	}

	if b[i].Total.Wins == b[j].Total.Wins {
		return b[i].Total.Score > b[j].Total.Score
	}

	return b[i].Total.Wins > b[j].Total.Wins
}

func SortByRank(ps []*Snapshot) []*Snapshot {
	sort.Sort(ByRank(ps))
	return ps
}
