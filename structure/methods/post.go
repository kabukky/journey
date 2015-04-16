package methods

import (
	"github.com/kabukky/journey/database"
	"github.com/kabukky/journey/structure"
	"time"
)

func SavePost(p *structure.Post) error {
	tagIds := make([]int64, 0)
	// Insert tags
	for _, tag := range p.Tags {
		// Tag slug might already be in database
		tagId, err := database.RetrieveTagIdBySlug(tag.Slug)
		if err != nil {
			// Tag is probably not in database yet
			tagId, err = database.InsertTag(tag.Name, tag.Slug, time.Now(), p.Author.Id)
			if err != nil {
				return err
			}
		}
		if tagId != 0 {
			tagIds = append(tagIds, tagId)
		}
	}
	// Insert post
	postId, err := database.InsertPost(p.Title, p.Slug, p.Markdown, p.Html, p.IsFeatured, p.IsPage, p.IsPublished, p.Image, *p.Date, p.Author.Id)
	if err != nil {
		return err
	}
	// Insert postTags
	for _, tagId := range tagIds {
		err = database.InsertPostTag(postId, tagId)
		if err != nil {
			return err
		}
	}
	return nil
}

func UpdatePost(p *structure.Post) error {
	tagIds := make([]int64, 0)
	// Insert tags
	for _, tag := range p.Tags {
		// Tag slug might already be in database
		tagId, err := database.RetrieveTagIdBySlug(tag.Slug)
		if err != nil {
			// Tag is probably not in database yet
			tagId, err = database.InsertTag(tag.Name, tag.Slug, time.Now(), p.Author.Id)
			if err != nil {
				return err
			}
		}
		if tagId != 0 {
			tagIds = append(tagIds, tagId)
		}
	}
	// Update post
	err := database.UpdatePost(p.Id, p.Title, p.Slug, p.Markdown, p.Html, p.IsFeatured, p.IsPage, p.IsPublished, p.Image, *p.Date, p.Author.Id)
	if err != nil {
		return err
	}
	// Delete old postTags
	err = database.DeletePostTagsForPostId(p.Id)
	// Insert postTags
	if err != nil {
		return err
	}
	for _, tagId := range tagIds {
		err = database.InsertPostTag(p.Id, tagId)
		if err != nil {
			return err
		}
	}
	return nil
}
