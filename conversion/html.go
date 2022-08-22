package conversion

import (
	"bytes"
	"regexp"
)

var tagChecker = regexp.MustCompile("<.*?>")
var whitespaceChecker = regexp.MustCompile("\\s{2,}")

// StripTagsFromHTML removes tags from HTML, and also changes
// whitespace to spaces
func StripTagsFromHTML(input []byte) []byte {
	output := tagChecker.ReplaceAll(input, []byte{})
	output = bytes.Replace(output, []byte("\n"), []byte(" "), -1)
	output = bytes.Replace(output, []byte("\t"), []byte(" "), -1)
	output = whitespaceChecker.ReplaceAll(output, []byte(" "))
	return output
}
