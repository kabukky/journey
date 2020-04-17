package https

import (
	"log"

	"github.com/kabukky/httpscerts"
	"github.com/kabukky/journey/configuration"
	"github.com/kabukky/journey/filenames"
)

func checkCertificates() {
	// Check https certificates. If they are not available generate temporary ones for testing.
	if err := httpscerts.Check(filenames.HTTPSCertFilename, filenames.HTTPSKeyFilename); err != nil {
		log.Println("Warning: couldn't load https certs. Generating new ones. Replace " + filenames.HTTPSCertFilename + " and " + filenames.HTTPSKeyFilename + " with your own certificates as soon as possible!")
		if err := httpscerts.Generate(filenames.HTTPSCertFilename, filenames.HTTPSKeyFilename, configuration.Config.HTTPSUrl); err != nil {
			log.Fatal("Error: Couldn't create https certificates.")
			return
		}
	}
}
