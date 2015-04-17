package templates

import (
	"bytes"
	"errors"
	"github.com/kabukky/journey/database"
	"github.com/kabukky/journey/filenames"
	"github.com/kabukky/journey/structure"
	"github.com/kabukky/journey/structure/methods"
	"net/http"
	"path/filepath"
	"sync"
)

type Templates struct {
	sync.RWMutex
	m map[string]*Helper
}

func newTemplates() *Templates { return &Templates{m: make(map[string]*Helper)} }

// Global compiled templates - thread safe and accessible from all packages
var compiledTemplates = newTemplates()

func ShowPostTemplate(writer http.ResponseWriter, slug string) error {
	compiledTemplates.RLock()
	defer compiledTemplates.RUnlock()
	blog, err := methods.GenerateBlog()
	if err != nil {
		return err
	}
	post, err := database.RetrievePostBySlug(slug)
	if err != nil {
		return err
	} else if !post.IsPublished { // Make sure the post is published before rendering it
		return errors.New("Post not published.")
	}
	requestData := structure.RequestData{Posts: make([]structure.Post, 1), Blog: blog, CurrentTemplate: 1} // CurrentTemplate = post
	requestData.Posts[0] = *post
	// If the post is a page and the page template is available, use the page template
	if post.IsPage {
		if template, ok := compiledTemplates.m["page"]; ok {
			_, err = writer.Write(executeHelper(template, &requestData, 1)) // context = post
			return err
		}
	}
	_, err = writer.Write(executeHelper(compiledTemplates.m["post"], &requestData, 1)) // context = post
	return err
}

func ShowAuthorTemplate(writer http.ResponseWriter, slug string, page int) error {
	compiledTemplates.RLock()
	defer compiledTemplates.RUnlock()
	postIndex := int64(page - 1)
	if postIndex < 0 {
		postIndex = 0
	}
	blog, err := methods.GenerateBlog()
	if err != nil {
		return err
	}
	author, err := database.RetrieveUserBySlug(slug)
	if err != nil {
		return err
	}
	posts, err := database.RetrievePostsByUser(author.Id, blog.PostsPerPage, (blog.PostsPerPage * postIndex))
	if err != nil {
		return err
	}
	requestData := structure.RequestData{Posts: posts, Blog: blog, CurrentIndexPage: page, CurrentTemplate: 3} // CurrentTemplate = author
	if template, ok := compiledTemplates.m["author"]; ok {
		_, err = writer.Write(executeHelper(template, &requestData, 0)) // context = index
	} else {
		_, err = writer.Write(executeHelper(compiledTemplates.m["index"], &requestData, 0)) // context = index
	}
	return err
}

func ShowTagTemplate(writer http.ResponseWriter, slug string, page int) error {
	compiledTemplates.RLock()
	defer compiledTemplates.RUnlock()
	postIndex := int64(page - 1)
	if postIndex < 0 {
		postIndex = 0
	}
	blog, err := methods.GenerateBlog()
	if err != nil {
		return err
	}
	tag, err := database.RetrieveTagBySlug(slug)
	if err != nil {
		return err
	}
	posts, err := database.RetrievePostsByTag(tag.Id, blog.PostsPerPage, (blog.PostsPerPage * postIndex))
	if err != nil {
		return err
	}
	requestData := structure.RequestData{Posts: posts, Blog: blog, CurrentIndexPage: page, CurrentTag: tag, CurrentTemplate: 2} // CurrentTemplate = tag
	if template, ok := compiledTemplates.m["tag"]; ok {
		_, err = writer.Write(executeHelper(template, &requestData, 0)) // context = index
	} else {
		_, err = writer.Write(executeHelper(compiledTemplates.m["index"], &requestData, 0)) // context = index
	}
	return err
}

func ShowIndexTemplate(writer http.ResponseWriter, page int) error {
	compiledTemplates.RLock()
	defer compiledTemplates.RUnlock()
	postIndex := int64(page - 1)
	if postIndex < 0 {
		postIndex = 0
	}
	blog, err := methods.GenerateBlog()
	if err != nil {
		return err
	}
	posts, err := database.RetrievePostsForIndex(blog.PostsPerPage, (blog.PostsPerPage * postIndex))
	if err != nil {
		return err
	}
	requestData := structure.RequestData{Posts: posts, Blog: blog, CurrentIndexPage: page, CurrentTemplate: 0} // CurrentTemplate = index
	_, err = writer.Write(executeHelper(compiledTemplates.m["index"], &requestData, 0))                        // context = index
	return err
}

func GetAllThemes() []string {
	themes := make([]string, 0)
	files, _ := filepath.Glob(filepath.Join(filenames.ThemesFilepath, "*"))
	for _, file := range files {
		if isDirectory(file) {
			themes = append(themes, filepath.Base(file))
		}
	}
	return themes
}

func executeHelper(helper *Helper, values *structure.RequestData, context int) []byte {
	// Set context and set it back to the old value once fuction returns
	defer setCurrentHelperContext(values, values.CurrentHelperContext)
	values.CurrentHelperContext = context

	block := helper.Block
	indexTracker := 0
	extended := false
	var extendHelper *Helper
	for index, child := range helper.Children {
		// Handle extend helper
		if index == 0 && child.Name == "!<" {
			extended = true
			extendHelper = compiledTemplates.m[string(child.Function(&child, values))]
		} else {
			var buffer bytes.Buffer
			toAdd := child.Function(&child, values)
			buffer.Write(block[:child.Position+indexTracker])
			buffer.Write(toAdd)
			buffer.Write(block[child.Position+indexTracker:])
			block = buffer.Bytes()
			indexTracker += len(toAdd)
		}
	}
	if extended {
		extendHelper.BodyHelper.Block = block
		return executeHelper(extendHelper, values, values.CurrentHelperContext) // TODO: not sure if context = values.CurrentHelperContext is right.
	}
	return block
}

func setCurrentHelperContext(values *structure.RequestData, context int) {
	values.CurrentHelperContext = context
}
