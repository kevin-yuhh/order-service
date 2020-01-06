package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDatabase_QueryFileById(t *testing.T) {
	database := PrepareTestDatabase()

	// Query exists row.
	file, err := database.QueryFileById(1)
	assert.NoError(t, err)

	t.Log(file)
}

func TestDatabase_QueryFileByUk(t *testing.T) {
	database := PrepareTestDatabase()

	// Query exists row.
	file, err := database.QueryFileByUk(1, 1)
	assert.NoError(t, err)

	t.Log(file)
}

func TestDatabase_QueryMaxExpireByHash(t *testing.T) {
	database := PrepareTestDatabase()

	maxExpireTime, err := database.QueryMaxExpireByHash(1)
	assert.NoError(t, err)

	t.Log(maxExpireTime)
}

func TestInsertFileInfo(t *testing.T) {
	database := PrepareTestDatabase()

	session := database.DB.NewSession()
	err := session.Begin()
	assert.NoError(t, err)
	defer session.Close()

	id, err := InsertFileInfo(session, 2, 100, "hello.txt", int(time.Now().Local().Unix()))
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

func TestUpdateBtfsFileId(t *testing.T) {
	database := PrepareTestDatabase()

	session := database.DB.NewSession()
	err := session.Begin()
	assert.NoError(t, err)
	defer session.Close()

	err = UpdateBtfsFileId(session, 1, 2, 1)
	if err != nil {
		err1 := session.Rollback()
		assert.NoError(t, err1)
		t.Error(err)
		return
	}

	err = session.Commit()
	assert.NoError(t, err)
}

func TestUpdateFileExpireTime(t *testing.T) {
	database := PrepareTestDatabase()

	session := database.DB.NewSession()
	err := session.Begin()
	assert.NoError(t, err)
	defer session.Close()

	err = UpdateFileExpireTime(session, time.Now().Unix(), 1, 1)
	if err != nil {
		err1 := session.Rollback()
		assert.NoError(t, err1)
		t.Error(err)
		return
	}

	err = session.Commit()
	assert.NoError(t, err)
}

func TestDeleteFile(t *testing.T) {
	database := PrepareTestDatabase()

	session := database.DB.NewSession()
	err := session.Begin()
	assert.NoError(t, err)
	defer session.Close()

	err = DeleteFile(session, 3, 1)
	if err != nil {
		err1 := session.Rollback()
		assert.NoError(t, err1)
		t.Error(err)
		return
	}

	err = session.Commit()
	assert.NoError(t, err)
}

func TestReopenFile(t *testing.T) {
	database := PrepareTestDatabase()

	session := database.DB.NewSession()
	err := session.Begin()
	assert.NoError(t, err)
	defer session.Close()

	err = ReopenFile(session, 4, 1, 1578284537)
	if err != nil {
		err1 := session.Rollback()
		assert.NoError(t, err1)
		t.Error(err)
		return
	}

	err = session.Commit()
	assert.NoError(t, err)
}
