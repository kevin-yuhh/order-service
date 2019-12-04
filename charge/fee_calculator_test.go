package charge

import (
	"testing"

	"github.com/TRON-US/soter-order-service/config"
	"github.com/TRON-US/soter-order-service/model"

	"github.com/stretchr/testify/assert"
)

func TestCalc(t *testing.T) {
	conf, err := config.NewConfiguration("config", "..")
	assert.NoError(t, err)

	database, err := model.NewDatabase(conf)
	assert.NoError(t, err)

	// Get configure from database by env.
	strategyId, _, err := database.QueryConfig(conf.Env)
	assert.NoError(t, err)

	// Get strategy from database by strategy id.
	strategy, err := database.QueryStrategyById(strategyId)
	assert.NoError(t, err)

	f := NewFeeCalculator(strategyId, strategy.Script)
	f.Reload()

	tables := []struct {
		size       int64
		totalTimes int
		time       int
		output     int64
	}{
		{50, 50, 1, 0},
		{50, 1500, 2, 4},
		{1011111, 50, 3, 6},
		{1011111, 1500, 4, 12},
		{200000000, 50, 5, 2000},
		{200000000, 1500, 6, 2406},
		{1500000000, 1500, 7, 10507},
	}
	for _, tt := range tables {
		f := f.Fee(tt.size, tt.totalTimes, tt.time)
		assert.Equal(t, f, tt.output)
	}
}
