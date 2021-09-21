//go:build noplugins
// +build noplugins

package plugins

import (
	"errors"
	"github.com/Landria/journey/structure"
	"sync"
)

// Global LState pool
var LuaPool *lStatePool

type lStatePool struct {
	m     sync.Mutex
	files map[string]string
	saved []map[string]*string
}

func Load() error {
	LuaPool = nil
	return errors.New("Plugin system is not compiled")
}

func Execute(helper *structure.Helper, values *structure.RequestData) ([]byte, error) {
	return []byte{}, nil
}

func (pl *lStatePool) Get(helper *structure.Helper, values *structure.RequestData) map[string]*string {
	return nil
}

func (pl *lStatePool) Put(L map[string]*string) {
}

func (pl *lStatePool) Shutdown() {
}
