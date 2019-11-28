package model

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/TRON-US/soter-order-service/config"
	"github.com/TRON-US/soter-order-service/logger"
)

func TestDatabase_QueryLedgerInfoByAddress(t *testing.T) {
	conf, err := config.NewConfiguration("config", "..")
	if err != nil {
		logger.Logger.Fatal(err)
	}

	fmt.Println(conf)

	database, err := NewDatabase(conf)
	if err != nil {
		logger.Logger.Fatal(err)
	}

	ledger, err := database.QueryLedgerInfoByAddress("TCJCq2S7QuC5ijzdZBF2uLjg8z8fBtwZdS")
	fmt.Println(err)
	fmt.Println(ledger)
}

func TestDatabase_UpdateUserBalance(t *testing.T) {
	conf, err := config.NewConfiguration("config", "..")
	if err != nil {
		logger.Logger.Fatal(err)
	}

	fmt.Println(conf)

	database, err := NewDatabase(conf)
	if err != nil {
		logger.Logger.Fatal(err)
	}

	err = database.UpdateUserBalance(7000, 2000, 2, 1, "TCJCq2S7QuC5ijzdZBF2uLjg8z8fBtwZdS", 1574853071)
	assert.NoError(t, err)
}

func TestDatabase_UpdateLedgerInfo(t *testing.T) {
	conf, err := config.NewConfiguration("config", "..")
	if err != nil {
		logger.Logger.Fatal(err)
	}

	fmt.Println(conf)

	database, err := NewDatabase(conf)
	if err != nil {
		logger.Logger.Fatal(err)
	}

	err = database.UpdateLedgerInfo(100, 7000, 1000, 1000, 4, 1, "TCJCq2S7QuC5ijzdZBF2uLjg8z8fBtwZdS", 1574853071)
	assert.NoError(t, err)
}
