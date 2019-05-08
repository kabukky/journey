package methods

import (
	"strings"

	"github.com/kabukky/journey/slug"
	"github.com/kabukky/journey/structure"
)

func GenerateTagsFromCommaString(input string) []structure.Tag {
	output := make([]structure.Tag, 0)
	tags := strings.Split(input, ",")
	for index, _ := range tags {
		tags[index] = strings.TrimSpace(tags[index])
	}
	for _, tag := range tags {
		if tag != "" {
			output = append(output, structure.Tag{Name: []byte(tag), Slug: slug.Generate(tag, "tags")})
		}
	}
	return output
}
