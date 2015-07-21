package conversion

import (
	"github.com/russross/blackfriday"
)

// Using blackfriday library for markdown to html conversion. At least for now.

const (
	htmlFlags = 0 |
		// We don't want blackfriday.HTML_USE_XHTML
		blackfriday.HTML_USE_SMARTYPANTS |
		blackfriday.HTML_SMARTYPANTS_FRACTIONS |
		blackfriday.HTML_SMARTYPANTS_LATEX_DASHES |
		blackfriday.HTML_FOOTNOTE_RETURN_LINKS

	extensions = 0 |
		blackfriday.EXTENSION_TABLES |
		blackfriday.EXTENSION_FENCED_CODE |
		blackfriday.EXTENSION_AUTOLINK |
		blackfriday.EXTENSION_STRIKETHROUGH |
		blackfriday.EXTENSION_SPACE_HEADERS |
		blackfriday.EXTENSION_HEADER_IDS |
		blackfriday.EXTENSION_BACKSLASH_LINE_BREAK |
		blackfriday.EXTENSION_FOOTNOTES
)

func GenerateHtmlFromMarkdown(input []byte) []byte {
	renderer := blackfriday.HtmlRenderer(htmlFlags, "", "")
	return blackfriday.Markdown(input, renderer, extensions)
}
