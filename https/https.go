package https

import (
	"net/http"

	"github.com/kabukky/journey/filenames"
	"github.com/kabukky/journey/flags"
)

func StartServer(addr string, handler http.Handler) error {
	if flags.UseLetsEncrypt {
		server := buildLetsEncryptServer(addr, handler)
		return server.ListenAndServeTLS("", "")
	} else {
		checkCertificates()
		return http.ListenAndServeTLS(addr, filenames.HttpsCertFilename, filenames.HttpsKeyFilename, handler)
	}
}
