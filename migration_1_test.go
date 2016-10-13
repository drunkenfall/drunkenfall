package main

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func fatalError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func TestMigration1(t *testing.T) {
	assert := assert.New(t)

	db := migrationDB(1)
	err := MigrateOriginalColorPreferredColor(db)
	assert.Nil(err)

	v, err := getVersion(db)
	assert.Nil(err)
	assert.Equal(2, v)
}
