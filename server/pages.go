package server

import (
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/dimfeld/httptreemux/v5"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/kabukky/journey/filenames"
	"github.com/kabukky/journey/helpers"
	"github.com/kabukky/journey/metrics"
	"github.com/kabukky/journey/templates"
)

func pagesHandler(w http.ResponseWriter, r *http.Request, params map[string]string) {
	metrics.JourneyHandler.With(prometheus.Labels{"handler": "pages"}).Inc()
	path := filepath.Join(filenames.PagesFilepath, params["filepath"])
	// If the path points to a directory, add a trailing slash to the path (needed if the page loads relative assets).
	if helpers.IsDirectory(path) && !strings.HasSuffix(r.RequestURI, "/") {
		http.Redirect(w, r, r.RequestURI+"/", 301)
		return
	}
	if !helpers.FileExists(path) {
		e404 := templates.ShowPostTemplate(w, r, "404")
		if e404 != nil {
			http.Error(w, "Nobody here but us chickens!", http.StatusNotFound)
			log.Println("404:", r.URL)
		}
		return
	}
	http.ServeFile(w, r, path)

}

// InitializePages serving standalone projects or pages saved in in content/pages
func InitializePages(router *httptreemux.TreeMux) {
	router.GET("/pages/*filepath", pagesHandler)
}
