package plugins

import (
	"github.com/kabukky/journey/structure"
	"github.com/yuin/gopher-lua"
	"log"
)

func Execute(stateMap map[string]*lua.LState, name string, values *structure.RequestData) ([]byte, error) {
	// Retrieve a lua state
	vm := stateMap[name]
	// Execute plugin
	err := vm.CallByParam(lua.P{Fn: vm.GetGlobal(name), NRet: 1, Protect: true})
	if err != nil {
		log.Println("Error while exec:", err)
		// SInce the vm threw an error, close it and don't put the map back into the pool
		vm.Close()
		return []byte{}, err
	}
	ret := vm.ToString(-1)
	log.Println("Successful!")
	// Put the map back into the pool
	LuaPool.Put(stateMap)
	return []byte(ret), nil
}
