package structure

import (
	"time"
)

type Post struct {
	Id          int64
	Uuid        []byte
	Title       []byte
	Slug        string
	Markdown    []byte
	Html        []byte
	IsFeatured  bool
	IsPage      bool
	IsPublished bool
	Date        time.Time
	Tags        []Tag
	Author      *User
	Image       []byte
}
