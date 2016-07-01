package templates

import (
	"bytes"
	"encoding/json"
	"github.com/kabukky/journey/conversion"
	"github.com/kabukky/journey/database"
	"github.com/kabukky/journey/date"
	"github.com/kabukky/journey/plugins"
	"github.com/kabukky/journey/structure"
	"github.com/kabukky/journey/structure/methods"
	"html"
	"log"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// Helper fuctions
func nullFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	// Check if the helper was defined in a plugin
	if plugins.LuaPool != nil {
		// Get a state map to execute and attach it to the request data
		if values.PluginVMs == nil {
			values.PluginVMs = plugins.LuaPool.Get(helper, values)
		}
		if values.PluginVMs[helper.Name] != nil {
			pluginResult, err := plugins.Execute(helper, values)
			if err != nil {
				return []byte{}
			}
			return evaluateEscape(pluginResult, helper.Unescaped)
		} else {
			// This helper is not implemented in a plugin. Get rid of the Lua VMs
			plugins.LuaPool.Put(values.PluginVMs)
			values.PluginVMs = nil
		}
	}
	log.Println("Warning: This helper is not implemented:", helper.Name)
	return []byte{}
}

func slugFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	if len(values.Blog.NavigationItems) != 0 {
		return evaluateEscape([]byte(values.Blog.NavigationItems[values.CurrentNavigationIndex].Slug), helper.Unescaped)
	}
	return []byte{}
}

var pageInUrlRegex = regexp.MustCompile("/page/[0-9]+/$")

func currentFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	if len(values.Blog.NavigationItems) != 0 {
		url := values.Blog.NavigationItems[values.CurrentNavigationIndex].Url
		// Since the router rewrites all urls with a trailing slash, add / to url if not already there
		if !strings.HasSuffix(url, "/") {
			url = url + "/"
		}
		currentPath := pageInUrlRegex.ReplaceAllString(values.CurrentPath, "/")
		if currentPath == url {
			return []byte{1}
		}
	}
	return []byte{}
}

func navigationFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	if len(values.Blog.NavigationItems) == 0 {
		return []byte{}
	} else if templateHelper, ok := compiledTemplates.m["navigation"]; ok {
		return executeHelper(templateHelper, values, values.CurrentHelperContext)
	}
	return []byte{}
}

func labelFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	if len(values.Blog.NavigationItems) != 0 {
		return evaluateEscape([]byte(values.Blog.NavigationItems[values.CurrentNavigationIndex].Label), helper.Unescaped)
	}
	return []byte{}
}

func contentForFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	// If there is no array attached to the request data already, make one
	if values.ContentForHelpers == nil {
		values.ContentForHelpers = make([]structure.Helper, 0)
	}
	// Collect all contentFor helpers to use them with a block helper
	values.ContentForHelpers = append(values.ContentForHelpers, *helper)
	return []byte{}
}

func blockFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	if len(helper.Arguments) != 0 {
		// Loop through the collected contentFor helpers and execute the appropriate one
		for index, _ := range values.ContentForHelpers {
			if len(values.ContentForHelpers[index].Arguments) != 0 {
				if values.ContentForHelpers[index].Arguments[0].Name == helper.Arguments[0].Name {
					return executeHelper(&values.ContentForHelpers[index], values, values.CurrentHelperContext)
				}
			}
		}
	}
	return []byte{}
}

func paginationFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	if templateHelper, ok := compiledTemplates.m["pagination"]; ok {
		return executeHelper(templateHelper, values, values.CurrentHelperContext)
	}
	return []byte{}
}

func paginationDotTotalFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	if values.CurrentTemplate == 0 { // index
		return []byte(strconv.FormatInt(values.Blog.PostCount, 10))
	} else if values.CurrentTemplate == 3 { // author
		count, err := database.RetrieveNumberOfPostsByUser(values.Posts[values.CurrentPostIndex].Author.Id)
		if err != nil {
			log.Println("Couldn't get number of posts", err.Error())
			return []byte{}
		}
		return []byte(strconv.FormatInt(count, 10))
	} else if values.CurrentTemplate == 2 { // tag
		count, err := database.RetrieveNumberOfPostsByTag(values.CurrentTag.Id)
		if err != nil {
			log.Println("Couldn't get number of posts", err.Error())
			return []byte{}
		}
		return []byte(strconv.FormatInt(count, 10))
	}
	return []byte{}
}

func pluralFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	if len(helper.Arguments) != 0 {
		// Get the number calculated by executing the first argument
		countString := string(helper.Arguments[0].Function(helper, values))
		if countString == "" {
			log.Println("Couldn't get count in plural helper")
			return []byte{}
		}
		arguments := methods.ProcessHelperArguments(helper.Arguments)
		for key, value := range arguments {
			if countString == "0" && key == "empty" {
				output := value
				output = strings.Replace(output, "%", countString, -1)
				return []byte(output)
			} else if countString == "1" && key == "singular" {
				output := value
				output = strings.Replace(output, "%", countString, -1)
				return []byte(output)
			} else if countString != "0" && countString != "1" && key == "plural" {
				output := value
				output = strings.Replace(output, "%", countString, -1)
				return []byte(output)
			}
		}
	}
	return []byte{}
}

func prevFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	if values.CurrentIndexPage > 1 {
		return []byte{1}
	}
	return []byte{}
}

func nextFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	var count int64
	var err error
	if values.CurrentTemplate == 0 { // index
		count = values.Blog.PostCount
	} else if values.CurrentTemplate == 2 { // tag
		count, err = database.RetrieveNumberOfPostsByTag(values.CurrentTag.Id)
		if err != nil {
			log.Println("Couldn't get number of posts for tag", err.Error())
			return []byte{}
		}
	} else if values.CurrentTemplate == 3 { // author
		count, err = database.RetrieveNumberOfPostsByUser(values.Posts[values.CurrentPostIndex].Author.Id)
		if err != nil {
			log.Println("Couldn't get number of posts for author", err.Error())
			return []byte{}
		}
	}
	maxPages := positiveCeilingInt64(float64(count) / float64(values.Blog.PostsPerPage))
	if int64(values.CurrentIndexPage) < maxPages {
		return []byte{1}
	}
	return []byte{}
}

func pageFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	return []byte(strconv.Itoa(values.CurrentIndexPage))
}

func pagesFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	var count int64
	var err error
	if values.CurrentTemplate == 0 { // index
		count = values.Blog.PostCount
	} else if values.CurrentTemplate == 2 { // tag
		count, err = database.RetrieveNumberOfPostsByTag(values.CurrentTag.Id)
		if err != nil {
			log.Println("Couldn't get number of posts for tag", err.Error())
			return []byte{}
		}
	} else if values.CurrentTemplate == 3 { // author
		count, err = database.RetrieveNumberOfPostsByUser(values.Posts[values.CurrentPostIndex].Author.Id)
		if err != nil {
			log.Println("Couldn't get number of posts for author", err.Error())
			return []byte{}
		}
	}
	maxPages := positiveCeilingInt64(float64(count) / float64(values.Blog.PostsPerPage))
	// Output at least 1 (even if there are no posts in the database)
	if maxPages == 0 {
		maxPages = 1
	}
	return []byte(strconv.FormatInt(maxPages, 10))
}

func page_urlFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	if len(helper.Arguments) != 0 {
		if helper.Arguments[0].Name == "prev" || helper.Arguments[0].Name == "pagination.prev" {
			if values.CurrentIndexPage > 1 {
				var buffer bytes.Buffer
				if values.CurrentIndexPage == 2 {
					if values.CurrentTemplate == 3 { // author
						buffer.WriteString("/author/")
						//TODO: Error handling if there is no Posts[values.CurrentPostIndex]
						buffer.WriteString(values.Posts[values.CurrentPostIndex].Author.Slug)
					} else if values.CurrentTemplate == 2 { // tag
						buffer.WriteString("/tag/")
						//TODO: Error handling if there is no Posts[values.CurrentPostIndex]
						buffer.WriteString(values.CurrentTag.Slug)
					}
					buffer.WriteString("/")
				} else {
					if values.CurrentTemplate == 3 { // author
						buffer.WriteString("/author/")
						//TODO: Error handling if there is no Posts[values.CurrentPostIndex]
						buffer.WriteString(values.Posts[values.CurrentPostIndex].Author.Slug)
					} else if values.CurrentTemplate == 2 { // tag
						buffer.WriteString("/tag/")
						//TODO: Error handling if there is no Posts[values.CurrentPostIndex]
						buffer.WriteString(values.CurrentTag.Slug)
					}
					page := values.CurrentIndexPage - 1
					if page > 1 {
						buffer.WriteString("/page/")
						buffer.WriteString(strconv.Itoa(page))
					}
					buffer.WriteString("/")
				}
				return buffer.Bytes()
			}
		} else if helper.Arguments[0].Name == "next" || helper.Arguments[0].Name == "pagination.next" {
			var count int64
			var err error
			if values.CurrentTemplate == 0 { // index
				count = values.Blog.PostCount
			} else if values.CurrentTemplate == 2 { // tag
				count, err = database.RetrieveNumberOfPostsByTag(values.CurrentTag.Id)
				if err != nil {
					log.Println("Couldn't get number of posts for tag", err.Error())
					return []byte{}
				}
			} else if values.CurrentTemplate == 3 { // author
				count, err = database.RetrieveNumberOfPostsByUser(values.Posts[values.CurrentPostIndex].Author.Id)
				if err != nil {
					log.Println("Couldn't get number of posts for author", err.Error())
					return []byte{}
				}
			}
			maxPages := positiveCeilingInt64(float64(count) / float64(values.Blog.PostsPerPage))
			if int64(values.CurrentIndexPage) < maxPages {
				var buffer bytes.Buffer
				if values.CurrentTemplate == 3 { // author
					buffer.WriteString("/author/")
					// TODO: Error handling if there is no Posts[values.CurrentPostIndex]
					buffer.WriteString(values.Posts[values.CurrentPostIndex].Author.Slug)
				} else if values.CurrentTemplate == 2 { // tag
					buffer.WriteString("/tag/")
					// TODO: Error handling if there is no Posts[values.CurrentPostIndex]
					buffer.WriteString(values.CurrentTag.Slug)
				}
				page := values.CurrentIndexPage + 1
				if page > 1 {
					buffer.WriteString("/page/")
					buffer.WriteString(strconv.Itoa(page))
				}
				buffer.WriteString("/")
				return buffer.Bytes()
			}
		}
	}
	return []byte{}
}

func extendFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	if len(helper.Arguments) != 0 {
		return []byte(helper.Arguments[0].Name)
	}
	return []byte{}
}

func featuredFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	if values.Posts[values.CurrentPostIndex].IsFeatured {
		return []byte{1}
	}
	return []byte{}
}

func publishedFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	if values.Posts[values.CurrentPostIndex].IsPublished {
		return []byte{1}
	}
	return []byte{}
}

func body_classFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	if values.CurrentTemplate == 1 { // post
		// TODO: is there anything else that needs to be output here?
		var buffer bytes.Buffer
		buffer.WriteString("post-template")
		// If page
		if values.Posts[values.CurrentPostIndex].IsPage {
			buffer.WriteString(" page-template page")
		}
		for _, tag := range values.Posts[values.CurrentPostIndex].Tags {
			buffer.WriteString(" tag-")
			buffer.WriteString(tag.Slug)
		}
		return buffer.Bytes()
	} else if values.CurrentTemplate == 0 { // index
		if values.CurrentIndexPage == 1 {
			return []byte("home-template")
		} else {
			return []byte("paged archive-template")
		}
	} else if values.CurrentTemplate == 3 { // author
		var buffer bytes.Buffer
		buffer.WriteString("author-template author-")
		// TODO: Error handling if there is no Posts[values.CurrentPostIndex]
		buffer.WriteString(values.Posts[values.CurrentPostIndex].Author.Slug)
		if values.CurrentIndexPage > 1 {
			buffer.WriteString(" paged archive-template")
		}
		return buffer.Bytes()
	} else if values.CurrentTemplate == 2 { // tag
		var buffer bytes.Buffer
		buffer.WriteString("tag-template tag-")
		buffer.WriteString(values.CurrentTag.Slug)
		if values.CurrentIndexPage > 1 {
			buffer.WriteString(" paged archive-template")
		}
		return buffer.Bytes()
	}
	// TODO: Delete this. Probably not needed.
	return []byte("post-template")
}

func ghost_headFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	// SEO stuff:
	currentUrl := string(evaluateEscape(values.Blog.Url, helper.Unescaped)) + values.CurrentPath
	// Output canonical url
	var buffer bytes.Buffer
	buffer.WriteString("<link rel=\"canonical\" href=\"")
	buffer.WriteString(currentUrl)
	buffer.WriteString("\">\n")
	// Output structured data
	// values.CurrentPostIndex = 0 // current
	structuredData := map[string]string{
		"og:site_name": string(evaluateEscape(values.Blog.Title, helper.Unescaped)),
		"og:type": "website",
		"og:title": string(meta_titleFunc(helper, values)),
		"og:description": string(meta_descriptionFunc(helper, values)),
		"og:url": currentUrl,
		"og:image": string(imageFunc(helper, values)),
		"twitter:card": "summary", // summary or summary_large_image
		"twitter:title": string(meta_titleFunc(helper, values)),
		"twitter:description": string(meta_descriptionFunc(helper, values)),
		"twitter:url": currentUrl,
		"twitter:image:src": string(imageFunc(helper, values)),
	}
	schema := map[string]string{
		"@context": "http://schema.org",
		"@type": "Website",
		"publisher": string(evaluateEscape(values.Blog.Title, helper.Unescaped)),
		"headline": string(meta_titleFunc(helper, values)),
		"url": currentUrl,
		"image": string(imageFunc(helper, values)),
		// "keywords": ,
		"description": string(meta_descriptionFunc(helper, values)),
	};
	if values.CurrentTemplate == 1 { // post
		publicationDate := values.Posts[values.CurrentPostIndex].Date.Format("2006-01-02T15:04:05Z")
		structuredData["og:type"] = "article"
		structuredData["article:published_time"] = values.Posts[values.CurrentPostIndex].Date.Format("2006-01-02T15:04:05Z")
		// structuredData["article:modified_time"]
		// structuredData["article:tag"]
		schema["@type"] = "Article"
		schema["datePublished"] = publicationDate
		// schema["dateModified"]
	} else if values.CurrentTemplate == 2 { // tag
		schema["@type"] = "Series"
	} else if values.CurrentTemplate == 3 { // author
		structuredData["og:type"] = "profile"
		schema["@type"] = "Person"
	}
	structuredDataKeys := make([]string, len(structuredData))
	i := 0
	for key, _ := range structuredData {
		structuredDataKeys[i] = key
		i++
	}
	sort.Strings(structuredDataKeys)
	for _, key := range structuredDataKeys {
		buffer.WriteString("<meta property=\"")
		buffer.WriteString(key)
		buffer.WriteString("\" content=\"")
		buffer.WriteString(structuredData[key])
		buffer.WriteString("\">\n")
	}
	buffer.WriteString("<script type=\"application/ld+json\">\n")
	json, _ := json.MarshalIndent(schema, "", "    ")
	buffer.Write(json)
	buffer.WriteString("\n</script>\n")
	return buffer.Bytes()
}

func ghost_footFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	// TODO: customized code injection
	return []byte{}
}

func meta_titleFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	if values.CurrentTemplate == 1 { // post or page
		return evaluateEscape(values.Posts[values.CurrentPostIndex].Title, helper.Unescaped)
	} else if values.CurrentTemplate == 3 { // author
		var buffer bytes.Buffer
		// TODO: Error handling if there is no Posts[values.CurrentPostIndex]
		buffer.Write(values.Posts[values.CurrentPostIndex].Author.Name)
		buffer.WriteString(" - ")
		buffer.Write(values.Blog.Title)
		return evaluateEscape(buffer.Bytes(), helper.Unescaped)
	} else if values.CurrentTemplate == 2 { // tag
		var buffer bytes.Buffer
		// TODO: Error handling if there is no Posts[values.CurrentPostIndex]
		buffer.Write(values.CurrentTag.Name)
		buffer.WriteString(" - ")
		buffer.Write(values.Blog.Title)
		return evaluateEscape(buffer.Bytes(), helper.Unescaped)
	}
	// index
	return evaluateEscape(values.Blog.Title, helper.Unescaped)
}

func meta_descriptionFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	// TODO: Finish this
	if values.CurrentTemplate == 1 || values.CurrentHelperContext == 1 { // post
		return evaluateEscape(values.Posts[values.CurrentPostIndex].MetaDescription, helper.Unescaped)
	} else {
		return evaluateEscape(values.Blog.Description, helper.Unescaped)
	}
	// Nothing on post yet
	return []byte{}
}

func bodyFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	return helper.Block
}

func insertFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	if len(helper.Arguments) != 0 {
		if templateHelper, ok := compiledTemplates.m[helper.Arguments[0].Name]; ok {
			return executeHelper(templateHelper, values, values.CurrentHelperContext)
		}
	}
	return []byte{}
}

func encodeFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	if len(helper.Arguments) != 0 {
		return []byte(url.QueryEscape(string(helper.Arguments[0].Function(&helper.Arguments[0], values))))
	}
	return []byte{}
}

func authorFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	// Check if helper is block helper
	if len(helper.Block) != 0 {
		return executeHelper(helper, values, 3) // context = author
	}
	// Else return author name (as link)
	arguments := methods.ProcessHelperArguments(helper.Arguments)
	for key, value := range arguments {
		// If link is set to false, just return the name
		if key == "autolink" && value == "false" {
			return evaluateEscape(values.Posts[values.CurrentPostIndex].Author.Name, helper.Unescaped)
		}
	}
	var buffer bytes.Buffer
	buffer.WriteString("<a href=\"")
	buffer.WriteString("/author/")
	// TODO: Error handling if there i no Posts[values.CurrentPostIndex]
	buffer.WriteString(values.Posts[values.CurrentPostIndex].Author.Slug)
	buffer.WriteString("/\">")
	// TODO: Error handling if there is no Posts[values.CurrentPostIndex]
	buffer.Write(evaluateEscape(values.Posts[values.CurrentPostIndex].Author.Name, helper.Unescaped))
	buffer.WriteString("</a>")
	return buffer.Bytes()
}

func authorDotNameFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	return evaluateEscape(values.Posts[values.CurrentPostIndex].Author.Name, helper.Unescaped)
}

func bioFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	// TODO: Error handling if there is no Posts[values.CurrentPostIndex]
	return evaluateEscape(values.Posts[values.CurrentPostIndex].Author.Bio, helper.Unescaped)
}

func emailFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	// TODO: Error handling if there is no Posts[values.CurrentPostIndex]
	return evaluateEscape(values.Posts[values.CurrentPostIndex].Author.Email, helper.Unescaped)
}

func websiteFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	// TODO: Error handling if there is no Posts[values.CurrentPostIndex]
	return evaluateEscape(values.Posts[values.CurrentPostIndex].Author.Website, helper.Unescaped)
}

func imageFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	if values.CurrentHelperContext == 1 { // post
		// TODO: Error handling if there is no Posts[values.CurrentPostIndex]
		return evaluateEscape(values.Posts[values.CurrentPostIndex].Image, helper.Unescaped)
	} else if values.CurrentHelperContext == 3 { // author
		// TODO: Error handling if there is no Posts[values.CurrentPostIndex]
		return evaluateEscape(values.Posts[values.CurrentPostIndex].Author.Image, helper.Unescaped)
	}
	return []byte{}
}

func authorDotImageFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	// TODO: Error handling if there is no Posts[values.CurrentPostIndex]
	return evaluateEscape(values.Posts[values.CurrentPostIndex].Author.Image, helper.Unescaped)
}

func coverFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	// TODO: Error handling if there is no Posts[values.CurrentPostIndex]
	return evaluateEscape(values.Posts[values.CurrentPostIndex].Author.Cover, helper.Unescaped)
}

func locationFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	// TODO: Error handling if there is no Posts[values.CurrentPostIndex]
	return evaluateEscape(values.Posts[values.CurrentPostIndex].Author.Location, helper.Unescaped)
}

func twitterFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	// TODO: Error handling if there is no Posts[values.CurrentPostIndex]
	return evaluateEscape(values.Posts[values.CurrentPostIndex].Author.Twitter, helper.Unescaped)
}

func facebookFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	// TODO: Error handling if there is no Posts[values.CurrentPostIndex]
	return evaluateEscape(values.Posts[values.CurrentPostIndex].Author.Facebook, helper.Unescaped)
}

func postFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	values.CurrentPostIndex = 0 // the current post
	return executeHelper(helper, values, 1) // context = post
}

func prevPostFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	if values.Posts[1].Id == 0 {
		return []byte{}
	}
	values.CurrentPostIndex = 1 // the previous post
	var buffer bytes.Buffer
	buffer.Write(executeHelper(helper, values, 1)) // context = post
	values.CurrentPostIndex = 0 // the current post
	return buffer.Bytes()
}

func nextPostFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	if values.Posts[2].Id == 0 {
		return []byte{}
	}
	values.CurrentPostIndex = 2 // the next post
	var buffer bytes.Buffer
	buffer.Write(executeHelper(helper, values, 1)) // context = post
	values.CurrentPostIndex = 0 // the current post
	return buffer.Bytes()
}

func postsFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	if len(values.Posts) > 0 {
		return []byte{1}
	}
	return []byte{}
}

func tagsFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	if len(values.Posts[values.CurrentPostIndex].Tags) > 0 {
		separator := ", "
		suffix := ""
		prefix := ""
		makeLink := true
		if len(helper.Arguments) != 0 {
			arguments := methods.ProcessHelperArguments(helper.Arguments)
			for key, value := range arguments {
				if key == "separator" {
					separator = value
				} else if key == "suffix" {
					suffix = value
				} else if key == "prefix" {
					prefix = value
				} else if key == "autolink" {
					if value == "false" {
						makeLink = false
					}
				}
			}
		}
		var buffer bytes.Buffer
		if prefix != "" {
			buffer.WriteString(prefix)
			buffer.WriteString(" ")
		}
		for index, tag := range values.Posts[values.CurrentPostIndex].Tags {
			if index != 0 {
				buffer.WriteString(separator)
			}
			if makeLink {
				buffer.WriteString("<a href=\"")
				buffer.WriteString("/tag/")
				buffer.WriteString(tag.Slug)
				buffer.WriteString("/\">")
			}
			buffer.Write(evaluateEscape(tag.Name, helper.Unescaped))
			if makeLink {
				buffer.WriteString("</a>")
			}
		}
		if suffix != "" {
			buffer.WriteString(" ")
			buffer.WriteString(suffix)
		}
		return buffer.Bytes()
	}
	return []byte{}
}

func post_classFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	var buffer bytes.Buffer
	buffer.WriteString("post")
	if values.Posts[values.CurrentPostIndex].IsFeatured {
		buffer.WriteString(" featured")
	}
	if values.Posts[values.CurrentPostIndex].IsPage {
		buffer.WriteString(" page")
	}
	for _, tag := range values.Posts[values.CurrentPostIndex].Tags {
		buffer.WriteString(" tag-")
		buffer.WriteString(tag.Slug)
	}
	return evaluateEscape(buffer.Bytes(), helper.Unescaped)
}

func urlFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	var buffer bytes.Buffer
	if len(helper.Arguments) != 0 {
		arguments := methods.ProcessHelperArguments(helper.Arguments)
		for key, value := range arguments {
			if key == "absolute" {
				if value == "true" {
					// Only write the blog url if navigation url does not begin with http/https
					if values.CurrentHelperContext == 4 && (!strings.HasPrefix(values.Blog.NavigationItems[values.CurrentNavigationIndex].Url, "http://") && !strings.HasPrefix(values.Blog.NavigationItems[values.CurrentNavigationIndex].Url, "https://")) { // navigation
						buffer.Write(values.Blog.Url)
					} else if values.CurrentHelperContext != 4 {
						buffer.Write(values.Blog.Url)
					}
				}
			}
		}
	}
	if values.CurrentHelperContext == 1 { // post
		buffer.WriteString("/")
		buffer.WriteString(values.Posts[values.CurrentPostIndex].Slug)
		buffer.WriteString("/")
		return evaluateEscape(buffer.Bytes(), helper.Unescaped)
	} else if values.CurrentHelperContext == 3 { // author
		buffer.WriteString("/author/")
		// TODO: Error handling if there is no Posts[values.CurrentPostIndex]
		buffer.WriteString(values.Posts[values.CurrentPostIndex].Author.Slug)
		buffer.WriteString("/")
		return evaluateEscape(buffer.Bytes(), helper.Unescaped)
	} else if values.CurrentHelperContext == 4 { // navigation
		buffer.WriteString(values.Blog.NavigationItems[values.CurrentNavigationIndex].Url)
		return evaluateEscape(buffer.Bytes(), helper.Unescaped)
	}
	return []byte{}
}

func titleFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	return evaluateEscape(values.Posts[values.CurrentPostIndex].Title, helper.Unescaped)
}

func contentFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	// TODO: is content always unescaped? seems like it...
	return values.Posts[values.CurrentPostIndex].Html
}

func excerptFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	if values.CurrentHelperContext == 1 { // post
		if len(helper.Arguments) != 0 {
			arguments := methods.ProcessHelperArguments(helper.Arguments)
			for key, value := range arguments {
				if key == "words" {
					number, err := strconv.Atoi(value)
					if err == nil {
						excerpt := conversion.StripTagsFromHtml(values.Posts[values.CurrentPostIndex].Html)
						words := bytes.Fields(excerpt)
						if len(words) < number {
							return excerpt
						}
						return bytes.Join(words[:number], []byte(" "))
					}
				} else if key == "characters" {
					number, err := strconv.Atoi(value)
					if err == nil {
						// Use runes for UTF-8 support
						runes := []rune(string(conversion.StripTagsFromHtml(values.Posts[values.CurrentPostIndex].Html)))
						if len(runes) < number {
							return []byte(string(runes))
						}
						return []byte(string(runes[:number]))
					}
				}
			}
		}
		// Default to 50 words excerpt
		excerpt := conversion.StripTagsFromHtml(values.Posts[values.CurrentPostIndex].Html)
		words := bytes.Fields(excerpt)
		if len(words) < 50 {
			return excerpt
		}
		return bytes.Join(words[:50], []byte(" "))
	}
	return []byte{}
}

func dateFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	showPublicationDate := false
	timeFormat := "MMM Do, YYYY" // Default time format
	// If in scope of a post, change default to published date
	if values.CurrentHelperContext == 1 { // post
		showPublicationDate = true
	}
	// Get the date
	if len(helper.Arguments) != 0 {
		arguments := methods.ProcessHelperArguments(helper.Arguments)
		for key, value := range arguments {
			if key == "published_at" {
				showPublicationDate = true
			} else if key == "timeago" {
				if value == "true" {
					// Compute time ago
					return evaluateEscape(date.GenerateTimeAgo(values.Posts[values.CurrentPostIndex].Date), helper.Unescaped)
				}
			} else if key == "format" {
				timeFormat = value
			}
		}
	}
	if showPublicationDate {
		return evaluateEscape(date.FormatDate(timeFormat, values.Posts[values.CurrentPostIndex].Date), helper.Unescaped)
	}
	currentDate := date.GetCurrentTime()
	return evaluateEscape(date.FormatDate(timeFormat, &currentDate), helper.Unescaped)
}

func atFirstFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	if values.CurrentHelperContext == 1 { // post
		if values.CurrentPostIndex == 0 {
			return []byte{1}
		}
		return []byte{}
	}
	if values.CurrentHelperContext == 2 { // tag
		if values.CurrentTagIndex == 0 {
			return []byte{1}
		}
		return []byte{}
	}
	return []byte{}
}

func atLastFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	if values.CurrentHelperContext == 1 { // post
		if values.CurrentPostIndex == (len(values.Posts) - 1) {
			return []byte{1}
		}
		return []byte{}
	}
	if values.CurrentHelperContext == 2 { // tag
		if values.CurrentTagIndex == (len(values.Posts[values.CurrentPostIndex].Tags) - 1) {
			return []byte{1}
		}
		return []byte{}
	}
	return []byte{}
}

func atEvenFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	if values.CurrentHelperContext == 1 { // post
		// First post (index 0) needs to be odd
		if values.CurrentPostIndex%2 == 1 {
			return []byte{1}
		}
		return []byte{}
	}
	if values.CurrentHelperContext == 2 { // tag
		// First tag (index 0) needs to be odd
		if values.CurrentTagIndex%2 == 1 {
			return []byte{1}
		}
		return []byte{}
	}
	return []byte{}
}

func atOddFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	if values.CurrentHelperContext == 1 { // post
		// First post (index 0) needs to be odd
		if values.CurrentPostIndex%2 == 0 {
			return []byte{1}
		}
		return []byte{}
	}
	if values.CurrentHelperContext == 2 { // tag
		// First tag (index 0) needs to be odd
		if values.CurrentTagIndex%2 == 0 {
			return []byte{1}
		}
		return []byte{}
	}
	return []byte{}
}

func nameFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	// If tag (commented out the code for generating a link. Ghost doesn't seem to do that either).
	if values.CurrentHelperContext == 2 { // tag
		//var buffer bytes.Buffer
		//buffer.WriteString("<a href=\"")
		//buffer.WriteString("/tag/")
		//buffer.WriteString(values.Posts[values.CurrentPostIndex].Tags[values.CurrentTagIndex].Slug)
		//buffer.WriteString("/\">")
		//buffer.Write(evaluateEscape([]byte(values.Posts[values.CurrentPostIndex].Tags[values.CurrentTagIndex].Name), helper.Unescaped))
		//buffer.WriteString("</a>")
		//return buffer.Bytes()
		return evaluateEscape(values.Posts[values.CurrentPostIndex].Tags[values.CurrentTagIndex].Name, helper.Unescaped)
	}
	// If author (commented out the code for generating a link. Ghost doesn't seem to do that).
	//var buffer bytes.Buffer
	//buffer.WriteString("<a href=\"")
	//buffer.WriteString("/author/")
	//buffer.WriteString(values.Author.Slug)
	//buffer.WriteString("\">")
	//buffer.Write(evaluateEscape([]byte(values.Author.Name), helper.Unescaped))
	//buffer.WriteString("</a>")
	//return buffer.Bytes()
	//TODO: Error handling if there is no Posts[values.CurrentPostIndex]
	return evaluateEscape(values.Posts[values.CurrentPostIndex].Author.Name, helper.Unescaped)
}

func tagDotNameFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	if len(values.CurrentTag.Name) != 0 {
		return evaluateEscape(values.CurrentTag.Name, helper.Unescaped)
	} else {
		return evaluateEscape(values.Posts[values.CurrentPostIndex].Tags[values.CurrentTagIndex].Name, helper.Unescaped)
	}
}

func tagDotSlugFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	if values.CurrentTag.Slug != "" {
		return evaluateEscape([]byte(values.CurrentTag.Slug), helper.Unescaped)
	} else {
		return evaluateEscape([]byte(values.Posts[values.CurrentPostIndex].Tags[values.CurrentTagIndex].Slug), helper.Unescaped)
	}
}

func idFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	return []byte(strconv.FormatInt(values.Posts[values.CurrentPostIndex].Id, 10))
}

func assetFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	if len(helper.Arguments) != 0 {
		var buffer bytes.Buffer
		buffer.Write(values.Blog.AssetPath)
		buffer.WriteString(helper.Arguments[0].Name)
		return buffer.Bytes()
	}
	return []byte{}
}

func foreachFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	if len(helper.Arguments) != 0 {
		switch helper.Arguments[0].Name {
		case "posts":
			var buffer bytes.Buffer
			for index, _ := range values.Posts {
				//if values.Posts[index].Id != 0 { // If post is not empty (Commented out for now. This was only neccessary in previous versions, when the array length was always the postsPerPage length)
				values.CurrentPostIndex = index
				buffer.Write(executeHelper(helper, values, 1)) // context = post
				//}
			}
			return buffer.Bytes()
		case "tags":
			var buffer bytes.Buffer
			for index, _ := range values.Posts[values.CurrentPostIndex].Tags {
				//if values.Posts[values.CurrentPostIndex].Tags[index].Id != 0 { // If tag is not empty (Commented out for now. Not neccessary.)
				values.CurrentTagIndex = index
				buffer.Write(executeHelper(helper, values, 2)) // context = tag
				//}
			}
			return buffer.Bytes()
		case "navigation":
			var buffer bytes.Buffer
			for index, _ := range values.Blog.NavigationItems {
				values.CurrentNavigationIndex = index
				buffer.Write(executeHelper(helper, values, 4)) // context = navigation
			}
			return buffer.Bytes()
		default:
			return []byte{}
		}
	}
	return []byte{}
}

func ifFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	if len(helper.Arguments) != 0 {
		if len(helper.Arguments[0].Function(&helper.Arguments[0], values)) != 0 {
			// If the evaluation is true, execute the if helper
			return executeHelper(helper, values, values.CurrentHelperContext)
		} else {
			// Else execute the else helper which is always at the last index of the if helper Arguments
			if helper.Arguments[len(helper.Arguments)-1].Name == "else" {
				if len(helper.Arguments[len(helper.Arguments)-1].Children) != 0 {
				}
				return executeHelper(&helper.Arguments[len(helper.Arguments)-1], values, values.CurrentHelperContext)
			}
		}
	}
	return []byte{}
}

func unlessFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	if len(helper.Arguments) != 0 {
		if len(helper.Arguments[0].Function(&helper.Arguments[0], values)) == 0 {
			// If the evaluation is false, execute the unless helper
			return executeHelper(helper, values, values.CurrentHelperContext)
		}
	}
	return []byte{}
}

func atBlogDotTitleFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	return evaluateEscape(values.Blog.Title, helper.Unescaped)
}

func atBlogDotUrlFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	var buffer bytes.Buffer
	buffer.Write(values.Blog.Url)
	buffer.WriteString("/")
	return evaluateEscape(buffer.Bytes(), helper.Unescaped)
}

func atBlogDotLogoFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	return evaluateEscape(values.Blog.Logo, helper.Unescaped)
}

func atBlogDotCoverFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	return evaluateEscape(values.Blog.Cover, helper.Unescaped)
}

func atBlogDotDescriptionFunc(helper *structure.Helper, values *structure.RequestData) []byte {
	return evaluateEscape(values.Blog.Description, helper.Unescaped)
}

func evaluateEscape(value []byte, unescaped bool) []byte {
	if unescaped {
		return value
	}
	return []byte(html.EscapeString(string(value)))
}

func positiveCeilingInt64(input float64) int64 {
	output := int64(input)
	if (input - float64(output)) > 0 {
		output++
	}
	return output
}
