package https

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/url"

	"github.com/Landria/journey/configuration"
	"github.com/Landria/journey/filenames"
	"golang.org/x/crypto/acme/autocert"
)

func buildLetsEncryptServer(addr string, handler http.Handler) *http.Server {
	// Get host from HTTPS URL
	httpsUrl, err := url.Parse(configuration.Config.HttpsUrl)
	if err != nil {
		log.Fatal("Fatal error: Couldn't parse HttpsUrl field in config.")
	}
	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(httpsUrl.Host),
		Cache:      autocert.DirCache(filenames.HttpsFilepath),
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
