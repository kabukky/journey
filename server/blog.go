package server

import (
	"github.com/dimfeld/httptreemux"
	"github.com/kabukky/journey/filenames"
	"github.com/kabukky/journey/structure/methods"
	"github.com/kabukky/journey/templates"
	"net/http"
	"path/filepath"
	"strconv"
)

func indexHandler(w http.ResponseWriter, r *http.Request, params map[string]string) {
	number := params["number"]
	if number == "" {
		// Render index template (first page)
		err := templates.ShowIndexTemplate(w, r, 1)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	page, err := strconv.Atoi(number)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	// Render index template
	err = templates.ShowIndexTemplate(w, r, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}

func authorHandler(w http.ResponseWriter, r *http.Request, params map[string]string) {
	slug := params["slug"]
	function := params["function"]
	number := params["number"]
	if function == "" {
		// Render author template (first page)
		err := templates.ShowAuthorTemplate(w, r, slug, 1)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	} else if function == "rss" {
		// Render author rss feed
		err := templates.ShowAuthorRss(w, slug)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	page, err := strconv.Atoi(number)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	// Render author template
	err = templates.ShowAuthorTemplate(w, r, slug, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}

func tagHandler(w http.ResponseWriter, r *http.Request, params map[string]string) {
	slug := params["slug"]
	function := params["function"]
	number := params["number"]
	if function == "" {
		// Render tag template (first page)
		err := templates.ShowTagTemplate(w, r, slug, 1)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	} else if function == "rss" {
		// Render tag rss feed
		err := templates.ShowTagRss(w, slug)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	page, err := strconv.Atoi(number)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	// Render tag template
	err = templates.ShowTagTemplate(w, r, slug, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}

func postHandler(w http.ResponseWriter, r *http.Request, params map[string]string) {
	slug := params["slug"]
	if slug == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	} else if slug == "rss" {
		// Render index rss feed
		err := templates.ShowIndexRss(w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	// Render post template
	err := templates.ShowPostTemplate(w, r, slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}

func assetsHandler(w http.ResponseWriter, r *http.Request, params map[string]string) {
	// Read lock global blog
	methods.Blog.RLock()
	defer methods.Blog.RUnlock()
	http.ServeFile(w, r, filepath.Join(filenames.ThemesFilepath, methods.Blog.ActiveTheme, "assets", params["filepath"]))
	return
}

func imagesHandler(w http.ResponseWriter, r *http.Request, params map[string]string) {
	http.ServeFile(w, r, filepath.Join(filenames.ImagesFilepath, params["filepath"]))
	return
}

func publicHandler(w http.ResponseWriter, r *http.Request, params map[string]string) {
	http.ServeFile(w, r, filepath.Join(filenames.PublicFilepath, params["filepath"]))
	return
}

func InitializeBlog(router *httptreemux.TreeMux) {
	// For index
	router.GET("/", checkHost(indexHandler))
	router.GET("/:slug/", checkHost(postHandler))
	router.GET("/page/:number/", checkHost(indexHandler))
	// For author
	router.GET("/author/:slug/", checkHost(authorHandler))
	router.GET("/author/:slug/:function/", checkHost(authorHandler))
	router.GET("/author/:slug/:function/:number/", checkHost(authorHandler))
	// For tag
	router.GET("/tag/:slug/", checkHost(tagHandler))
	router.GET("/tag/:slug/:function/", checkHost(tagHandler))
	router.GET("/tag/:slug/:function/:number/", checkHost(tagHandler))
	// For serving asset files
	router.GET("/assets/*filepath", checkHost(assetsHandler))
	router.GET("/images/*filepath", checkHost(imagesHandler))
	router.GET("/content/images/*filepath", checkHost(imagesHandler)) // This is here to keep compatibility with Ghost
	router.GET("/public/*filepath", checkHost(publicHandler))
}
