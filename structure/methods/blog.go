package methods

import (
	"github.com/kabukky/journey/configuration"
	"github.com/kabukky/journey/database"
	"github.com/kabukky/journey/structure"
)

var assetPath = []byte("/assets/")

func UpdateBlog(b *structure.Blog) error {
	err := database.UpdateSettings(b.Title, b.Description, b.Logo, b.Cover, b.PostsPerPage, b.ActiveTheme)
	if err != nil {
		return err
	}
	return nil
}

func GenerateBlog() (*structure.Blog, error) {
	// Generate blog from db
	blog, err := database.RetrieveBlog()
	if err != nil {
		return nil, err
	}
	// Add parameters that are not saved in db
	blog.Url = []byte(configuration.Config.HttpUrl)
	blog.AssetPath = assetPath
	return blog, nil
}
