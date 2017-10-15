package templates

import (
	"bytes"
	"errors"
	"github.com/kabukky/journey/database"
	"github.com/kabukky/journey/filenames"
	"github.com/kabukky/journey/helpers"
	"github.com/kabukky/journey/plugins"
	"github.com/kabukky/journey/structure"
	"github.com/kabukky/journey/structure/methods"
	"net/http"
	"path/filepath"
	"sync"
)

type Templates struct {
	sync.RWMutex
	m map[string]*structure.Helper
}

func newTemplates() *Templates { return &Templates{m: make(map[string]*structure.Helper)} }

// Global compiled templates - thread safe and accessible by all requests
var compiledTemplates = newTemplates()

func ShowPostTemplate(writer http.ResponseWriter, r *http.Request, slug string) error {
	// Read lock templates and global blog
	compiledTemplates.RLock()
	defer compiledTemplates.RUnlock()
	methods.Blog.RLock()
	defer methods.Blog.RUnlock()
	post, err := database.RetrievePostBySlug(slug)
	if err != nil {
		return err
	} else if !post.IsPublished { // Make sure the post is published before rendering it
		return errors.New("Post not published.")
	} else if post.Slug != slug {
		http.Redirect(writer, r, "/"+post.Slug+"/", 301)
		return nil
	}
	requestData := structure.RequestData{Posts: make([]structure.Post, 1), Blog: methods.Blog, CurrentTemplate: 1, CurrentPath: r.URL.Path} // CurrentTemplate = post
	requestData.Posts[0] = *post
	// Check if there's a custom page template available for this slug
	if template, ok := compiledTemplates.m["page-"+slug]; ok {
		_, err = writer.Write(executeHelper(template, &requestData, 1)) // context = post
		return err
	}
	// If the post is a page and the page template is available, use the page template
	if post.IsPage {
		if template, ok := compiledTemplates.m["page"]; ok {
			_, err = writer.Write(executeHelper(template, &requestData, 1)) // context = post
			return err
		}
	}
	_, err = writer.Write(executeHelper(compiledTemplates.m["post"], &requestData, 1)) // context = post
	if requestData.PluginVMs != nil {
		// Put the lua state map back into the pool
		plugins.LuaPool.Put(requestData.PluginVMs)
	}
	return err
}

func ShowAuthorTemplate(writer http.ResponseWriter, r *http.Request, slug string, page int) error {
	// Read lock templates and global blog
	compiledTemplates.RLock()
	defer compiledTemplates.RUnlock()
	methods.Blog.RLock()
	defer methods.Blog.RUnlock()
	postIndex := int64(page - 1)
	if postIndex < 0 {
		postIndex = 0
	}
	author, err := database.RetrieveUserBySlug(slug)
	if err != nil {
		return err
	}
	posts, err := database.RetrievePostsByUser(author.Id, methods.Blog.PostsPerPage, (methods.Blog.PostsPerPage * postIndex))
	if err != nil {
		return err
	}
	requestData := structure.RequestData{Posts: posts, Blog: methods.Blog, CurrentIndexPage: page, CurrentTemplate: 3, CurrentPath: r.URL.Path} // CurrentTemplate = author
	if template, ok := compiledTemplates.m["author"]; ok {
		_, err = writer.Write(executeHelper(template, &requestData, 0)) // context = index
	} else {
		_, err = writer.Write(executeHelper(compiledTemplates.m["index"], &requestData, 0)) // context = index
	}
	if requestData.PluginVMs != nil {
		// Put the lua state map back into the pool
		plugins.LuaPool.Put(requestData.PluginVMs)
	}
	return err
}

func ShowTagTemplate(writer http.ResponseWriter, r *http.Request, slug string, page int) error {
	// Read lock templates and global blog
	compiledTemplates.RLock()
	defer compiledTemplates.RUnlock()
	methods.Blog.RLock()
	defer methods.Blog.RUnlock()
	postIndex := int64(page - 1)
	if postIndex < 0 {
		postIndex = 0
	}
	tag, err := database.RetrieveTagBySlug(slug)
	if err != nil {
		return err
	}
	posts, err := database.RetrievePostsByTag(tag.Id, methods.Blog.PostsPerPage, (methods.Blog.PostsPerPage * postIndex))
	if err != nil {
		return err
	}
	requestData := structure.RequestData{Posts: posts, Blog: methods.Blog, CurrentIndexPage: page, CurrentTag: tag, CurrentTemplate: 2, CurrentPath: r.URL.Path} // CurrentTemplate = tag
	if template, ok := compiledTemplates.m["tag"]; ok {
		_, err = writer.Write(executeHelper(template, &requestData, 0)) // context = index
	} else {
		_, err = writer.Write(executeHelper(compiledTemplates.m["index"], &requestData, 0)) // context = index
	}
	if requestData.PluginVMs != nil {
		// Put the lua state map back into the pool
		plugins.LuaPool.Put(requestData.PluginVMs)
	}
	return err
}

func ShowIndexTemplate(w http.ResponseWriter, r *http.Request, page int) error {
	// Read lock templates and global blog
	compiledTemplates.RLock()
	defer compiledTemplates.RUnlock()
	methods.Blog.RLock()
	defer methods.Blog.RUnlock()
	postIndex := int64(page - 1)
	if postIndex < 0 {
		postIndex = 0
	}
	posts, err := database.RetrievePostsForIndex(methods.Blog.PostsPerPage, (methods.Blog.PostsPerPage * postIndex))
	if err != nil {
		return err
	}
	requestData := structure.RequestData{Posts: posts, Blog: methods.Blog, CurrentIndexPage: page, CurrentTemplate: 0, CurrentPath: r.URL.Path} // CurrentTemplate = index
	_, err = w.Write(executeHelper(compiledTemplates.m["index"], &requestData, 0))                                                              // context = index
	if requestData.PluginVMs != nil {
		// Put the lua state map back into the pool
		plugins.LuaPool.Put(requestData.PluginVMs)
	}
	return err
}

func GetAllThemes() []string {
	themes := make([]string, 0)
	files, _ := filepath.Glob(filepath.Join(filenames.ThemesFilepath, "*"))
	for _, file := range files {
		if helpers.IsDirectory(file) {
			themes = append(themes, filepath.Base(file))
		}
	}
	return themes
}

func executeHelper(helper *structure.Helper, values *structure.RequestData, context int) []byte {
	// Set context and set it back to the old value once fuction returns
	defer setCurrentHelperContext(values, values.CurrentHelperContext)
	values.CurrentHelperContext = context

	block := helper.Block
	indexTracker := 0
	extended := false
	var extendHelper *structure.Helper
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
