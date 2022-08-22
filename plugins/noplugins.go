//go:build noplugins
// +build noplugins

package plugins

import (
	"errors"
	"sync"

	"github.com/kabukky/journey/structure"
)

// Global LState pool
var LuaPool *lStatePool

type lStatePool struct {
	m     sync.Mutex
	files map[string]string
	saved []map[string]*string
}

// Load ...
func Load() error {
	LuaPool = nil
	return errors.New("Plugin system is not compiled")
}

// Execute ...
func Execute(helper *structure.Helper, values *structure.RequestData) ([]byte, error) {
	return []byte{}, nil
}

// Get ...
func (pl *lStatePool) Get(helper *structure.Helper, values *structure.RequestData) map[string]*string {
	return nil
}

// Put ...
func (pl *lStatePool) Put(L map[string]*string) {
}

// Shutdown ...
func (pl *lStatePool) Shutdown() {
}
