package filenames

import (
	"github.com/kabukky/journey/flags"
	"github.com/kardianos/osext"
	"log"
	"os"
	"path/filepath"
)

var (
	// Initialization of the working directory - needed to load relative assets
	_ = initializeWorkingDirectory()

	// For assets that are created, changed, our user-provided while running journey
	ConfigFilename   = filepath.Join(flags.CustomPath, "config.json")
	LogFilename      = filepath.Join(flags.CustomPath, "log.txt")
	DatabaseFilename = filepath.Join(flags.CustomPath, "content", "data", "journey.db")
	ThemesFilepath   = filepath.Join(flags.CustomPath, "content", "themes")
	ImagesFilepath   = filepath.Join(flags.CustomPath, "content", "images")
	ContentFilepath  = filepath.Join(flags.CustomPath, "content")
	PluginsFilepath  = filepath.Join(flags.CustomPath, "content", "plugins")

	// For https
	HttpsCertFilename = filepath.Join(flags.CustomPath, "content", "https", "cert.pem")
	HttpsKeyFilename  = filepath.Join(flags.CustomPath, "content", "https", "key.pem")

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
)

func init() {
	// Create content directories if they are not created already
	err := createDirectories()
	if err != nil {
		log.Fatalln("Error: Couldn't create directories: " + err.Error())
	}

}

func createDirectories() error {
	paths := []string{filepath.Join(flags.CustomPath, "content", "data"), filepath.Join(flags.CustomPath, "content", "themes"), filepath.Join(flags.CustomPath, "content", "images"), filepath.Join(flags.CustomPath, "content", "https")}
	for _, path := range paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			log.Println("Creating " + path)
			err := os.MkdirAll(path, 0776)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func initializeWorkingDirectory() error {
	// Set working directory to the path this executable is in
	executablePath, err := osext.ExecutableFolder()
	if err != nil {
		log.Fatal("Error: Couldn't determine working directory: " + err.Error())
		return err
	}
	os.Chdir(executablePath)
	return nil
}
