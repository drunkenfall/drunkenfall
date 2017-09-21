package migrations

import (
	"log"

	"github.com/boltdb/bolt"
	"github.com/drunkenfall/drunkenfall/towerfall"
)

func FlattenMatches(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		out, err := Convert(db, "flatten")

		if err != nil {
			return err
		}

		// Loop the tournament results and save them into the db
		for id, t := range out {
			err = tx.Bucket(towerfall.TournamentKey).Put([]byte(id), t)
			if err != nil {
				log.Fatal(err)
			}
		}
		return setVersion(tx, 6)
	})
}
