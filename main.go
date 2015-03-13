package main

import (
	"github.com/kabukky/journey/certificates"
	"github.com/kabukky/journey/configuration"
	"github.com/kabukky/journey/database"
	"github.com/kabukky/journey/filenames"
	"github.com/kabukky/journey/server"
	"github.com/kabukky/journey/templates"
	"log"
	"net/http"
	"os"
	"runtime"
)

func httpsRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://"+configuration.Config.HttpsUrl+r.RequestURI, http.StatusMovedPermanently)
	return
}

func checkHttpsCertificates() {
	// Check https certificates. If they are not available generate temporary ones for testing.
	err := certificates.Check(filenames.HttpsCertFilename, filenames.HttpsKeyFilename)
	if err != nil {
		log.Println("Warning: couldn't load https certs. Generating new ones. Replace " + filenames.HttpsCertFilename + " and " + filenames.HttpsKeyFilename + " with your own certificates as soon as possible!")
		err := certificates.Generate(filenames.HttpsCertFilename, filenames.HttpsKeyFilename, configuration.Config.HttpsUrl)
		if err != nil {
			log.Fatal("Error: Couldn't create https certificates.")
			return
		}
	}
}

func main() {
	// Setup
	runtime.GOMAXPROCS(runtime.NumCPU()) // Maybe not needed
	// Write log to file
	f, err := os.OpenFile(filenames.LogFilename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Error: Counldn't open log file: " + err.Error())
	}
	defer f.Close()
	log.SetOutput(f)

	// Configuration is read from config.json by loading the configuration package

	// Database
	err = database.Initialize()
	if err != nil {
		log.Fatal("Error: Couldn't initialize database: " + err.Error())
		return
	}

	// Templates
	err = templates.Generate()
	if err != nil {
		log.Fatal("Error: Couldn't compile templates: " + err.Error())
		return
	}

	// HTTP(S) Server
	// Determine the kind of https support (as set in the config.json)
	switch configuration.Config.HttpsUsage {
	case "AdminOnly":
		checkHttpsCertificates()
		httpMux := http.NewServeMux()
		httpsMux := http.NewServeMux()
		// Blog as http
		server.InitializeBlog(httpMux)
		// Blog as https
		server.InitializeBlog(httpsMux)
		// Admin as https and http redirect
		// Add redirection to http mux
		httpMux.Handle("/admin/", http.HandlerFunc(httpsRedirect))
		// Add routes to https mux
		server.InitializeAdmin(httpsMux)
		// Start https server
		log.Println("Starting https server on port " + configuration.Config.HttpsHostAndPort + "...")
		go http.ListenAndServeTLS(configuration.Config.HttpsHostAndPort, filenames.HttpsCertFilename, filenames.HttpsKeyFilename, httpsMux)
		// Start http server
		log.Println("Starting http server on port " + configuration.Config.HttpHostAndPort + "...")
		http.ListenAndServe(configuration.Config.HttpHostAndPort, httpMux)
	case "All":
		checkHttpsCertificates()
		httpsMux := http.NewServeMux()
		// Blog as https
		server.InitializeBlog(httpsMux)
		// Admin as https
		server.InitializeAdmin(httpsMux)
		// Start https server
		log.Println("Starting https server on port " + configuration.Config.HttpsHostAndPort + "...")
		go http.ListenAndServeTLS(configuration.Config.HttpsHostAndPort, filenames.HttpsCertFilename, filenames.HttpsKeyFilename, httpsMux)
		// Start http server
		log.Println("Starting http server on port " + configuration.Config.HttpHostAndPort + "...")
		http.ListenAndServe(configuration.Config.HttpHostAndPort, http.HandlerFunc(httpsRedirect))
	default: // This is configuration.HttpsUsage == "None"
		httpMux := http.NewServeMux()
		// Blog as http
		server.InitializeBlog(httpMux)
		// Admin as http
		server.InitializeAdmin(httpMux)
		// Start http server
		log.Println("Starting server without HTTPS support. Please enable HTTPS in " + filenames.ConfigFilename + " to improve security.")
		log.Println("Starting http server on port " + configuration.Config.HttpHostAndPort + "...")
		http.ListenAndServe(configuration.Config.HttpHostAndPort, httpMux)
	}
}
