package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsertBtfsFileInfo(t *testing.T) {
	database := PrepareTestDatabase()

	session := database.DB.NewSession()
	err := session.Begin()
	assert.NoError(t, err)
	defer session.Close()

	id, err := InsertBtfsFileInfo(session, "QmQUNM9ggZXNmuZXYxN2VAbDbmjRiRS9mJYX1ZDgvW8piW", 1597881600)
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

func TestDatabase_QueryBtfsFileByHash(t *testing.T) {
	database := PrepareTestDatabase()

	// Query exists row.
	file, err := database.QueryBtfsFileByHash("QmT78zSuBmuS4z925WZfrqQ1qHaJ56DQaTfyMUF7F8ff5o")
	assert.NoError(t, err)

	t.Log(file)
}

func TestUpdateBtfsFileExpireTime(t *testing.T) {
	database := PrepareTestDatabase()

	session := database.DB.NewSession()
	err := session.Begin()
	assert.NoError(t, err)
	defer session.Close()

	err = UpdateBtfsFileExpireTime(session, 1, 1, 1597881600)
	if err != nil {
		err1 := session.Rollback()
		assert.NoError(t, err1)
		t.Error(err)
		return
	}

	err = session.Commit()
	assert.NoError(t, err)
}
