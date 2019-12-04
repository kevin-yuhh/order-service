package charge

import (
	"fmt"
	"strconv"

	"github.com/TRON-US/soter-order-service/logger"

	lua "github.com/yuin/gopher-lua"
)

// Fee calculator struct.
type FeeCalculator struct {
	StrategyId int64
	script     string
	luaPool    *lStatePool
}

// Construct function of fee calculator.
func NewFeeCalculator(strategyId int64, script string) *FeeCalculator {
	return &FeeCalculator{StrategyId: strategyId, script: script, luaPool: &lStatePool{
		saved: make([]*lua.LState, 0, 10),
	}}
}

// Reload policy scripts and lua instances.
func (f *FeeCalculator) Reload() {
	f.luaPool.Reload(f.script)
}

// Calculate the cost, take an instance from the lua instance pool and call the calculation function.
func (f *FeeCalculator) Fee(size int64, totalTimes int, time int) int64 {
	l := f.luaPool.Get()
	if err := l.CallByParam(lua.P{
		Fn:      l.GetGlobal("calc"), // name of Lua function
		NRet:    1,                   // number of returned values
		Protect: true,                // return err or panic
	}, lua.LNumber(size), lua.LNumber(totalTimes), lua.LNumber(time)); err != nil {
		panic(err)
	}
	r, err := strconv.ParseFloat(l.Get(-1).String(), 64)
	if err != nil {
		errMessage := fmt.Sprintf("Parse string to float error, reasons: [%v]", err)
		logger.Logger.Errorw(errMessage, "function", "service.Fee")
	}
	l.Close()
	return int64(r)
}
