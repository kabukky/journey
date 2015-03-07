package conversion

import (
	"github.com/russross/blackfriday"
)

// Using blackfriday library for markdown to html conversion. At least for now.

func GenerateHtmlFromMarkdown(input []byte) []byte {
	return blackfriday.MarkdownCommon(input)
}
