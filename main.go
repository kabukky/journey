package main

import (
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/dimfeld/httptreemux"
	apachelog "github.com/lestrrat-go/apache-logformat"
	"github.com/rkuris/journey/configuration"
	"github.com/rkuris/journey/database"
	"github.com/rkuris/journey/flags"
	"github.com/rkuris/journey/https"
	"github.com/rkuris/journey/notifications"
	"github.com/rkuris/journey/plugins"
	"github.com/rkuris/journey/server"
	"github.com/rkuris/journey/structure/methods"
	"github.com/rkuris/journey/templates"
)

func httpsRedirect(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	http.Redirect(w, r, configuration.Config.HTTPSUrl+r.RequestURI, http.StatusMovedPermanently)
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

	// Notification system
	// expand the SMTP password from the environment
	configuration.Config.SMTP.SMTPPassword = os.ExpandEnv(configuration.Config.SMTP.SMTPPassword)
	if err = notifications.Initialize(&configuration.Config.SMTP); err != nil {
		log.Fatal("Could not initialize notification system:", err)
	}

	// HTTP(S) Server
	httpPort := configuration.Config.HTTPHostAndPort
	httpsPort := configuration.Config.HTTPSHostAndPort
	// Check if HTTP/HTTPS flags were provided
	if flags.HTTPPort != "" {
		components := strings.SplitAfterN(httpPort, ":", 2)
		httpPort = components[0] + flags.HTTPPort
	}
	if flags.HTTPSPort != "" {
		components := strings.SplitAfterN(httpsPort, ":", 2)
		httpsPort = components[0] + flags.HTTPSPort
	}
	// Determine the kind of https support (as set in the config.json)
	switch configuration.Config.HTTPSUsage {
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
			if err := https.StartServer(httpsPort, logWrapRouter(httpsRouter, configuration.Config)); err != nil {
				log.Fatal("Error: Couldn't start the HTTPS server:", err)
			}
		}()
		// Start http server
		log.Println("Starting http-redirect server on port " + httpPort + "...")
		if err := http.ListenAndServe(httpPort, httpRouter); err != nil {
			log.Fatal("Error: Couldn't start the HTTP server:", err)
		}
	case "All":
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
		log.Printf("Starting https server on port %q", httpsPort)
		go func() {
			if err := https.StartServer(httpsPort, logWrapRouter(httpsRouter, configuration.Config)); err != nil {
				log.Fatal("Couldn't start the HTTPS server", err)
			}
		}()
		// Start http server
		log.Printf("Starting http-redirect server on port %q", httpPort)
		if err := http.ListenAndServe(httpPort, httpRouter); err != nil {
			log.Fatal("Couldn't start the HTTP server:", err)
		}
	default: // This is configuration.HTTPSUsage == "None"
		httpRouter := httptreemux.New()
		// Blog and pages as http
		server.InitializeBlog(httpRouter)
		server.InitializePages(httpRouter)
		// Admin as http
		server.InitializeAdmin(httpRouter)
		// Start http server
		log.Printf("Starting http-only server on port %q", httpPort)
		if err := http.ListenAndServe(httpPort, logWrapRouter(httpRouter, configuration.Config)); err != nil {
			log.Fatal("Couldn't start the HTTP server:", err)
		}
	}
}
func logWrapRouter(httpRouter *httptreemux.TreeMux, config *configuration.Configuration) http.Handler {
	logfile := config.RequestLog
	format := config.RequestLogFormat
	if logfile != "" {
		var logFormat *apachelog.ApacheLog
		if format == "" {
			logFormat = apachelog.CombinedLog
		} else {
			var err error
			logFormat, err = apachelog.New(format)
			if err != nil {
				log.Fatalf("bad ApacheLog format: %s", err)
			}
		}
		fp, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error creating logfile: %s", err)
		}
		log.Printf("Logging to %q", logfile)
		return logFormat.Wrap(httpRouter, fp)
	}
	return httpRouter
}
