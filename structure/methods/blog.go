package methods

import (
	"encoding/json"
	"log"

	"github.com/rkuris/journey/configuration"
	"github.com/rkuris/journey/database"
	"github.com/rkuris/journey/date"
	"github.com/rkuris/journey/slug"
	"github.com/rkuris/journey/structure"
)

// Global blog - thread safe and accessible by all requests
var Blog *structure.Blog

var assetPath = []byte("/assets/")

func UpdateBlog(b *structure.Blog, userId int64) error {
	// Marshal navigation items to json string
	navigation, err := json.Marshal(b.NavigationItems)
	if err != nil {
		return err
	}
	err = database.UpdateSettings(b.Title, b.Description, b.Logo, b.Cover, b.PostsPerPage, b.ActiveTheme, navigation, date.GetCurrentTime(), userId)
	if err != nil {
		return err
	}
	// Generate new global blog
	err = GenerateBlog()
	if err != nil {
		log.Panic("Error: couldn't generate blog data:", err)
	}
	return nil
}

func UpdateActiveTheme(activeTheme string, userId int64) error {
	err := database.UpdateActiveTheme(activeTheme, date.GetCurrentTime(), userId)
	if err != nil {
		return err
	}
	// Generate new global blog
	err = GenerateBlog()
	if err != nil {
		log.Panic("Error: couldn't generate blog data:", err)
	}
	return nil
}

func GenerateBlog() error {
	// Write lock the global blog
	if Blog != nil {
		Blog.Lock()
		defer Blog.Unlock()
	}
	// Generate blog from db
	blog, err := database.RetrieveBlog()
	if err != nil {
		return err
	}
	// Add parameters that are not saved in db
	blog.URL = []byte(configuration.Config.URL)
	blog.AssetPath = assetPath
	// Create navigation slugs
	for index := range blog.NavigationItems {
		blog.NavigationItems[index].Slug = slug.Generate(blog.NavigationItems[index].Label, "navigation")
	}
	Blog = blog
	return nil
}
