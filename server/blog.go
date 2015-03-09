package server

import (
	"github.com/kabukky/journey/database"
	"github.com/kabukky/journey/filenames"
	"github.com/kabukky/journey/templates"
	"github.com/kabukky/journey/timer"
	"net/http"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

var validPostPath = regexp.MustCompile("^/([\\p{L}\\p{M}\\p{N}-]*)/?$")
var validIndexPath = regexp.MustCompile("^/(page/)?([\\p{L}\\p{M}\\p{N}-]*)/?$")
var validAuthorPath = regexp.MustCompile("^/author/([\\p{L}\\p{M}\\p{N}-]*)(/page|/rss)?/?(\\d+)?/?$")
var validTagPath = regexp.MustCompile("^/tag/([\\p{L}\\p{M}\\p{N}-]*)(/page|/rss)?/?(\\d+)?/?$")

func indexHandler(w http.ResponseWriter, r *http.Request) {
	defer timer.Track(time.Now(), "index generation")
	m := validIndexPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	if m[2] == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	page, err := strconv.Atoi(m[2])
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	// Render index template
	err = templates.ShowIndexTemplate(w, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return
}

func authorHandler(w http.ResponseWriter, r *http.Request) {
	defer timer.Track(time.Now(), "author generation")
	m := validAuthorPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	if m[2] == "" {
		// Render author template (first page)
		err := templates.ShowAuthorTemplate(w, m[1], 1)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	} else if m[2] == "/rss" {
		// Render author rss feed
		err := templates.ShowAuthorRss(w, m[1])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	page, err := strconv.Atoi(m[3])
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	// Render author template
	err = templates.ShowAuthorTemplate(w, m[1], page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return
}

func tagHandler(w http.ResponseWriter, r *http.Request) {
	defer timer.Track(time.Now(), "tag generation")
	m := validTagPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	if m[2] == "" {
		// Render tag template (first page)
		err := templates.ShowTagTemplate(w, m[1], 1)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	} else if m[2] == "/rss" {
		// Render tag rss feed
		err := templates.ShowTagRss(w, m[1])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	page, err := strconv.Atoi(m[3])
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	// Render tag template
	err = templates.ShowTagTemplate(w, m[1], page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	defer timer.Track(time.Now(), "post generation")
	m := validPostPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	if m[1] == "" {
		// Render index template (first page)
		err := templates.ShowIndexTemplate(w, 1)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	} else if m[1] == "rss" {
		// Render index rss feed
		err := templates.ShowIndexRss(w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	// Render post template
	err := templates.ShowPostTemplate(w, m[1])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return
}

func generateAssetHandler(mux *http.ServeMux) {
	// Function to be able to change the asset path at runtime (e.g. when the theme changes)
	mux.Handle("/assets/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: It might be possible to do this more efficently. Getting the theme from the database at every request? Seems like too much.
		activeTheme, err := database.RetrieveActiveTheme()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.ServeFile(w, r, filepath.Join(filenames.ThemesFilepath, *activeTheme, r.URL.Path))
		return
	}))
}

func InitializeBlog(mux *http.ServeMux) {
	mux.Handle("/", http.HandlerFunc(postHandler))
	mux.Handle("/page/", http.HandlerFunc(indexHandler))
	mux.Handle("/author/", http.HandlerFunc(authorHandler))
	mux.Handle("/tag/", http.HandlerFunc(tagHandler))
	// Handle for serving asset files
	generateAssetHandler(mux)
	mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir(filenames.ImagesFilepath))))
	mux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir(filenames.PublicFilepath))))
}
