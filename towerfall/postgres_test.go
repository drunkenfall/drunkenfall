package towerfall

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

func testDatabase(t *testing.T, c *Config) (*Database, func()) {
	t.Helper()

	pg, err := gorm.Open("postgres", "user=postgres dbname=test_playground sslmode=disable")

	if c.DbVerbose {
		// t.Log("DRUNKENFALL_DBVERBOSE is set; enabling logger")
		pg.LogMode(true)
	}

	// If this is set, will reuse the same database for all the tests.
	// This is significantly faster when running all of them since
	// database setup/teardown takes a while.
	if os.Getenv("DRUNKENFALL_DBFASTTEST") != "" {
		zlog, _ := zap.NewDevelopment()
		db := &Database{
			DB:  pg,
			log: zlog.With(zap.String("db", "test_read")),
		}
		return db, func() {}
	}

	if err != nil {
		log.Fatalf("failed to open postgres session: %+v", err)
	}

	name := strings.ToLower(t.Name())
	name = regexp.MustCompile("[^_a-zA-Z0-9]+").ReplaceAllString(name, "_")

	err = pg.Exec("DROP DATABASE IF EXISTS " + name + "").Error
	if err != nil {
		log.Fatalf("couldn't create database '%s': %+v", name, err)
	}

	err = pg.Exec("CREATE DATABASE " + name + " WITH TEMPLATE test_drunkenfall OWNER postgres").Error
	errs := pg.GetErrors()
	if len(errs) != 0 {
		log.Fatal(errs)
	}
	if err != nil {
		log.Fatalf("couldn't create database '%s': %+v", name, err)
	}

	err = pg.Close()
	if err != nil {
		log.Fatalf("couldn't close database '%s': %+v", name, err)
	}

	str := fmt.Sprintf("user=postgres dbname=%s sslmode=disable", name)
	pg, err = gorm.Open("postgres", str)
	if err != nil {
		log.Fatalf("failed to open postgres session to '%s': %+v", name, err)
	}

	zlog, _ := zap.NewDevelopment()
	db := &Database{
		DB:  pg,
		log: zlog.With(zap.String("db", "test_read")),
	}

	return db, func() {
		err = pg.Close()
		if err != nil {
			log.Fatalf("couldn't close database '%s': %+v", name, err)
		}

		pg, err := gorm.Open("postgres", "user=postgres dbname=test_drunkenfall sslmode=disable")
		err = pg.Exec("DROP DATABASE " + name + "").Error
		if err != nil {
			log.Fatalf("couldn't drop database '%s': %+v", name, err)
		}

		err = pg.Close()
		if err != nil {
			log.Fatalf("couldn't close database '%s': %+v", name, err)
		}
	}
}
