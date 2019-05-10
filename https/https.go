package https

import (
	"net/http"

	"github.com/kabukky/journey/configuration"
	"github.com/kabukky/journey/filenames"
)

func StartServer(addr string, handler http.Handler) error {
	if configuration.Config.UseLetsEncrypt {
		server := buildLetsEncryptServer(addr, handler)
		return server.ListenAndServeTLS("", "")
	}

	checkCertificates()
	return http.ListenAndServeTLS(addr, filenames.HttpsCertFilename, filenames.HttpsKeyFilename, handler)
}
