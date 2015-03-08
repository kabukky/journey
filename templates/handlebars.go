package templates

import (
	"bytes"
	"github.com/kabukky/journey/conversion"
	"github.com/kabukky/journey/database"
	"github.com/kabukky/journey/filenames"
	"github.com/kabukky/journey/structure"
	"html"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Ghost always includes a link to jquery in it's footer func. Mimic this.
var jqueryCodeForFooter = []byte("<script src=\"" + filenames.JqueryFilename + "\"></script>")

// Helper fuctions
func nullFunc(helper *Helper, values *structure.RequestData) []byte {
	log.Println("Warning: This helper is not implemented:", helper.Name)
	return []byte{}
}

func paginationDotTotalFunc(helper *Helper, values *structure.RequestData) []byte {
	if values.CurrentTemplate == 0 { // index
		return []byte(strconv.FormatInt(values.Blog.PostCount, 10))
	} else if values.CurrentTemplate == 3 { // author
		count, err := database.RetrieveNumberOfPostsByAuthor(values.Posts[values.CurrentPostIndex].Author.Id)
		if err != nil {
			log.Println("couldn't get number of posts", err.Error())
			return []byte{}
		}
		return []byte(strconv.FormatInt(count, 10))
	} else if values.CurrentTemplate == 2 { // tag
		count, err := database.RetrieveNumberOfPostsByTag(values.CurrentTag.Id)
		if err != nil {
			log.Println("couldn't get number of posts", err.Error())
			return []byte{}
		}
		return []byte(strconv.FormatInt(count, 10))
	}
	return []byte{}
}

func pluralFunc(helper *Helper, values *structure.RequestData) []byte {
	countString := string(helper.Arguments[0].Function(helper, values))
	if countString == "" {
		log.Println("Couldn't get count in plural helper")
		return []byte{}
	}
	for _, argument := range helper.Arguments[1:] {
		if countString == "0" && strings.HasPrefix(argument.Name, "empty") {
			output := argument.Name[len("empty"):]
			output = strings.Replace(output, "%", countString, -1)
			return []byte(output)
		} else if countString == "1" && strings.HasPrefix(argument.Name, "singular") {
			output := argument.Name[len("singular"):]
			output = strings.Replace(output, "%", countString, -1)
			return []byte(output)
		} else if countString != "0" && countString != "1" && strings.HasPrefix(argument.Name, "plural") {
			output := argument.Name[len("plural"):]
			output = strings.Replace(output, "%", countString, -1)
			return []byte(output)
		}
	}
	return []byte{}
}

func prevFunc(helper *Helper, values *structure.RequestData) []byte {
	if values.CurrentIndexPage > 1 {
		return []byte{1}
	}
	return []byte{}
}

func nextFunc(helper *Helper, values *structure.RequestData) []byte {
	var count int64
	var err error
	if values.CurrentTemplate == 0 { // index
		count = values.Blog.PostCount
	} else if values.CurrentTemplate == 2 { // tag
		count, err = database.RetrieveNumberOfPostsByTag(values.CurrentTag.Id)
		if err != nil {
			log.Println("couldn't get number of posts for tag", err.Error())
			return []byte{}
		}
	} else if values.CurrentTemplate == 3 { // author
		count, err = database.RetrieveNumberOfPostsByAuthor(values.Posts[values.CurrentPostIndex].Author.Id)
		if err != nil {
			log.Println("couldn't get number of posts for author", err.Error())
			return []byte{}
		}
	}
	maxPages := int64((float64(count) / float64(values.Blog.PostsPerPage)) + 0.5)
	if int64(values.CurrentIndexPage) < maxPages {
		return []byte{1}
	}
	return []byte{}
}

func pageFunc(helper *Helper, values *structure.RequestData) []byte {
	return []byte(strconv.Itoa(values.CurrentIndexPage))
}

func pagesFunc(helper *Helper, values *structure.RequestData) []byte {
	var count int64
	var err error
	if values.CurrentTemplate == 0 { // index
		count = values.Blog.PostCount	
	} else if values.CurrentTemplate == 2 { // tag
		count, err = database.RetrieveNumberOfPostsByTag(values.CurrentTag.Id)
		if err != nil {
			log.Println("couldn't get number of posts for tag", err.Error())
			return []byte{}
		}
	} else if values.CurrentTemplate == 3 { // author
		count, err = database.RetrieveNumberOfPostsByAuthor(values.Posts[values.CurrentPostIndex].Author.Id)
		if err != nil {
			log.Println("couldn't get number of posts for author", err.Error())
			return []byte{}
		}
	}
	maxPages := int64((float64(count) / float64(values.Blog.PostsPerPage)) + 0.5)
	return []byte(strconv.FormatInt(maxPages, 10))
}

func page_urlFunc(helper *Helper, values *structure.RequestData) []byte {
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
					buffer.WriteString("/page/")
					buffer.WriteString(strconv.Itoa(values.CurrentIndexPage - 1))
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
					log.Println("couldn't get number of posts for tag", err.Error())
					return []byte{}
				}
			} else if values.CurrentTemplate == 3 { // author
				count, err = database.RetrieveNumberOfPostsByAuthor(values.Posts[values.CurrentPostIndex].Author.Id)
				if err != nil {
					log.Println("couldn't get number of posts for author", err.Error())
					return []byte{}
				}
			}
			maxPages := int64((float64(count) / float64(values.Blog.PostsPerPage)) + 0.5)
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
				buffer.WriteString("/page/")
				buffer.WriteString(strconv.Itoa(values.CurrentIndexPage + 1))
				buffer.WriteString("/")
				return buffer.Bytes()
			}
		}
	}
	return []byte{}
}

func extendFunc(helper *Helper, values *structure.RequestData) []byte {
	if len(helper.Arguments) != 0 {
		return []byte(helper.Arguments[0].Name)
	}
	return []byte{}
}

func featuredFunc(helper *Helper, values *structure.RequestData) []byte {
	if values.Posts[values.CurrentPostIndex].IsFeatured {
		return []byte{1}
	}
	return []byte{}
}

func body_classFunc(helper *Helper, values *structure.RequestData) []byte {
	if values.CurrentTemplate == 1 { // post
		// TODO: is there anything else that needs to get output here?
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

func ghost_headFunc(helper *Helper, values *structure.RequestData) []byte {
	// TODO: Implement
	return []byte{}
}

func ghost_footFunc(helper *Helper, values *structure.RequestData) []byte {
	// TODO: This seems to just output a jquery link in ghost. Keep for compatibility?
	return jqueryCodeForFooter
}

func meta_titleFunc(helper *Helper, values *structure.RequestData) []byte {
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

func meta_descriptionFunc(helper *Helper, values *structure.RequestData) []byte {
	// TODO: Finish this
	if values.CurrentTemplate != 1 { // not post
		return evaluateEscape(values.Blog.Description, helper.Unescaped)
	}
	// Nothing on post yet
	return []byte{}
}

func bodyFunc(helper *Helper, values *structure.RequestData) []byte {
	return helper.Block
}

func insertFunc(helper *Helper, values *structure.RequestData) []byte {
	if len(helper.Arguments) != 0 {
		if templateHelper, ok := compiledTemplates.m[helper.Arguments[0].Name]; ok {
			return executeHelper(templateHelper, values, values.CurrentHelperContext)
		}
	}
	return []byte{}
}

func encodeFunc(helper *Helper, values *structure.RequestData) []byte {
	if len(helper.Arguments) != 0 {
		return []byte(url.QueryEscape(string(helper.Arguments[0].Function(&helper.Arguments[0], values))))
	}
	return []byte{}
}

func authorFunc(helper *Helper, values *structure.RequestData) []byte {
	// Check if helper is block helper
	if len(helper.Block) != 0 {
		return executeHelper(helper, values, 3) // context = author
	}
	// Else return author.name
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

func authorDotNameFunc(helper *Helper, values *structure.RequestData) []byte {
	var buffer bytes.Buffer
	buffer.WriteString("<a href=\"")
	buffer.WriteString("/author/")
	// TODO: Error handling if there is no Posts[values.CurrentPostIndex]
	buffer.WriteString(values.Posts[values.CurrentPostIndex].Author.Slug)
	buffer.WriteString("\">")
	// TODO: Error handling if there is no Posts[values.CurrentPostIndex]
	buffer.Write(evaluateEscape(values.Posts[values.CurrentPostIndex].Author.Name, helper.Unescaped))
	buffer.WriteString("</a>")
	return buffer.Bytes()
}

func bioFunc(helper *Helper, values *structure.RequestData) []byte {
	// TODO: Error handling if there is no Posts[values.CurrentPostIndex]
	return evaluateEscape(values.Posts[values.CurrentPostIndex].Author.Bio, helper.Unescaped)
}

func emailFunc(helper *Helper, values *structure.RequestData) []byte {
	// TODO: Error handling if there is no Posts[values.CurrentPostIndex]
	return evaluateEscape(values.Posts[values.CurrentPostIndex].Author.Email, helper.Unescaped)
}

func websiteFunc(helper *Helper, values *structure.RequestData) []byte {
	// TODO: Error handling if there is no Posts[values.CurrentPostIndex]
	return evaluateEscape(values.Posts[values.CurrentPostIndex].Author.Website, helper.Unescaped)
}

func imageFunc(helper *Helper, values *structure.RequestData) []byte {
	if values.CurrentHelperContext == 1 { // post
		// TODO: Error handling if there is no Posts[values.CurrentPostIndex]
		return evaluateEscape(values.Posts[values.CurrentPostIndex].Image, helper.Unescaped)
	} else if values.CurrentHelperContext == 3 { // author
		// TODO: Error handling if there is no Posts[values.CurrentPostIndex]
		return evaluateEscape(values.Posts[values.CurrentPostIndex].Author.Image, helper.Unescaped)
	}
	return []byte{}
}

func authorDotImageFunc(helper *Helper, values *structure.RequestData) []byte {
	// TODO: Error handling if there is no Posts[values.CurrentPostIndex]
	return evaluateEscape(values.Posts[values.CurrentPostIndex].Author.Image, helper.Unescaped)
}

func coverFunc(helper *Helper, values *structure.RequestData) []byte {
	// TODO: Error handling if there is no Posts[values.CurrentPostIndex]
	return evaluateEscape(values.Posts[values.CurrentPostIndex].Author.Cover, helper.Unescaped)
}

func locationFunc(helper *Helper, values *structure.RequestData) []byte {
	// TODO: Error handling if there is no Posts[values.CurrentPostIndex]
	return evaluateEscape(values.Posts[values.CurrentPostIndex].Author.Location, helper.Unescaped)
}

func postFunc(helper *Helper, values *structure.RequestData) []byte {
	return executeHelper(helper, values, 1) // context = post
}

func postsFunc(helper *Helper, values *structure.RequestData) []byte {
	if len(values.Posts) > 0 {
		return []byte{1}
	}
	return []byte{}
}

func tagsFunc(helper *Helper, values *structure.RequestData) []byte {
	if len(values.Posts[values.CurrentPostIndex].Tags) > 0 {
		separator := ", "
		suffix := ""
		prefix := ""
		makeLink := true
		if len(helper.Arguments) != 0 {
			for _, argument := range helper.Arguments {
				if strings.HasPrefix(argument.Name, "separator") {
					separator = argument.Name[len("separator"):]

				} else if strings.HasPrefix(argument.Name, "suffix") {
					suffix = argument.Name[len("suffix"):]
				} else if strings.HasPrefix(argument.Name, "prefix") {
					prefix = argument.Name[len("prefix"):]
				} else if strings.HasPrefix(argument.Name, "autolink") {
					if argument.Name[len("autolink"):] == "false" {
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

func post_classFunc(helper *Helper, values *structure.RequestData) []byte {
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

func urlFunc(helper *Helper, values *structure.RequestData) []byte {
	if values.CurrentHelperContext == 1 { // post
		var buffer bytes.Buffer
		buffer.WriteString("/")
		buffer.WriteString(values.Posts[values.CurrentPostIndex].Slug)
		buffer.WriteString("/")
		return evaluateEscape(buffer.Bytes(), helper.Unescaped)
	} else if values.CurrentHelperContext == 3 { // author
		var buffer bytes.Buffer
		buffer.WriteString("/author/")
		// TODO: Error handling if there is no Posts[values.CurrentPostIndex]
		buffer.WriteString(values.Posts[values.CurrentPostIndex].Author.Slug)
		buffer.WriteString("/")
		return evaluateEscape(buffer.Bytes(), helper.Unescaped)
	}
	return []byte{}
}

func titleFunc(helper *Helper, values *structure.RequestData) []byte {
	return evaluateEscape(values.Posts[values.CurrentPostIndex].Title, helper.Unescaped)
}

func contentFunc(helper *Helper, values *structure.RequestData) []byte {
	// TODO: is content always unescaped? seems like it...
	return values.Posts[values.CurrentPostIndex].Html
}

func excerptFunc(helper *Helper, values *structure.RequestData) []byte {
	if values.CurrentHelperContext == 1 { // post
		if len(helper.Arguments) != 0 {
			if strings.HasPrefix(helper.Arguments[0].Name, "words") {
				number, err := strconv.Atoi(helper.Arguments[0].Name[len("words"):])
				if err == nil {
					excerpt := conversion.StripTagsFromHtml(values.Posts[values.CurrentPostIndex].Html)
					words := bytes.Fields(excerpt)
					if len(words) < number {
						return excerpt
					}
					return bytes.Join(words[:number], []byte(" "))
				}
			} else if strings.HasPrefix(helper.Arguments[0].Name, "characters") {
				number, err := strconv.Atoi(helper.Arguments[0].Name[len("characters"):])
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

func dateFunc(helper *Helper, values *structure.RequestData) []byte {
	showPublicationDate := false
	timeFormat := "MMM Do, YYYY" // Default time format
	// If in scope of a post, change default to published date
	if values.CurrentHelperContext == 1 { // post
		showPublicationDate = true
	}
	// Get the date
	if len(helper.Arguments) != 0 {
		for index, _ := range helper.Arguments {
			if helper.Arguments[index].Name == "published_at" {
				showPublicationDate = true
			} else if strings.HasPrefix(helper.Arguments[index].Name, "timeago") {
				if helper.Arguments[index].Name[len("timeago"):] == "true" {
					// Compute time ago
					return evaluateEscape(generateTimeAgo(&values.Posts[values.CurrentPostIndex].Date), helper.Unescaped)
				}
			} else if strings.HasPrefix(helper.Arguments[index].Name, "format") {
				timeFormat = helper.Arguments[index].Name[len("format"):]
			}
		}
	}
	if showPublicationDate {
		return evaluateEscape(formatDate(timeFormat, &values.Posts[values.CurrentPostIndex].Date), helper.Unescaped)
	}
	date := time.Now()
	return evaluateEscape(formatDate(timeFormat, &date), helper.Unescaped)
}

func atFirstFunc(helper *Helper, values *structure.RequestData) []byte {
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

func atLastFunc(helper *Helper, values *structure.RequestData) []byte {
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

func atEvenFunc(helper *Helper, values *structure.RequestData) []byte {
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

func atOddFunc(helper *Helper, values *structure.RequestData) []byte {
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

func nameFunc(helper *Helper, values *structure.RequestData) []byte {
	// If tag (commented out the code for generating a link. Ghost doesn't seem to do that.
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
	// If author (commented out the code for generating a link. Ghost doesn't seem to do that.
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

func tagDotNameFunc(helper *Helper, values *structure.RequestData) []byte {
	if len(values.CurrentTag.Name) != 0 {
		return evaluateEscape(values.CurrentTag.Name, helper.Unescaped)
	} else {
		return evaluateEscape(values.Posts[values.CurrentPostIndex].Tags[values.CurrentTagIndex].Name, helper.Unescaped)
	}
}

func tagDotSlugFunc(helper *Helper, values *structure.RequestData) []byte {
	if values.CurrentTag.Slug != "" {
		return evaluateEscape([]byte(values.CurrentTag.Slug), helper.Unescaped)
	} else {
		return evaluateEscape([]byte(values.Posts[values.CurrentPostIndex].Tags[values.CurrentTagIndex].Slug), helper.Unescaped)
	}
}

func paginationFunc(helper *Helper, values *structure.RequestData) []byte {
	if template, ok := compiledTemplates.m["pagination"]; ok { // If the theme has a pagination.hbs
		return executeHelper(template, values, values.CurrentHelperContext)
	}
	var count int64
	var err error
	if values.CurrentTemplate == 0 { // index
		count = values.Blog.PostCount
	} else if values.CurrentTemplate == 2 { // tag
		count, err = database.RetrieveNumberOfPostsByTag(values.CurrentTag.Id)
		if err != nil {
			log.Println("couldn't get number of posts for tag", err.Error())
			return []byte{}
		}
	} else if values.CurrentTemplate == 3 { // author
		count, err = database.RetrieveNumberOfPostsByAuthor(values.Posts[values.CurrentPostIndex].Author.Id)
		if err != nil {
			log.Println("couldn't get number of posts for author", err.Error())
			return []byte{}
		}
	}
	if count > values.Blog.PostsPerPage {
		maxPages := int64((float64(count) / float64(values.Blog.PostsPerPage)) + 0.5)
		var buffer bytes.Buffer
		buffer.WriteString("<nav class=\"pagination\" role=\"navigation\">")
		// If this is not the first index page, display a back link
		if values.CurrentIndexPage > 1 {
			buffer.WriteString("\n\t\t<a class=\"newer-posts\" href=\"")
			if values.CurrentIndexPage == 2 {
				if values.CurrentTemplate == 3 { // author
					buffer.WriteString("/author/")
					// TODO: Error handling if there is no Posts[values.CurrentPostIndex]
					buffer.WriteString(values.Posts[values.CurrentPostIndex].Author.Slug)
				} else if values.CurrentTemplate == 2 { // tag
					buffer.WriteString("/tag/")
					// TODO: Error handling if there is no Posts[values.CurrentPostIndex]
					buffer.WriteString(values.CurrentTag.Slug)
				}
				buffer.WriteString("/")
			} else {
				if values.CurrentTemplate == 3 { // author
					buffer.WriteString("/author/")
					// TODO: Error handling if there is no Posts[values.CurrentPostIndex]
					buffer.WriteString(values.Posts[values.CurrentPostIndex].Author.Slug)
				} else if values.CurrentTemplate == 2 { // tag
					buffer.WriteString("/tag/")
					// TODO: Error handling if there is no Posts[values.CurrentPostIndex]
					buffer.WriteString(values.CurrentTag.Slug)
				}
				buffer.WriteString("/page/")
				buffer.WriteString(strconv.Itoa(values.CurrentIndexPage - 1))
				buffer.WriteString("/")
			}
			buffer.WriteString("\">&larr; Newer Posts</a>")
		}
		buffer.WriteString("\n\t<span class=\"page-number\">Page ")
		buffer.WriteString(strconv.Itoa(values.CurrentIndexPage))
		buffer.WriteString(" of ")
		buffer.WriteString(strconv.FormatInt(maxPages, 10))
		buffer.WriteString("</span>")
		if int64(values.CurrentIndexPage) < maxPages {
			buffer.WriteString("\n\t\t<a class=\"older-posts\" href=\"")
			if values.CurrentTemplate == 3 { // author
				buffer.WriteString("/author/")
				// TODO: Error handling if there is no Posts[values.CurrentPostIndex]
				buffer.WriteString(values.Posts[values.CurrentPostIndex].Author.Slug)
			} else if values.CurrentTemplate == 2 { // tag
				buffer.WriteString("/tag/")
				// TODO: Error handling if there is no Posts[values.CurrentPostIndex]
				buffer.WriteString(values.CurrentTag.Slug)
			}
			buffer.WriteString("/page/")
			buffer.WriteString(strconv.Itoa(values.CurrentIndexPage + 1))
			buffer.WriteString("/\">Older Posts &rarr;</a>")
		}
		buffer.WriteString("\n</nav>")
		return buffer.Bytes()
	} else {
		return []byte("<nav class=\"pagination\" role=\"navigation\">\n\t<span class=\"page-number\">Page 1 of 1</span>\n</nav>")
	}
}

func idFunc(helper *Helper, values *structure.RequestData) []byte {
	return []byte(strconv.FormatInt(values.Posts[values.CurrentPostIndex].Id, 10))
}

func assetFunc(helper *Helper, values *structure.RequestData) []byte {
	if len(helper.Arguments) != 0 {
		var buffer bytes.Buffer
		buffer.Write(values.Blog.AssetPath)
		buffer.WriteString(helper.Arguments[0].Name)
		return buffer.Bytes()
	}
	return []byte{}
}

func foreachFunc(helper *Helper, values *structure.RequestData) []byte {
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
		default:
			return []byte{}
		}
	}
	return []byte{}
}

func ifFunc(helper *Helper, values *structure.RequestData) []byte {
	if len(helper.Arguments) != 0 {
		if len(helper.Arguments[0].Function(&helper.Arguments[0], values)) != 0 {
			return executeHelper(helper, values, values.CurrentHelperContext)
		} else {
			// Execute the else helper which is always at the last index of the if helper Arguments
			if helper.Arguments[len(helper.Arguments)-1].Name == "else" {
				if len(helper.Arguments[len(helper.Arguments)-1].Children) != 0 {
				}
				return executeHelper(&helper.Arguments[len(helper.Arguments)-1], values, values.CurrentHelperContext)
			}
		}
	}
	return []byte{}
}

func unlessFunc(helper *Helper, values *structure.RequestData) []byte {
	if len(helper.Arguments) != 0 {
		if len(helper.Arguments[0].Function(&helper.Arguments[0], values)) == 0 {
			return executeHelper(helper, values, values.CurrentHelperContext)
		} else {
			// Execute the else helper which is always at the last index of the if helper Arguments
			if helper.Arguments[len(helper.Arguments)-1].Name == "else" {
				if len(helper.Arguments[len(helper.Arguments)-1].Children) != 0 {
				}
				return executeHelper(&helper.Arguments[len(helper.Arguments)-1], values, values.CurrentHelperContext)
			}
		}
	}
	return []byte{}
}

func atBlogDotTitleFunc(helper *Helper, values *structure.RequestData) []byte {
	return evaluateEscape(values.Blog.Title, helper.Unescaped)
}

func atBlogDotUrlFunc(helper *Helper, values *structure.RequestData) []byte {
	var buffer bytes.Buffer
	// Write // in front of url to be protocol agnostic
	buffer.WriteString("//")
	buffer.Write(values.Blog.Url)
	return evaluateEscape(buffer.Bytes(), helper.Unescaped)
}

func atBlogDotLogoFunc(helper *Helper, values *structure.RequestData) []byte {
	return evaluateEscape(values.Blog.Logo, helper.Unescaped)
}

func atBlogDotCoverFunc(helper *Helper, values *structure.RequestData) []byte {
	return evaluateEscape(values.Blog.Cover, helper.Unescaped)
}

func atBlogDotDescriptionFunc(helper *Helper, values *structure.RequestData) []byte {
	return evaluateEscape(values.Blog.Description, helper.Unescaped)
}

func evaluateEscape(value []byte, unescaped bool) []byte {
	if unescaped {
		return value
	}
	return []byte(html.EscapeString(string(value)))
}
