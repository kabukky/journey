package main

import (
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"encoding/json"

	"github.com/dimfeld/httptreemux"
	"github.com/thoas/stats"
	"github.com/kabukky/journey/authentication"
	"github.com/kabukky/journey/configuration"
	"github.com/kabukky/journey/database"
	"github.com/kabukky/journey/filenames"
	"github.com/kabukky/journey/flags"
	"github.com/kabukky/journey/https"
	"github.com/kabukky/journey/plugins"
	"github.com/kabukky/journey/server"
	"github.com/kabukky/journey/structure/methods"
	"github.com/kabukky/journey/templates"
)

func httpsRedirect(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	http.Redirect(w, r, configuration.Config.HttpsUrl+r.RequestURI, http.StatusMovedPermanently)
	return
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
	if err = database.Initialize(); err != nil {
		log.Fatal("Error: Couldn't initialize database:", err)
		return
	}

	// Global blog data
	if err = methods.GenerateBlog(); err != nil {
		log.Fatal("Error: Couldn't generate blog data:", err)
		return
	}

	// Templates
	if err = templates.Generate(); err != nil {
		log.Fatal("Error: Couldn't compile templates:", err)
		return
	}

	// Plugins
	if err = plugins.Load(); err == nil {
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
			if err := https.StartServer(httpsPort, httpsRouter); err != nil {
				log.Fatal("Error: Couldn't start the HTTPS server:", err)
			}
		}()
		// Start http server
		log.Println("Starting http server on port " + httpPort + "...")
		if err := http.ListenAndServe(httpPort, httpRouter); err != nil {
			log.Fatal("Error: Couldn't start the HTTP server:", err)
		}
	case "All":
		httpsStats := stats.New()
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
		httpsRouter.GET("/admin/stats", func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
			userName := authentication.GetUserName(r)
			if userName != "" {
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				s, err := json.Marshal(httpsStats.Data())
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				w.Write(s)
			} else {
				http.Error(w, "Not logged in!", http.StatusInternalServerError)
				return
			}
		})
		go func() {
			if err := https.StartServer(httpsPort, httpsStats.Handler(httpsRouter)); err != nil {
				log.Fatal("Error: Couldn't start the HTTPS server:", err)
			}
		}()
		// Start http server
		log.Println("Starting http server on port " + httpPort + "...")
		if err := http.ListenAndServe(httpPort, httpRouter); err != nil {
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
		if err := http.ListenAndServe(httpPort, httpRouter); err != nil {
			log.Fatal("Error: Couldn't start the HTTP server:", err)
		}
	}
}
