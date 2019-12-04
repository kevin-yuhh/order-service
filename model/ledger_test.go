package model

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDatabase_QueryLedgerInfoByAddress(t *testing.T) {
	database := PrepareTestDatabase()

	// Query exists row.
	ledger1, err := database.QueryLedgerInfoByAddress("TUsf2groYouQ7RzMkGcJH3PnSxFcwJCvrh")
	assert.NoError(t, err)

	// Query non-existing row.
	_, err = database.QueryLedgerInfoByAddress("TTCXimHXjen9BdTFW5JvcLKGWNm3SSuECF")
	assert.Error(t, err, errors.New("sql: no rows in result set"))

	t.Log(ledger1)
}

func TestUpdateUserBalance(t *testing.T) {
	database := PrepareTestDatabase()

	session := database.DB.NewSession()
	err := session.Begin()
	assert.NoError(t, err)
	defer session.Close()

	err = UpdateUserBalance(session, 9000, 2000, 1, 1, "TUsf2groYouQ7RzMkGcJH3PnSxFcwJCvrh", int(time.Now().Local().Unix()))
	if err != nil {
		err1 := session.Rollback()
		assert.NoError(t, err1)
		t.Error(err)
		return
	}

	err = session.Commit()
	assert.NoError(t, err)
}

func TestUpdateLedgerInfo(t *testing.T) {
	database := PrepareTestDatabase()

	session := database.DB.NewSession()
	err := session.Begin()
	assert.NoError(t, err)
	defer session.Close()

	err = UpdateLedgerInfo(session, 100, 9000, 1000, 1000, 2, 2, "TCJCq2S7QuC5ijzdZBF2uLjg8z8fBtwZdS", int(time.Now().Local().Unix()))
	if err != nil {
		err1 := session.Rollback()
		assert.NoError(t, err1)
		t.Error(err)
		return
	}

	err = session.Commit()
	assert.NoError(t, err)
}
