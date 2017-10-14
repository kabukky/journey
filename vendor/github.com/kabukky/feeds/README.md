## gorilla/feeds

Web feed generator library.

[![Build Status](https://travis-ci.org/gorilla/feeds.png?branch=master)](https://travis-ci.org/gorilla/feeds)

### Goals

 * simple interface to create both Atom & RSS 2.0 feeds
 * full support for Atom/RSS2.0 spec elements
 * ability to modify particulars for each spec

### Usage

```go

import (
    "fmt"
    "time"
    "github.com/gorilla/feeds"
)

now := time.Now()
feed := &feeds.Feed{
    Title:       "jmoiron.net blog",
    Link:        &feeds.Link{Href: "http://jmoiron.net/blog"},
    Description: "discussion about tech, footie, photos",
    Author:      &feeds.Author{"Jason Moiron", "jmoiron@jmoiron.net"},
    Created:     now,
}

feed.Items = []*feeds.Item{
    &feeds.Item{
        Title:       "Limiting Concurrency in Go",
        Link:        &feeds.Link{Href: "http://jmoiron.net/blog/limiting-concurrency-in-go/"},
        Description: "A discussion on controlled parallelism in golang",
        Author:      &feeds.Author{"Jason Moiron", "jmoiron@jmoiron.net"},
        Created:     now,
    },
    &feeds.Item{
        Title:       "Logic-less Template Redux",
        Link:        &feeds.Link{Href: "http://jmoiron.net/blog/logicless-template-redux/"},
        Description: "More thoughts on logicless templates",
        Created:     now,
    },
    &feeds.Item{
        Title:       "Idiomatic Code Reuse in Go",
        Link:        &feeds.Link{Href: "http://jmoiron.net/blog/idiomatic-code-reuse-in-go/"},
        Description: "How to use interfaces <em>effectively</em>",
        Created:     now,
    },
}

atom, err := feed.ToAtom()
rss, err := feed.ToRss()

fmt.Println(atom, "\n", rss)

```

Outputs:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  <title>jmoiron.net blog</title>
  <link href="http://jmoiron.net/blog"></link>
  <id>http://jmoiron.net/blog</id>
  <updated>2013-01-16T03:26:01-05:00</updated>
  <summary>discussion about tech, footie, photos</summary>
  <entry>
    <title>Limiting Concurrency in Go</title>
    <link href="http://jmoiron.net/blog/limiting-concurrency-in-go/"></link>
    <updated>2013-01-16T03:26:01-05:00</updated>
    <id>tag:jmoiron.net,2013-01-16:/blog/limiting-concurrency-in-go/</id>
    <summary type="html">A discussion on controlled parallelism in golang</summary>
    <author>
      <name>Jason Moiron</name>
      <email>jmoiron@jmoiron.net</email>
    </author>
  </entry>
  <entry>
    <title>Logic-less Template Redux</title>
    <link href="http://jmoiron.net/blog/logicless-template-redux/"></link>
    <updated>2013-01-16T03:26:01-05:00</updated>
    <id>tag:jmoiron.net,2013-01-16:/blog/logicless-template-redux/</id>
    <summary type="html">More thoughts on logicless templates</summary>
    <author></author>
  </entry>
  <entry>
    <title>Idiomatic Code Reuse in Go</title>
    <link href="http://jmoiron.net/blog/idiomatic-code-reuse-in-go/"></link>
    <updated>2013-01-16T03:26:01-05:00</updated>
    <id>tag:jmoiron.net,2013-01-16:/blog/idiomatic-code-reuse-in-go/</id>
    <summary type="html">How to use interfaces &lt;em&gt;effectively&lt;/em&gt;</summary>
    <author></author>
  </entry>
</feed>

<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>jmoiron.net blog</title>
    <link>http://jmoiron.net/blog</link>
    <description>discussion about tech, footie, photos</description>
    <managingEditor>jmoiron@jmoiron.net (Jason Moiron)</managingEditor>
    <pubDate>2013-01-16T03:22:24-05:00</pubDate>
    <item>
      <title>Limiting Concurrency in Go</title>
      <link>http://jmoiron.net/blog/limiting-concurrency-in-go/</link>
      <description>A discussion on controlled parallelism in golang</description>
      <pubDate>2013-01-16T03:22:24-05:00</pubDate>
    </item>
    <item>
      <title>Logic-less Template Redux</title>
      <link>http://jmoiron.net/blog/logicless-template-redux/</link>
      <description>More thoughts on logicless templates</description>
      <pubDate>2013-01-16T03:22:24-05:00</pubDate>
    </item>
    <item>
      <title>Idiomatic Code Reuse in Go</title>
      <link>http://jmoiron.net/blog/idiomatic-code-reuse-in-go/</link>
      <description>How to use interfaces &lt;em&gt;effectively&lt;/em&gt;</description>
      <pubDate>2013-01-16T03:22:24-05:00</pubDate>
    </item>
  </channel>
</rss>

```

