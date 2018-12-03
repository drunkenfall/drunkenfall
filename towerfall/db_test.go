package towerfall

import (
	"log"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/go-pg/pg"
	"go.uber.org/zap"
)

func testDatabase(t *testing.T, c *Config) (*Database, func()) {
	t.Helper()

	pgdb := pg.Connect(&pg.Options{
		User:     "postgres",
		Database: "test_drunkenfall",
	})

	if c.DbVerbose {
		t.Log("DRUNKENFALL_DBVERBOSE is set; enabling logger")
		pgdb.OnQueryProcessed(func(event *pg.QueryProcessedEvent) {
			query, err := event.FormattedQuery()
			if err != nil {
				panic(err)
			}

			log.Printf("%s %s", time.Since(event.StartTime), query)
		})
	}

	// If this is set, will reuse the same database for all the tests.
	// This is significantly faster when running all of them since
	// database setup/teardown takes a while.
	if os.Getenv("DRUNKENFALL_DBFASTTEST") != "" {
		zlog, _ := zap.NewDevelopment()
		db := &Database{
			DB:  pgdb,
			log: zlog.With(zap.String("pgdb", "test_read")),
		}
		globalDB = db
		return db, func() {}
	}

	name := strings.ToLower(t.Name())
	name = regexp.MustCompile("[^_a-zA-Z0-9]+").ReplaceAllString(name, "_")

	_, err := pgdb.Exec("DROP DATABASE IF EXISTS " + name + "")
	if err != nil {
		log.Fatalf("couldn't create database '%s': %+v", name, err)
	}

	_, err = pgdb.Exec("CREATE DATABASE " + name + " WITH TEMPLATE test_drunkenfall OWNER postgres")
	if err != nil {
		log.Fatalf("couldn't create database '%s': %+v", name, err)
	}

	err = pgdb.Close()
	if err != nil {
		log.Fatalf("couldn't close database '%s': %+v", name, err)
	}

	pgdb = pg.Connect(&pg.Options{
		User:     "postgres",
		Database: name,
	})
	if err != nil {
		log.Fatalf("failed to open postgres session to '%s': %+v", name, err)
	}

	zlog, _ := zap.NewDevelopment()
	db := &Database{
		DB:  pgdb,
		log: zlog.With(zap.String("pg", "test_read")),
	}

	globalDB = db

	return db, func() {
		err = db.Close()
		if err != nil {
			log.Fatalf("couldn't close database '%s': %+v", name, err)
		}

		pgdb := pg.Connect(&pg.Options{
			User:     "postgres",
			Database: "test_drunkenfall",
		})

		_, err = pgdb.Exec("DROP DATABASE " + name + "")
		if err != nil {
			log.Fatalf("couldn't drop database '%s': %+v", name, err)
		}

		err = pgdb.Close()
		if err != nil {
			log.Fatalf("couldn't close database '%s': %+v", name, err)
		}
	}
}
