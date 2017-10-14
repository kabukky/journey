package feeds

import (
	"encoding/xml"
	"io"
	"strings"
	"time"
)

type Image struct {
	Url, Title, Link string
	Width, Height    int
}

type Link struct {
	Href, Rel string
}

type Author struct {
	Name, Email string
}

type Item struct {
	Title       string
	Link        *Link
	Author      *Author
	Description string // used as description in rss, summary in atom
	Id          string // used as guid in rss, id in atom
	Updated     time.Time
	Created     time.Time
	Image       *Image
}

type Feed struct {
	Title       string
	Link        *Link
	Description string
	Author      *Author
	Updated     time.Time
	Created     time.Time
	Id          string
	Subtitle    string
	Items       []*Item
	Copyright   string
	Image       *Image
	Url         string
}

// add a new Item to a Feed
func (f *Feed) Add(item *Item) {
	f.Items = append(f.Items, item)
}

// returns the first non-zero time formatted as a string or ""
func anyTimeFormat(format string, times ...time.Time) string {
	for _, t := range times {
		if !t.IsZero() {
			// Always return GMT time by converting to UTC and then replacing UTC with GMT in the output string (RSS doesn't allow UTC)
			timeFormatted := t.UTC().Format(format)
			return strings.Replace(timeFormatted, "UTC", "GMT", -1)
		}
	}
	return ""
}

// interface used by ToXML to get a object suitable for exporting XML.
type XmlFeed interface {
	FeedXml() interface{}
}

// turn a feed object (either a Feed, AtomFeed, or RssFeed) into xml
// returns an error if xml marshaling fails
func ToXML(feed XmlFeed) (string, error) {
	x := feed.FeedXml()
	data, err := xml.MarshalIndent(x, "", "  ")
	if err != nil {
		return "", err
	}
	// strip empty line from default xml header
	s := xml.Header[:len(xml.Header)-1] + string(data)
	return s, nil
}

// Write a feed object (either a Feed, AtomFeed, or RssFeed) as XML into
// the writer. Returns an error if XML marshaling fails.
func WriteXML(feed XmlFeed, w io.Writer) error {
	x := feed.FeedXml()
	// write default xml header, without the newline
	if _, err := w.Write([]byte(xml.Header)); err != nil {
		return err
	}
	e := xml.NewEncoder(w)
	e.Indent("", "  ")
	return e.Encode(x)
}

// creates an Atom representation of this feed
func (f *Feed) ToAtom() (string, error) {
	a := &Atom{f}
	return ToXML(a)
}

// Writes an Atom representation of this feed to the writer.
func (f *Feed) WriteAtom(w io.Writer) error {
	return WriteXML(&Atom{f}, w)
}

// creates an Rss representation of this feed
func (f *Feed) ToRss() (string, error) {
	r := &Rss{f}
	return ToXML(r)
}

// Writes an RSS representation of this feed to the writer.
func (f *Feed) WriteRss(w io.Writer) error {
	return WriteXML(&Rss{f}, w)
}
