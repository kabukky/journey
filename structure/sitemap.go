package structure

import (
	"encoding/xml"
)

type SmURL struct {
	Loc     string `xml:"loc"`
	LastMod string `xml:"lastmod,omitempty"`
}

type Sitemap struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr" default:"http://www.sitemaps.org/schemas/sitemap/0.9"`

	URLs []SmURL `xml:"url"`

	Minify bool `xml:"-"`
}
