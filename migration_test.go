package main

import (
	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"
)

func testDB(fn string) *bolt.DB {
	// Reset it
	_ = os.Remove(fn)

	db, err := bolt.Open(fn, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func TestInitialMigration(t *testing.T) {
	assert := assert.New(t)

	db := testDB("test/migration1.db")
	err := InitialMigration(db)
	assert.Nil(err)

	v, err := getVersion(db)
	assert.Nil(err)
	assert.Equal(1, v)
}

func TestBackup(t *testing.T) {
	assert := assert.New(t)
	path := "test/backups/"
	_ = os.RemoveAll(path)

	db := testDB("test/backup.db")

	err := backup(db, 255, path)
	assert.Nil(err)

	files, err := ioutil.ReadDir(path)
	assert.Nil(err)
	assert.Equal(1, len(files))

	assert.True(
		strings.Contains(files[0].Name(), "_"+strconv.Itoa(255)+"-"),
		"Backup file did not contain the version number",
	)
}
