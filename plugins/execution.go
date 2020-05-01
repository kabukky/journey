// +build !noplugins

package plugins

import (
	"log"

	"github.com/rkuris/journey/structure"
	lua "github.com/yuin/gopher-lua"
)

// Execute ...
func Execute(helper *structure.Helper, values *structure.RequestData) ([]byte, error) {
	// Retrieve the lua state
	vm := values.PluginVMs[helper.Name]
	// Execute plugin
	err := vm.CallByParam(lua.P{Fn: vm.GetGlobal(helper.Name), NRet: 1, Protect: true})
	if err != nil {
		log.Println("Error while executing plugin for helper "+helper.Name+":", err)
		// Since the vm threw an error, close all vms and don't put the map back into the pool
		for _, luavm := range values.PluginVMs {
			luavm.Close()
		}
		values.PluginVMs = nil
		return []byte{}, err
	}
	// Get return value from vm
	ret := vm.ToString(-1)
	return []byte(ret), nil
}
