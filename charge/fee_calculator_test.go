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
		{50, 1500, 2, 10},
		{1011111, 50, 3, 0},
		{1011111, 1500, 4, 20},
		{200000000, 50, 5, 1990},
		{200000000, 1500, 6, 2418},
		{1500000000, 1500, 7, 17521},
	}
	for _, tt := range tables {
		f := f.Fee(tt.size, tt.totalTimes, tt.time)
		assert.Equal(t, f, tt.output)
	}
}
