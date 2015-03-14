package templates

import (
	"bytes"
	"github.com/gorilla/feeds"
	"github.com/kabukky/journey/database"
	"github.com/kabukky/journey/structure"
	"net/http"
	"time"
)

func ShowIndexRss(writer http.ResponseWriter) error {
	// 15 posts in rss for now
	posts, err := database.RetrievePostsForIndex(15, 0)
	if err != nil {
		return err
	}
	blog, err := database.RetrieveBlog()
	if err != nil {
		return err
	}
	blogData := &structure.RequestData{Posts: posts, Blog: blog}
	feed := createFeed(blogData)
	err = feed.WriteRss(writer)
	return err
}

func ShowTagRss(writer http.ResponseWriter, slug string) error {
	tag, err := database.RetrieveTagBySlug(slug)
	if err != nil {
		return err
	}
	// 15 posts in rss for now
	posts, err := database.RetrievePostsByTag(tag.Id, 15, 0)
	if err != nil {
		return err
	}
	blog, err := database.RetrieveBlog()
	if err != nil {
		return err
	}
	blogData := &structure.RequestData{Posts: posts, Blog: blog}
	feed := createFeed(blogData)
	err = feed.WriteRss(writer)
	return err
}

func ShowAuthorRss(writer http.ResponseWriter, slug string) error {
	author, err := database.RetrieveUserBySlug(slug)
	if err != nil {
		return err
	}
	// 15 posts in rss for now
	posts, err := database.RetrievePostsByUser(author.Id, 15, 0)
	if err != nil {
		return err
	}
	blog, err := database.RetrieveBlog()
	if err != nil {
		return err
	}
	blogData := &structure.RequestData{Posts: posts, Blog: blog}
	feed := createFeed(blogData)
	err = feed.WriteRss(writer)
	return err
}

func createFeed(values *structure.RequestData) *feeds.Feed {
	now := time.Now()
	feed := &feeds.Feed{
		Title:       string(makeCdata(values.Blog.Title)),
		Description: string(makeCdata(values.Blog.Description)),
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
				Title:       string(makeCdata(values.Posts[i].Title)),
				Description: string(makeCdata(values.Posts[i].Html)),
				Link:        &feeds.Link{Href: buffer.String()},
				Id:          string(values.Posts[i].Uuid),
				Author:      &feeds.Author{Name: string(values.Posts[i].Author.Name), Email: ""},
				Created:     values.Posts[i].Date,
			})
		}
	}

	return feed
}

func makeCdata(input []byte) []byte {
	var buffer bytes.Buffer
	buffer.WriteString("<![CDATA[")
	buffer.Write(input)
	buffer.WriteString("]]>")
	return buffer.Bytes()
}
