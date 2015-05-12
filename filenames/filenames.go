package filenames

import (
	"github.com/kabukky/journey/flags"
	"github.com/kardianos/osext"
	"log"
	"os"
	"path/filepath"
)

var (
	// Determine the path the Journey executable is in - needed to load relative assets
	ExecutablePath = determineExecutablePath()

	// Determine the path the the content folder
	ContentPath = determineContentPath()

	// For assets that are created, changed, our user-provided while running journey
	ConfigFilename   = filepath.Join(ExecutablePath, "config.json")
	DatabaseFilepath = filepath.Join(ContentPath, "data")
	DatabaseFilename = filepath.Join(ContentPath, "data", "journey.db")
	ThemesFilepath   = filepath.Join(ContentPath, "themes")
	ImagesFilepath   = filepath.Join(ContentPath, "images")
	PluginsFilepath  = filepath.Join(ContentPath, "plugins")
	PagesFilepath    = filepath.Join(ContentPath, "pages")

	// For https
	HttpsFilepath     = filepath.Join(ContentPath, "https")
	HttpsCertFilename = filepath.Join(ContentPath, "https", "cert.pem")
	HttpsKeyFilename  = filepath.Join(ContentPath, "https", "key.pem")

	//For built-in files (e.g. the admin interface)
	AdminFilepath  = filepath.Join(ExecutablePath, "built-in", "admin")
	PublicFilepath = filepath.Join(ExecutablePath, "built-in", "public")
	HbsFilepath    = filepath.Join(ExecutablePath, "built-in", "hbs")

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
		log.Fatal("Error: Couldn't create directories:", err)
	}

}

func createDirectories() error {
	paths := []string{DatabaseFilepath, ThemesFilepath, ImagesFilepath, HttpsFilepath, PluginsFilepath, PagesFilepath}
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

func determineContentPath() string {
	contentPath := ""
	if flags.CustomPath != "" {
		contentPath = flags.CustomPath
	} else {
		contentPath = determineExecutablePath()
	}
	return filepath.Join(contentPath, "content")
}

func determineExecutablePath() string {
	// Get the path this executable is located in
	executablePath, err := osext.ExecutableFolder()
	if err != nil {
		log.Fatal("Error: Couldn't determine what directory this executable is in:", err)
	}
	return executablePath
}
