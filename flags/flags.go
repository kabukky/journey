package flags

import (
	"log"

	flag "github.com/spf13/pflag"
)

var (
	Log         = ""
	CustomPath  = ""
	IsInDevMode = false
	HttpPort    = ""
	HttpsPort   = ""
)

func init() {
	// Parse all flags
	parseFlags()
	if IsInDevMode {
		// TODO: move message closer to place where it matters
		log.Println("Starting Journey in developer mode...")
	}
}

func parseFlags() {
	// Check if the log should be output to a file
	flag.StringVar(&Log, "log", "", "Use this option to save to log output to a file. Note: Journey needs create, read, and write access to that file. Example: --log=path/to/log.txt")
	// Check if a custom content path has been provided by the user
	flag.StringVar(&CustomPath, "custom-path", "", "Specify a custom path to store content files. Note: Journey needs read and write access to that path. A theme folder needs to be located in the custom path under content/themes. Example: --custom-path=/absolute/path/to/custom/folder")
	// Check if the development mode flag was provided by the user
	flag.BoolVar(&IsInDevMode, "dev", false, "Use this flag flag to put Journey in developer mode. Features of developer mode: Themes and plugins will be recompiled immediately after changes to the files. Example: --dev")
	// Check if the http port that was set in the config was overridden by the user
	flag.StringVar(&HttpPort, "http-port", "", "Use this option to override the HTTP port that was set in the config.json. Example: --http-port=8080")
	// Check if the http port that was set in the config was overridden by the user
	flag.StringVar(&HttpsPort, "https-port", "", "Use this option to override the HTTPS port that was set in the config.json. Example: --https-port=8081")
	flag.Parse()
}
