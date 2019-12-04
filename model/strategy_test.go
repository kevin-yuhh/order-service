package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatabase_QueryStrategyById(t *testing.T) {
	database := PrepareTestDatabase()

	strategy, err := database.QueryStrategyById(1)
	assert.NoError(t, err)

	t.Log(strategy)
}
