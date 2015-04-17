package flags

import (
	"flag"
	"log"
)

var (
	CustomPath  = ""
	IsInDevMode = false
)

func init() {
	// Parse all flags
	parseFlags()
	if IsInDevMode {
		log.Println("Starting Journey in developer mode...")
	}
}

func parseFlags() {
	// Check if a custom content path has been provided by the user
	flag.StringVar(&CustomPath, "custom-path", "", "Specify a custom path to store content files. Note: Journey needs read and write access to that path. Example: -custom-path=/absolute/path/to/custom/folder")
	flag.BoolVar(&IsInDevMode, "dev", false, "Set this flag to put Journey in developer mode. Features of developer mode: Themes will be recompiled immediately after changes to the files.")
	flag.Parse()
}
