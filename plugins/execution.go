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
		vm.Close()
		return []byte{}, err
	}
	ret := vm.ToString(-1)
	log.Println("Successful!")
	LuaPool.Put(stateMap)
	return []byte(ret), nil
}
