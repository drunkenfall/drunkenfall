package migrations

import "github.com/boltdb/bolt"

// These are all here just to be placeholders for the struct-based
// migrations. They worked well but were cumbersome to work with. They
// also should never have to be run again, so we should be good.

func MigrateOriginalColorPreferredColor(db *bolt.DB) error    { return nil }
func MigrateTournamentRunnerupStringPerson(db *bolt.DB) error { return nil }
func MigrateMatchScoreOrderKillOrder(db *bolt.DB) error       { return nil }
func MigrateMatchCommitToRound(db *bolt.DB) error             { return nil }
