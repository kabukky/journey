package filenames

import (
	"github.com/kabukky/journey/flags"
	"github.com/kardianos/osext"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	// Determine the path the Journey executable is in - needed to load relative assets
	ExecutablePath = determineExecutablePath()

	// Determine the path to the assets folder (default: Journey root folder)
	AssetPath = determineAssetPath()

	// For assets that are created, changed, our user-provided while running journey
	ConfigFilename   = filepath.Join(AssetPath, "config.json")
	ContentFilepath  = filepath.Join(AssetPath, "content")
	DatabaseFilepath = filepath.Join(ContentFilepath, "data")
	DatabaseFilename = filepath.Join(ContentFilepath, "data", "journey.db")
	ThemesFilepath   = filepath.Join(ContentFilepath, "themes")
	ImagesFilepath   = filepath.Join(ContentFilepath, "images")
	PluginsFilepath  = filepath.Join(ContentFilepath, "plugins")
	PagesFilepath    = filepath.Join(ContentFilepath, "pages")

	// For https
	HttpsFilepath     = filepath.Join(ContentFilepath, "https")
	HttpsCertFilename = filepath.Join(ContentFilepath, "https", "cert.pem")
	HttpsKeyFilename  = filepath.Join(ContentFilepath, "https", "key.pem")

	//For built-in files (e.g. the admin interface)
	AdminFilepath  = filepath.Join(ExecutablePath, "built-in", "admin")
	PublicFilepath = filepath.Join(ExecutablePath, "built-in", "public")
	HbsFilepath    = filepath.Join(ExecutablePath, "built-in", "hbs")

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

func bashPath(contentPath string) string {
	if strings.HasPrefix(contentPath, "~") {
		return strings.Replace(contentPath, "~", os.Getenv("HOME"), 1)
	}
	return contentPath
}

func determineAssetPath() string {
	if flags.CustomPath != "" {
		contentPath, err := filepath.Abs(bashPath(flags.CustomPath))
		if err != nil {
			log.Fatal("Error: Couldn't read from custom path:", err)
		}
		return contentPath
	}
	return determineExecutablePath()
}

func determineExecutablePath() string {
	// Get the path this executable is located in
	executablePath, err := osext.ExecutableFolder()
	if err != nil {
		log.Fatal("Error: Couldn't determine what directory this executable is in:", err)
	}
	return executablePath
}
