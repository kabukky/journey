package templates

import (
	"bytes"
	"github.com/gorilla/feeds"
	"github.com/kabukky/journey/database"
	"github.com/kabukky/journey/structure"
	"github.com/kabukky/journey/structure/methods"
	"net/http"
	"time"
)

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

func ShowTagRss(writer http.ResponseWriter, slug string) error {
	// Read lock global blog
	methods.Blog.RLock()
	defer methods.Blog.RUnlock()
	tag, err := database.RetrieveTagBySlug(slug)
	if err != nil {
		return err
	}
	// 15 posts in rss for now
	posts, err := database.RetrievePostsByTag(tag.Id, 15, 0)
	if err != nil {
		return err
	}
	blogData := &structure.RequestData{Posts: posts, Blog: methods.Blog}
	feed := createFeed(blogData)
	err = feed.WriteRss(writer)
	return err
}

func ShowAuthorRss(writer http.ResponseWriter, slug string) error {
	// Read lock global blog
	methods.Blog.RLock()
	defer methods.Blog.RUnlock()
	author, err := database.RetrieveUserBySlug(slug)
	if err != nil {
		return err
	}
	// 15 posts in rss for now
	posts, err := database.RetrievePostsByUser(author.Id, 15, 0)
	if err != nil {
		return err
	}
	blogData := &structure.RequestData{Posts: posts, Blog: methods.Blog}
	feed := createFeed(blogData)
	err = feed.WriteRss(writer)
	return err
}

func createFeed(values *structure.RequestData) *feeds.Feed {
	now := time.Now()
	feed := &feeds.Feed{
		Title:       string(values.Blog.Title),
		Description: string(values.Blog.Description),
		Link:        &feeds.Link{Href: string(values.Blog.Url)},
		Created:     now,
	}
	for i := 0; i < len(values.Posts); i++ {
		if values.Posts[i].Id != 0 {
			// Make link
			var buffer bytes.Buffer
			buffer.Write(values.Blog.Url)
			buffer.WriteString("/")
			buffer.WriteString(values.Posts[i].Slug)
			feed.Items = append(feed.Items, &feeds.Item{
				Title:       string(values.Posts[i].Title),
				Description: string(values.Posts[i].Html),
				Link:        &feeds.Link{Href: buffer.String()},
				Id:          string(values.Posts[i].Uuid),
				Author:      &feeds.Author{Name: string(values.Posts[i].Author.Name), Email: ""},
				Created:     *values.Posts[i].Date,
			})
		}
	}

	return feed
}
