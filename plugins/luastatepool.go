package plugins

import (
	"github.com/yuin/gopher-lua"
	"log"
	"sync"
)

// Global lStateMap pool
var LuaPool *lStatePool

// Struct to hold the map of plugin helpers
type lStatePool struct {
	sync.Mutex
	template map[string]*lua.LState // key (string) = name of helper, value (*lua.LState) = pointer to the vm with the lua script to execute the helper
	saved    []map[string]*lua.LState
}

// This function will only be called when a .lua script was found in the plugins directory.
func makeLuaPool() {
	LuaPool = &lStatePool{
		saved: make([]map[string]*lua.LState, 0, 4),
	}
}

func (pl *lStatePool) Get() map[string]*lua.LState {
	pl.Lock()
	defer pl.Unlock()
	n := len(pl.saved)
	if n == 0 {
		return pl.New()
	}
	stateMap := pl.saved[n-1]
	pl.saved = pl.saved[0 : n-1]
	return stateMap
}

func (pl *lStatePool) New() map[string]*lua.LState {
	stateMap := make(map[string]*lua.LState, 0)
	// If the state map is already created, copy the whole map to the new one
	if LuaPool.template != nil {
		wasAlreadyAssigned := false
		for key, value := range LuaPool.template {
			// Loop a second time to see if the pointer was already used in this assignment
			for key2, value2 := range LuaPool.template {
				if value2 == value {
					// Check if key2 was already assigned. If so, use that pointer
					if stateMap[key2] != nil {
						log.Println("Duplicate pointer:", key, key2)
						stateMap[key] = stateMap[key2]
						wasAlreadyAssigned = true
						break
					}
				}
			}
			// Assign a copy of the lua.LState struct by dereferencing it.
			if !wasAlreadyAssigned {
				vm := *value
				stateMap[key] = &vm
			}
		}
	}
	return stateMap
}

func (pl *lStatePool) Put(stateMap map[string]*lua.LState) {
	pl.Lock()
	defer pl.Unlock()
	pl.saved = append(pl.saved, stateMap)
}

func (pl *lStatePool) Shutdown() {
	for _, stateMap := range pl.saved {
		for _, vm := range stateMap {
			vm.Close()
		}
	}
}

func (pl *lStatePool) Size() int {
	return len(pl.saved)
}
