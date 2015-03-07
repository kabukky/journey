package methods

import (
	"github.com/kabukky/journey/database"
	"github.com/kabukky/journey/structure"
	"time"
)

func SaveAuthor(a *structure.Author, hashedPassword string, createdBy int64) error {
	err := database.InsertUser(a.Name, a.Slug, hashedPassword, a.Email, a.Image, a.Cover, time.Now(), createdBy)
	if err != nil {
		return err
	}
	return nil
}

func UpdateAuthor(a *structure.Author, updatedById int64) error {
	err := database.UpdateUser(a.Id, a.Email, a.Image, a.Cover, a.Bio, a.Website, a.Location, time.Now(), updatedById)
	if err != nil {
		return err
	}
	return nil
}
