package server

import (
	"github.com/dimfeld/httptreemux"
	"github.com/kabukky/journey/filenames"
	"net/http"
	"path/filepath"
)

func pagesHandler(w http.ResponseWriter, r *http.Request, params map[string]string) {
	http.ServeFile(w, r, filepath.Join(filenames.PagesFilepath, params["filepath"]))
	return
}

func InitializePages(router *httptreemux.TreeMux) {
	// For serving standalone projects or pages saved in in content/pages
	router.GET("/pages/*filepath/", pagesHandler)
}
