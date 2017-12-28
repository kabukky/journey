package server

import (
	"github.com/dimfeld/httptreemux"
	"github.com/kabukky/journey/helpers"
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
	if err != nil || page <= 1 || number[0] == '0' {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	// Render index template
	err = templates.ShowIndexTemplate(w, r, page)
	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		errorHandler(w, r, http.StatusNotFound)
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
			// http.Error(w, err.Error(), http.StatusInternalServerError)
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		return
	} else if function == "page" {
		page, err := strconv.Atoi(number)
		if err != nil || page <= 1 || number[0] == '0' {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		// Render author template
		err = templates.ShowAuthorTemplate(w, r, slug, page)
		if err != nil {
			// http.Error(w, err.Error(), http.StatusInternalServerError)
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		return
	} else if function == "rss" {
		// Render author rss feed
		w.Header().Set("Cache-Control", "public, max-age=86400") // 1 day
		err := templates.ShowAuthorRss(w, slug)
		if err != nil {
			// http.Error(w, err.Error(), http.StatusInternalServerError)
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		return
	} else {
		errorHandler(w, r, http.StatusNotFound)
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
			// http.Error(w, err.Error(), http.StatusInternalServerError)
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		return
	} else if function == "page" {
		page, err := strconv.Atoi(number)
		if err != nil || page <= 1 || number[0] == '0' {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		// Render tag template
		err = templates.ShowTagTemplate(w, r, slug, page)
		if err != nil {
			// http.Error(w, err.Error(), http.StatusInternalServerError)
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		return
	} else if function == "rss" {
		// Render tag rss feed
		w.Header().Set("Cache-Control", "public, max-age=86400") // 1 day
		err := templates.ShowTagRss(w, slug)
		if err != nil {
			// http.Error(w, err.Error(), http.StatusInternalServerError)
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		return
	} else {
		errorHandler(w, r, http.StatusNotFound)
		return
	}
	return
}

func postHandler(w http.ResponseWriter, r *http.Request, params map[string]string) {
	slug := params["slug"]
	if slug == "rss" {
		// Render index rss feed
		w.Header().Set("Cache-Control", "public, max-age=86400") // 1 day
		err := templates.ShowIndexRss(w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	uuid := params["uuid"]
	uuidAsSlug := false
	if uuid != "" {
		slug = uuid
		uuidAsSlug = true
	}
	if slug == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	} 
	// Render post template
	err := templates.ShowPostTemplate(w, r, slug, uuidAsSlug)
	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		errorHandler(w, r, http.StatusNotFound)
		return
	}
	return
}

func assetsHandler(w http.ResponseWriter, r *http.Request, params map[string]string) {
	// Read lock global blog
	methods.Blog.RLock()
	defer methods.Blog.RUnlock()
	path := filepath.Join(filenames.ThemesFilepath, methods.Blog.ActiveTheme, "assets", params["filepath"])
	if helpers.IsDirectory(path) {
		errorHandler(w, r, http.StatusNotFound)
	}
	w.Header().Set("Cache-Control", "public, max-age=864000") // 10 days
	http.ServeFile(w, r, path)
	return
}

func imagesHandler(w http.ResponseWriter, r *http.Request, params map[string]string) {
	path := filepath.Join(filenames.ImagesFilepath, params["filepath"])
	if helpers.IsDirectory(path) {
		errorHandler(w, r, http.StatusNotFound)
	}
	w.Header().Set("Cache-Control", "public, max-age=864000") // 10 days
	http.ServeFile(w, r, path)
	return
}

func publicHandler(w http.ResponseWriter, r *http.Request, params map[string]string) {
	path := filepath.Join(filenames.PublicFilepath, params["filepath"])
	if helpers.IsDirectory(path) {
		errorHandler(w, r, http.StatusNotFound)
	}
	w.Header().Set("Cache-Control", "public, max-age=864000") // 10 days
	http.ServeFile(w, r, path)
	return
}

func faviconHandler(w http.ResponseWriter, r *http.Request, params map[string]string) {
	w.Header().Set("Cache-Control", "public, max-age=864000") // 10 days
	http.ServeFile(w, r, filepath.Join(filenames.ImagesFilepath, "favicon.ico"))
	return
}

func robotsHandler(w http.ResponseWriter, r *http.Request, params map[string]string) {
	// Read lock global blog
	methods.Blog.RLock()
	defer methods.Blog.RUnlock()
	w.Header().Set("Cache-Control", "public, max-age=864000") // 10 days
	http.ServeFile(w, r, filepath.Join(filenames.ThemesFilepath, methods.Blog.ActiveTheme, "assets", "robots.txt"))
	return
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	err := templates.ShowPostTemplate(w, r, "404", false) // TODO status might not always be 404
	if err != nil {
		http.NotFound(w, r)
		return
	}
}

func InitializeBlog(router *httptreemux.TreeMux) {
	// For index
	router.GET("/", indexHandler)
	router.GET("/favicon.ico", faviconHandler)
	router.GET("/robots.txt", robotsHandler)
	router.GET("/:slug/", postHandler)
	router.GET("/p/:uuid/", postHandler)
	router.GET("/page/:number/", indexHandler)
	// For author
	router.GET("/author/:slug/", authorHandler)
	router.GET("/author/:slug/:function/", authorHandler)
	router.GET("/author/:slug/:function/:number/", authorHandler)
	// For tag
	router.GET("/tag/:slug/", tagHandler)
	router.GET("/tag/:slug/:function/", tagHandler)
	router.GET("/tag/:slug/:function/:number/", tagHandler)
	// For serving asset files
	router.GET("/assets/*filepath", assetsHandler)
	router.GET("/images/*filepath", imagesHandler)
	// router.GET("/content/images/*filepath", imagesHandler) // This is here to keep compatibility with Ghost
	router.GET("/public/*filepath", publicHandler)
}
