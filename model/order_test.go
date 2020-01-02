package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsertOrderInfo(t *testing.T) {
	database := PrepareTestDatabase()

	session := database.DB.NewSession()
	err := session.Begin()
	assert.NoError(t, err)
	defer session.Close()

	id, err := InsertOrderInfo(session, 2, 2, 100, 1, "14feb899-54b3-4025-8327-9b3c7168460c", "R", "P", 90)
	if err != nil {
		err1 := session.Rollback()
		assert.NoError(t, err1)
		t.Error(err)
		return
	}

	err = session.Commit()
	assert.NoError(t, err)

	t.Log(id)
}

func TestUpdateOrderStatus(t *testing.T) {
	database := PrepareTestDatabase()

	session := database.DB.NewSession()
	err := session.Begin()
	assert.NoError(t, err)
	defer session.Close()

	err = UpdateOrderStatus(session, 1, "S")
	if err != nil {
		err1 := session.Rollback()
		assert.NoError(t, err1)
		t.Error(err)
		return
	}

	err = session.Commit()
	assert.NoError(t, err)
}

func TestUpdateOrderFileIdById(t *testing.T) {
	database := PrepareTestDatabase()

	session := database.DB.NewSession()
	err := session.Begin()
	assert.NoError(t, err)
	defer session.Close()

	err = UpdateOrderFileIdById(session, 2, 1)
	if err != nil {
		err1 := session.Rollback()
		assert.NoError(t, err1)
		t.Error(err)
		return
	}

	err = session.Commit()
	assert.NoError(t, err)
}

func TestDatabase_QueryOrderInfoById(t *testing.T) {
	database := PrepareTestDatabase()

	// Query exists row.
	order, err := database.QueryOrderInfoById(1)
	assert.NoError(t, err)

	t.Log(order)
}

func TestDatabase_QueryOrderInfoByRequestId(t *testing.T) {
	database := PrepareTestDatabase()

	// Query exists row.
	order, err := database.QueryOrderInfoByRequestId("322ed182-187c-4b87-8dd3-e309c8f9cbe0", "TUsf2groYouQ7RzMkGcJH3PnSxFcwJCvrh")
	assert.NoError(t, err)

	t.Log(order)
}
