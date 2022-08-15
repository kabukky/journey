//go:build !noplugins
// +build !noplugins

package plugins

import (
	"sync"

	"github.com/rkuris/journey/structure"
	lua "github.com/yuin/gopher-lua"
)

// LuaPool is the Global LState pool
var LuaPool *lStatePool

type lStatePool struct {
	m     sync.Mutex
	files map[string]string
	saved []map[string]*lua.LState
}

// Get ...
func (pl *lStatePool) Get(helper *structure.Helper, values *structure.RequestData) map[string]*lua.LState {
	pl.m.Lock()
	defer pl.m.Unlock()
	n := len(pl.saved)
	if n == 0 {
		x := pl.New()
		// Since these are new lua states, do the lua file.
		for key, value := range x {
			setUpVM(value, helper, values, LuaPool.files[key])
			value.DoFile(LuaPool.files[key])
		}
		return x
	}
	x := pl.saved[n-1]
	// Set the new values for this request in every lua state
	for key, value := range x {
		setUpVM(value, helper, values, LuaPool.files[key])
	}
	pl.saved = pl.saved[0 : n-1]
	return x
}

// New ...
func (pl *lStatePool) New() map[string]*lua.LState {
	stateMap := make(map[string]*lua.LState, 0)
	for key := range LuaPool.files {
		L := lua.NewState()
		stateMap[key] = L
	}
	return stateMap
}

// Put ...
func (pl *lStatePool) Put(L map[string]*lua.LState) {
	pl.m.Lock()
	defer pl.m.Unlock()
	pl.saved = append(pl.saved, L)
}

// Shutdown ...
func (pl *lStatePool) Shutdown() {
	for _, stateMap := range pl.saved {
		for _, value := range stateMap {
			value.Close()
		}
	}
}

func newLuaPool() *lStatePool {
	return &lStatePool{saved: make([]map[string]*lua.LState, 0, 4)}
}
