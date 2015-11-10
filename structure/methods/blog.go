package methods

import (
	"encoding/json"
	"github.com/kabukky/journey/configuration"
	"github.com/kabukky/journey/database"
	"github.com/kabukky/journey/slug"
	"github.com/kabukky/journey/structure"
	"log"
	"time"
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
	err = database.UpdateSettings(b.Title, b.Description, b.Logo, b.Cover, b.PostsPerPage, b.ActiveTheme, navigation, time.Now().UTC(), userId)
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
	err := database.UpdateActiveTheme(activeTheme, time.Now().UTC(), userId)
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
	blog.Url = []byte(configuration.Config.Url)
	blog.AssetPath = assetPath
	// Create navigation slugs
	for index, _ := range blog.NavigationItems {
		blog.NavigationItems[index].Slug = slug.Generate(blog.NavigationItems[index].Label, "navigation")
	}
	Blog = blog
	return nil
}
