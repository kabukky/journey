package filenames

import (
	"flag"
	"github.com/kardianos/osext"
	"log"
	"os"
	"path/filepath"
)

var (
	customPath = ""
	// Initialization of the working directory - needed to load relative assets
	_ = initializeWorkingDirectory()

	// For assets that are created or changed while running journey
	ConfigFilename   = filepath.Join(customPath, "config.json")
	LogFilename      = filepath.Join(customPath, "log.txt")
	DatabaseFilename = filepath.Join(customPath, "content", "data", "journey.db")
	ThemesFilepath   = filepath.Join(customPath, "content", "themes")
	ImagesFilepath   = filepath.Join(customPath, "content", "images")
	ContentFilepath  = filepath.Join(customPath, "content")

	// For https
	HttpsCertFilename = filepath.Join(customPath, "content", "https", "cert.pem")
	HttpsKeyFilename  = filepath.Join(customPath, "content", "https", "key.pem")

	//For built-in files (e.g. the admin interface)
	AdminFilepath  = filepath.Join("built-in", "admin")
	PublicFilepath = filepath.Join("built-in", "public")

	// For handlebars (this is a url string)
	JqueryFilename = "/public/jquery/jquery.js"

	// For blog  (this is a url string)
	// TODO: This is not used at the moment because it is still hard-coded into the create database string
	DefaultBlogLogoFilename  = "/public/images/blog-logo.jpg"
	DefaultBlogCoverFilename = "/public/images/blog-cover.jpg"

	// For users (this is a url string)
	DefaultUserImageFilename = "/public/images/user-image.jpg"
	DefaultUserCoverFilename = "/public/images/user-cover.jpg"

	// Create content directories if they are not created already
	_ = createDirectories()
)

func createDirectories() error {
	paths := []string{filepath.Join(customPath, "content", "data"), filepath.Join(customPath, "content", "themes"), filepath.Join(customPath, "content", "images"), filepath.Join(customPath, "content", "https")}
	for _, path := range paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			log.Println("Creating " + path)
			err := os.MkdirAll(path, 0776)
			if err != nil {
				log.Fatal("Error: Couldn't create directory " + path + ": " + err.Error())
				return err
			}
		}
	}
	return nil
}

func initializeWorkingDirectory() error {
	// Check if a custom content path has been provided by the user
	flag.StringVar(&customPath, "custom-path", "", "Specify a custom path to store content files. Note: Journey needs read and write access to that path.")
	flag.Parse()

	// Set working directory to the path this executable is in
	executablePath, err := osext.ExecutableFolder()
	if err != nil {
		log.Fatal("Error: Couldn't determine working directory: " + err.Error())
		return err
	}
	os.Chdir(executablePath)
	return nil
}
