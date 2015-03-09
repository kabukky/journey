package filenames

import (
	"path/filepath"
)

var ConfigFilename = "config.json"
var DatabaseFilename = filepath.Join("content", "data", "journey.db")
var ThemesFilepath = filepath.Join("content", "themes")
var ImagesFilepath = filepath.Join("content", "images")
var AdminFilepath = filepath.Join("built-in", "admin")
var PublicFilepath = filepath.Join("built-in", "public")
var ContentFilepath = "content"

// For https
var HttpsCertFilename = filepath.Join("content", "https", "cert.pem")
var HttpsKeyFilename = filepath.Join("content", "https", "key.pem")

// For handlebars
var JqueryFilename = "/public/jquery/jquery.js"

// For blog
// TODO: This is not used at the moment because it is still hard-coded into the create database string
var DefaultBlogLogoFilename = "/public/images/blog-logo.jpg"
var DefaultBlogCoverFilename = "/public/images/blog-cover.jpg"

// For users
var DefaultUserImageFilename = "/public/images/user-image.jpg"
var DefaultUserCoverFilename = "/public/images/user-cover.jpg"
