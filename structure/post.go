package structure

import (
	"time"
)

// Post ...
type Post struct {
	ID              int64
	UUID            []byte
	Title           []byte
	Slug            string
	Markdown        []byte
	HTML            []byte
	IsFeatured      bool
	IsPage          bool
	IsPublished     bool
	Date            *time.Time
	Tags            []Tag
	Author          *User
	MetaDescription []byte
	Image           []byte
}
