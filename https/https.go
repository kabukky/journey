package https

import (
	"net/http"

	"github.com/rkuris/journey/configuration"
	"github.com/rkuris/journey/filenames"
)

// StartServer ...
func StartServer(addr string, handler http.Handler) error {
	if configuration.Config.UseLetsEncrypt {
		server := buildLetsEncryptServer(addr, handler)
		return server.ListenAndServeTLS("", "")
	}
	checkCertificates()
	return http.ListenAndServeTLS(addr, filenames.HTTPSCertFilename, filenames.HTTPSKeyFilename, handler)
}
