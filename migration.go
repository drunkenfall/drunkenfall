package main

import (
	"errors"
	"github.com/boltdb/bolt"
	"log"
	"os"
	"strconv"
)

// Migration is the current level pf patches that the database knows of.
// This should always be the total amount of migrations.
//
// If you add a new migration, increase this number to make sure that your
// migration is actually applied.
const Migration = 1

var (
	errNoMigrationsYet = errors.New("no migrations have been added yet")
	levelBucket        = []byte("level")
)

// migrations is a list of all the migration functions.
// Their indexes are used to determine when they are to be applied.
var migrations = []func(db *bolt.DB) error{
	InitialMigration,
}

// Migrate is the main migration entrypoint
//
// When called, it will check the database to see what migrations have already
// been applied. If that is lower than the const Migration, all migrations up
// to that point will sequentially be applied.
func Migrate(fn string) error {
	var lvl int

	db, err := bolt.Open(fn, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = db.View(func(tx *bolt.Tx) error {
		var err error
		b := tx.Bucket(MigrationKey)
		if b == nil {
			return errNoMigrationsYet
		}

		x := b.Get(levelBucket)
		lvl, err = strconv.Atoi(string(x))
		if err != nil {
			return err
		}

		return nil
	})

	if err == errNoMigrationsYet {
		// No migrations have been done yet, so lets do the first one by adding
		// the migration bucket.
		lvl = 0

	} else if err != nil {
		// Something actually went wrong. Oh noes.
		return err
	}

	return nil
}

// backup creates a backup of the database to be migrated
func backup(fn, name string) error {
	os.Mkdir("db-migration-backup/", 0755)
	return nil
}

// InitialMigration adds the db_version bucket with the sole entry
func InitialMigration(db *bolt.DB) error {
	if err := backup(db.Path(), "0-init"); err != nil {
		log.Printf("backup() failed: %s", err)
		return err
	}

	ret := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(MigrationKey)
		if err != nil {
			return err
		}
		b.Cursor()
		// b.Put(levelBucket)

		return nil
	})

	return ret
}
