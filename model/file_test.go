package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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

func TestUpdateFileHash(t *testing.T) {
	database := PrepareTestDatabase()

	session := database.DB.NewSession()
	err := session.Begin()
	assert.NoError(t, err)
	defer session.Close()

	err = UpdateFileHash(session, 3, "QmT78zSuBmuS4z925WZfrqQ1qHaJ56DQaTfyMUF7F8ff51")
	if err != nil {
		err1 := session.Rollback()
		assert.NoError(t, err1)
		t.Error(err)
		return
	}

	err = session.Commit()
	assert.NoError(t, err)
}
