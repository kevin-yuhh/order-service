package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatabase_QueryActivityByUserId(t *testing.T) {
	database := PrepareTestDatabase()

	// Query exists row.
	strategy, err := database.QueryActivityByUserId(1)
	assert.NoError(t, err)

	t.Log(strategy)
}
