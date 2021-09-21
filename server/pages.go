package server

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/Landria/journey/filenames"
	"github.com/Landria/journey/helpers"
	"github.com/dimfeld/httptreemux"
)

func pagesHandler(w http.ResponseWriter, r *http.Request, params map[string]string) {
	path := filepath.Join(filenames.PagesFilepath, params["filepath"])
	// If the path points to a directory, add a trailing slash to the path (needed if the page loads relative assets).
	if helpers.IsDirectory(path) && !strings.HasSuffix(r.RequestURI, "/") {
		http.Redirect(w, r, r.RequestURI+"/", 301)
		return
	}
	http.ServeFile(w, r, path)
	return
}

func InitializePages(router *httptreemux.TreeMux) {
	// For serving standalone projects or pages saved in in content/pages
	router.GET("/pages/*filepath", pagesHandler)
}
