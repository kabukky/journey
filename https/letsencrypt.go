package https

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/url"

	"github.com/rkuris/journey/configuration"
	"github.com/rkuris/journey/filenames"
	"golang.org/x/crypto/acme/autocert"
)

func buildLetsEncryptServer(addr string, handler http.Handler) *http.Server {
	// Get host from HTTPS URL
	httpsURL, err := url.Parse(configuration.Config.HTTPSUrl)
	if err != nil {
		log.Fatal("Fatal error: Couldn't parse HTTPSUrl field in config.")
	}
	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(httpsURL.Host),
		Cache:      autocert.DirCache(filenames.HTTPSFilepath),
	}
	server := &http.Server{
		Addr:    addr,
		Handler: handler,
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
	}
	return server
}
