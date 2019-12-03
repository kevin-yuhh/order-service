package model

import (
	"log"
	"os"
	"testing"

	"github.com/TRON-US/soter-order-service/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"gopkg.in/testfixtures.v2"
)

func TestMain(m *testing.M) {
	conf, err := config.NewConfiguration("config", "..")
	if err != nil {
		log.Fatal(err)
	}

	database, err := NewDatabase(conf)
	if err != nil {
		log.Fatal(err)
	}

	Fixtures, err = testfixtures.NewFolder(database.DB.DB().DB, &testfixtures.MySQL{}, "../fixtures")
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

func TestNewDatabase(t *testing.T) {
	conf, err := config.NewConfiguration("config", "..")
	assert.NoError(t, err)

	database, err := NewDatabase(conf)
	assert.NoError(t, err)

	err = database.DB.Ping()
	assert.NoError(t, err)
}
