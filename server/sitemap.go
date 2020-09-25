package server

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"

	"github.com/dimfeld/httptreemux/v5"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/kabukky/journey/configuration"
	"github.com/kabukky/journey/database"
	"github.com/kabukky/journey/metrics"
	"github.com/kabukky/journey/structure"
)

func sitemapPrefix(SmURLS []structure.SmURL) []structure.SmURL {
	urls := make([]structure.SmURL, 0)
	for _, url := range SmURLS {
		url.Loc = fmt.Sprintf("%s/%s/", configuration.Config.Url, url.Loc)
		urls = append(urls, url)
	}
	return urls
}

func sitemapHandler(w http.ResponseWriter, r *http.Request, params map[string]string) {
	metrics.JourneyHandler.With(prometheus.Labels{"handler": "sitemap"}).Inc()
	var sitemap structure.Sitemap
	urls, err := database.RetrieveSitemap()
	if err != nil {
		log.Printf("%s", err)
	}
	urls = sitemapPrefix(urls)
	sitemap.URLs = urls
	sitemap.Xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"
	output, _ := xml.MarshalIndent(sitemap, " ", "  ")

	w.Header().Set("Content-Type", "application/xml")
	_, _ = w.Write(output)
}

func InitializeSitemap(router *httptreemux.TreeMux) {
	router.GET("/sitemap.xml", sitemapHandler)
}
