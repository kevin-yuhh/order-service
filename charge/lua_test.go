package charge

import (
	"fmt"
	"sync"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

var (
	s1 = `
		function c1()
			return 1
		end
	`
	s2 = `
		function c1()
			return 2
		end
	`
)

func TestPrecompiled(t *testing.T) {
	p := &lStatePool{
		script: s1,
		saved:  make([]*lua.LState, 0, 4),
	}
	do(p, 1000)
	p.Reload(s2)
	do(p, 1000)
}

func do(p *lStatePool, times int) {
	var wg sync.WaitGroup
	wg.Add(times)
	for i := 0; i < times; i++ {
		go func(i int) {
			L := p.Get()
			if err := L.CallByParam(lua.P{
				Fn:      L.GetGlobal("c1"),
				NRet:    1,
				Protect: true,
			}, lua.LNumber(1000000000*i), lua.LNumber(10)); err != nil {
				panic(err)
			}
			ret := L.Get(-1)
			fmt.Println(ret.String())
			wg.Done()
		}(i)
	}
	wg.Wait()
}
