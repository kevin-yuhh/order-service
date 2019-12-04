package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatabase_QueryConfig(t *testing.T) {
	database := PrepareTestDatabase()

	strategyId, time, err := database.QueryConfig("TEST")
	assert.NoError(t, err)

	t.Log(strategyId, time)
}
