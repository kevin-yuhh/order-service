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

	id, err := InsertOrderInfo(session, 2, 100, 1, "14feb899-54b3-4025-8327-9b3c7168460c", 90)
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

func TestUpdateOrderInfo(t *testing.T) {
	database := PrepareTestDatabase()

	session := database.DB.NewSession()
	err := session.Begin()
	assert.NoError(t, err)
	defer session.Close()

	err = UpdateOrderInfo(session, 1, 2, "S")
	if err != nil {
		err1 := session.Rollback()
		assert.NoError(t, err1)
		t.Error(err)
		return
	}

	err = session.Commit()
	assert.NoError(t, err)
}
