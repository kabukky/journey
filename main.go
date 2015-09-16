package main

import (
	"github.com/dimfeld/httptreemux"
	"github.com/kabukky/httpscerts"
	"github.com/kabukky/journey/configuration"
	"github.com/kabukky/journey/database"
	"github.com/kabukky/journey/filenames"
	"github.com/kabukky/journey/flags"
	"github.com/kabukky/journey/plugins"
	"github.com/rmrio/journey/server"
	"github.com/kabukky/journey/structure/methods"
	"github.com/kabukky/journey/templates"
	"github.com/justinas/alice"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
)

func httpsRedirect(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	http.Redirect(w, r, configuration.Config.HttpsUrl+r.RequestURI, http.StatusMovedPermanently)
	return
}

func checkHttpsCertificates() {
	// Check https certificates. If they are not available generate temporary ones for testing.
	err := httpscerts.Check(filenames.HttpsCertFilename, filenames.HttpsKeyFilename)
	if err != nil {
		log.Println("Warning: couldn't load https certs. Generating new ones. Replace " + filenames.HttpsCertFilename + " and " + filenames.HttpsKeyFilename + " with your own certificates as soon as possible!")
		err := httpscerts.Generate(filenames.HttpsCertFilename, filenames.HttpsKeyFilename, configuration.Config.HttpsUrl)
		if err != nil {
			log.Fatal("Error: Couldn't create https certificates.")
			return
		}
	}
}

func main() {
	// Setup
	var err error

	// GOMAXPROCS - Maybe not needed
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Write log to file if the log flag was provided
	if flags.Log != "" {
		logFile, err := os.OpenFile(flags.Log, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal("Error: Couldn't open log file: " + err.Error())
		}
		defer logFile.Close()
		log.SetOutput(logFile)
	}

	// Configuration is read from config.json by loading the configuration package

	// Database
	err = database.Initialize()
	if err != nil {
		log.Fatal("Error: Couldn't initialize database:", err)
		return
	}

	// Global blog data
	err = methods.GenerateBlog()
	if err != nil {
		log.Fatal("Error: Couldn't generate blog data:", err)
		return
	}

	// Templates
	err = templates.Generate()
	if err != nil {
		log.Fatal("Error: Couldn't compile templates:", err)
		return
	}

	// Plugins
	err = plugins.Load()
	if err == nil {
		// Close LuaPool at the end
		defer plugins.LuaPool.Shutdown()
		log.Println("Plugins loaded.")
	}

	// HTTP(S) Server
	httpPort := configuration.Config.HttpHostAndPort
	httpsPort := configuration.Config.HttpsHostAndPort
	// Check if HTTP/HTTPS flags were provided
	if flags.HttpPort != "" {
		components := strings.SplitAfterN(httpPort, ":", 2)
		httpPort = components[0] + flags.HttpPort
	}
	if flags.HttpsPort != "" {
		components := strings.SplitAfterN(httpsPort, ":", 2)
		httpsPort = components[0] + flags.HttpsPort
	}
	// Determine the kind of https support (as set in the config.json)
	switch configuration.Config.HttpsUsage {
	case "AdminOnly":
		checkHttpsCertificates()
		httpRouter := httptreemux.New()
		httpsRouter := httptreemux.New()
		// Blog and pages as http
		server.InitializeBlog(httpRouter)
		server.InitializePages(httpRouter)
		// Blog and pages as https
		server.InitializeBlog(httpsRouter)
		server.InitializePages(httpsRouter)
		// Admin as https and http redirect
		// Add redirection to http router
		httpRouter.GET("/admin/", httpsRedirect)
		httpRouter.GET("/admin/*path", httpsRedirect)
		// Add routes to https router
		server.InitializeAdmin(httpsRouter)
		// Start https server
		log.Println("Starting https server on port " + httpsPort + "...")
		go func() {
			chain := alice.New(server.CheckHost).Then(httpsRouter)
			err := http.ListenAndServeTLS(httpsPort, filenames.HttpsCertFilename, filenames.HttpsKeyFilename, chain)
			if err != nil {
				log.Fatal("Error: Couldn't start the HTTPS server:", err)
			}
		}()
		// Start http server
		log.Println("Starting http server on port " + httpPort + "...")
		chain := alice.New(server.CheckHost).Then(httpRouter)
		err := http.ListenAndServe(httpPort, chain)
		if err != nil {
			log.Fatal("Error: Couldn't start the HTTP server:", err)
		}
	case "All":
		checkHttpsCertificates()
		httpsRouter := httptreemux.New()
		httpRouter := httptreemux.New()
		// Blog and pages as https
		server.InitializeBlog(httpsRouter)
		server.InitializePages(httpsRouter)
		// Admin as https
		server.InitializeAdmin(httpsRouter)
		// Add redirection to http router
		httpRouter.GET("/", httpsRedirect)
		httpRouter.GET("/*path", httpsRedirect)
		// Start https server
		log.Println("Starting https server on port " + httpsPort + "...")
		go func() {
			chain := alice.New(server.CheckHost).Then(httpsRouter)
			err := http.ListenAndServeTLS(httpsPort, filenames.HttpsCertFilename, filenames.HttpsKeyFilename, chain)
			if err != nil {
				log.Fatal("Error: Couldn't start the HTTPS server:", err)
			}
		}()
		// Start http server
		log.Println("Starting http server on port " + httpPort + "...")
		chain := alice.New(server.CheckHost).Then(httpRouter)
		err := http.ListenAndServe(httpPort, chain)
		if err != nil {
			log.Fatal("Error: Couldn't start the HTTP server:", err)
		}
	default: // This is configuration.HttpsUsage == "None"
		httpRouter := httptreemux.New()
		// Blog and pages as http
		server.InitializeBlog(httpRouter)
		server.InitializePages(httpRouter)
		// Admin as http
		server.InitializeAdmin(httpRouter)
		// Start http server
		log.Println("Starting server without HTTPS support. Please enable HTTPS in " + filenames.ConfigFilename + " to improve security.")
		log.Println("Starting http server on port " + httpPort + "...")
		chain := alice.New(server.CheckHost).Then(httpRouter)
		err := http.ListenAndServe(httpPort, chain)
		if err != nil {
			log.Fatal("Error: Couldn't start the HTTP server:", err)
		}
	}
}
