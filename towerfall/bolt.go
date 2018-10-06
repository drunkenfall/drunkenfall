package towerfall

import (
	"errors"
	"sort"
	"sync"

	"github.com/boltdb/bolt"
)

type boltWriter struct {
	DB *bolt.DB
}

type boltReader struct {
	DB *bolt.DB
}

var (
	// TournamentKey defines the tournament buckets
	TournamentKey = []byte("tournaments")
	// PeopleKey defines the bucket of people
	PeopleKey = []byte("people")
	// MigrationKey defines the bucket of migration levels
	MigrationKey = []byte("migration")
)

var tournamentMutex = &sync.Mutex{}
var personMutex = &sync.Mutex{}

// Used to signal that a current tournament was found and that the
// scanner should stop iterating.
var ErrTournamentFound = errors.New("found")

// ByScheduleDate is a sort.Interface that sorts tournaments according
// to when they were scheduled.
type ByScheduleDate []*Tournament

func (s ByScheduleDate) Len() int {
	return len(s)
}
func (s ByScheduleDate) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]

}
func (s ByScheduleDate) Less(i, j int) bool {
	return s[i].Scheduled.Before(s[j].Scheduled)
}

// SortByScheduleDate returns a list in order of schedule date
func SortByScheduleDate(ps []*Tournament) []*Tournament {
	sort.Sort(ByScheduleDate(ps))
	return ps
}
