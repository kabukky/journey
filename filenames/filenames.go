package filenames

import (
	"github.com/kardianos/osext"
	"log"
	"os"
	"path/filepath"
)

var (
	// Initialization of the working directory - needed to load relative assets
	_ = initializeWorkingDirectory()

	// For assets
	ConfigFilename   = "config.json"
	LogFilename      = "log.txt"
	DatabaseFilename = filepath.Join("content", "data", "journey.db")
	ThemesFilepath   = filepath.Join("content", "themes")
	ImagesFilepath   = filepath.Join("content", "images")
	AdminFilepath    = filepath.Join("built-in", "admin")
	PublicFilepath   = filepath.Join("built-in", "public")
	ContentFilepath  = "content"

	// For https
	HttpsCertFilename = filepath.Join("content", "https", "cert.pem")
	HttpsKeyFilename  = filepath.Join("content", "https", "key.pem")

	// For handlebars
	JqueryFilename = "/public/jquery/jquery.js"

	// For blog
	// TODO: This is not used at the moment because it is still hard-coded into the create database string
	DefaultBlogLogoFilename  = "/public/images/blog-logo.jpg"
	DefaultBlogCoverFilename = "/public/images/blog-cover.jpg"

	// For users
	DefaultUserImageFilename = "/public/images/user-image.jpg"
	DefaultUserCoverFilename = "/public/images/user-cover.jpg"
)

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
