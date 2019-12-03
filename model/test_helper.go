package model

import (
	"database/sql"
	"log"

	"github.com/TRON-US/soter-order-service/config"

	"gopkg.in/testfixtures.v2"
)

var (
	Db       *sql.DB
	Fixtures *testfixtures.Context
)

// Get test database connection.
func PrepareTestDatabase() *Database {
	if err := Fixtures.Load(); err != nil {
		log.Fatal(err)
	}

	conf, err := config.NewConfiguration("config", "..")
	if err != nil {
		log.Fatal(err)
	}

	database, err := NewDatabase(conf)
	if err != nil {
		log.Fatal(err)
	}

	return database
}
