package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatabase_QueryConfig(t *testing.T) {
	database := PrepareTestDatabase()

	strategyId, time, err := database.QueryConfig("DEV")
	assert.NoError(t, err)

	t.Log(strategyId, time)
}
