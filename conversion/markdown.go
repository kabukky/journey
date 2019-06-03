package conversion

import (
	"github.com/russross/blackfriday"
)

// Using blackfriday library for markdown to html conversion. At least for now.

const (
	htmlFlags = 0 |
		// We don't want blackfriday.HTML_USE_XHTML
		blackfriday.Smartypants |
		blackfriday.SmartypantsFractions |
		blackfriday.SmartypantsLatexDashes |
		blackfriday.FootnoteReturnLinks

	extensions = 0 |
		blackfriday.Tables |
		blackfriday.FencedCode |
		blackfriday.Autolink |
		blackfriday.Strikethrough |
		blackfriday.SpaceHeadings |
		blackfriday.HeadingIDs |
		blackfriday.BackslashLineBreak |
		blackfriday.Footnotes
)

func GenerateHtmlFromMarkdown(input []byte) []byte {
	renderParameters := blackfriday.HTMLRendererParameters{}
	renderParameters.Flags = htmlFlags
	renderer := blackfriday.NewHTMLRenderer(renderParameters)
	return blackfriday.Run(input, blackfriday.WithRenderer(renderer), blackfriday.WithExtensions(extensions))
}
