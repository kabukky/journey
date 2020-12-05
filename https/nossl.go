// +build nossl

package https

import (
	"net/http"
)

func StartServer(addr string, handler http.Handler) error {
	return nil
}
