package templates

import (
	"bytes"
	"net/http"

	"github.com/kabukky/feeds"
	"github.com/rkuris/journey/database"
	"github.com/rkuris/journey/date"
	"github.com/rkuris/journey/structure"
	"github.com/rkuris/journey/structure/methods"
)

// ShowIndexRss shows the index rss
func ShowIndexRss(writer http.ResponseWriter) error {
	// Read lock global blog
	methods.Blog.RLock()
	defer methods.Blog.RUnlock()
	// 15 posts in rss for now
	posts, err := database.RetrievePostsForIndex(15, 0)
	if err != nil {
		return err
	}
	blogData := &structure.RequestData{Posts: posts, Blog: methods.Blog}
	feed := createFeed(blogData)
	err = feed.WriteRss(writer)
	return err
}

// ShowTagRss shows the tag rss
func ShowTagRss(writer http.ResponseWriter, slug string) error {
	// Read lock global blog
	methods.Blog.RLock()
	defer methods.Blog.RUnlock()
	tag, err := database.RetrieveTagBySlug(slug)
	if err != nil {
		return err
	}
	// 15 posts in rss for now
	posts, err := database.RetrievePostsByTag(tag.ID, 15, 0)
	if err != nil {
		return err
	}
	blogData := &structure.RequestData{Posts: posts, Blog: methods.Blog}
	feed := createFeed(blogData)
	err = feed.WriteRss(writer)
	return err
}

// ShowAuthorRss shows the author rss
func ShowAuthorRss(writer http.ResponseWriter, slug string) error {
	// Read lock global blog
	methods.Blog.RLock()
	defer methods.Blog.RUnlock()
	author, err := database.RetrieveUserBySlug(slug)
	if err != nil {
		return err
	}
	// 15 posts in rss for now
	posts, err := database.RetrievePostsByUser(author.ID, 15, 0)
	if err != nil {
		return err
	}
	blogData := &structure.RequestData{Posts: posts, Blog: methods.Blog}
	feed := createFeed(blogData)
	err = feed.WriteRss(writer)
	return err
}

func createFeed(values *structure.RequestData) *feeds.Feed {
	now := date.GetCurrentTime()
	feed := &feeds.Feed{
		Title:       string(values.Blog.Title),
		Description: string(values.Blog.Description),
		Link:        &feeds.Link{Href: string(values.Blog.URL)},
		Updated:     now,
		Image: &feeds.Image{
			Url:   string(values.Blog.URL) + string(values.Blog.Logo),
			Title: string(values.Blog.Title),
			Link:  string(values.Blog.URL),
		},
		Url: string(values.Blog.URL) + "/rss/",
	}
	for i := 0; i < len(values.Posts); i++ {
		if values.Posts[i].ID != 0 {
			// Make link
			var buffer bytes.Buffer
			buffer.Write(values.Blog.URL)
			buffer.WriteString("/")
			buffer.WriteString(values.Posts[i].Slug)
			item := &feeds.Item{
				Title:       string(values.Posts[i].Title),
				Description: string(values.Posts[i].HTML),
				Link:        &feeds.Link{Href: buffer.String()},
				Id:          string(values.Posts[i].UUID),
				Author:      &feeds.Author{Name: string(values.Posts[i].Author.Name), Email: ""},
				Created:     *values.Posts[i].Date,
			}
			// If the post has a cover image, add it to the item
			image := string(values.Posts[i].Image)
			if image != "" {
				item.Image = &feeds.Image{
					Url: string(values.Blog.URL) + image,
				}
			}
			feed.Items = append(feed.Items, item)
		}
	}

	return feed
}
