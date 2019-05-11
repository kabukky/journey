package server

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/dimfeld/httptreemux"

	"github.com/kabukky/journey/filenames"
	"github.com/kabukky/journey/helpers"
)

func pagesHandler(w http.ResponseWriter, r *http.Request, params map[string]string) {
	path := filepath.Join(filenames.PagesFilepath, params["filepath"])
	// If the path points to a directory, add a trailing slash to the path (needed if the page loads relative assets).
	if helpers.IsDirectory(path) && !strings.HasSuffix(r.RequestURI, "/") {
		http.Redirect(w, r, r.RequestURI+"/", 301)
		return
	}
	if !helpers.FileExists(path) {
		http.Error(w, "Still lost?", http.StatusNotFound)
		return
	}
	http.ServeFile(w, r, path)

}

// InitializePages serving standalone projects or pages saved in in content/pages
func InitializePages(router *httptreemux.TreeMux) {
	router.GET("/pages/*filepath", pagesHandler)
}
