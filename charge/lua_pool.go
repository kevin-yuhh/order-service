package charge

import (
	"sync"

	goLua "github.com/yuin/gopher-lua"
)

// Lua instance pool.
type lStatePool struct {
	script string
	m      sync.Mutex
	saved  []*goLua.LState
}

// Get lua state from pool.
func (pl *lStatePool) Get() *goLua.LState {
	pl.m.Lock()
	defer pl.m.Unlock()
	n := len(pl.saved)
	if n == 0 {
		return pl.New(pl.script)
	}
	x := pl.saved[n-1]
	pl.saved = pl.saved[0 : n-1]
	return x
}

// Reload lua state pool.
func (pl *lStatePool) Reload(script string) {
	pl.m.Lock()
	defer pl.m.Unlock()
	pl.script = script
	pl.Shutdown()
}

// New lua script.
func (pl *lStatePool) New(script string) *goLua.LState {
	L := goLua.NewState()
	_ = L.DoString(script)
	return L
}

// Put lua state into state pool.
func (pl *lStatePool) Put(L *goLua.LState) {
	pl.m.Lock()
	defer pl.m.Unlock()
	pl.saved = append(pl.saved, L)
}

// Shutdown lua state in pool.
func (pl *lStatePool) Shutdown() {
	for _, L := range pl.saved {
		L.Close()
	}
}
