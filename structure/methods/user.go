package methods

import (
	"github.com/kabukky/journey/database"
	"github.com/kabukky/journey/date"
	"github.com/kabukky/journey/structure"
)

func SaveUser(u *structure.User, hashedPassword string, createdBy int64) error {
	userId, err := database.InsertUser(u.Name, u.Slug, hashedPassword, u.Email, u.Image, u.Cover, date.GetCurrentTime(), createdBy)
	if err != nil {
		return err
	}
	err = database.InsertRoleUser(u.Role, userId)
	if err != nil {
		return err
	}
	return nil
}

func UpdateUser(u *structure.User, updatedById int64) error {
	err := database.UpdateUser(u.Id, u.Name, u.Slug, u.Email, u.Image, u.Cover, u.Bio, u.Website, u.Location, date.GetCurrentTime(), updatedById)
	if err != nil {
		return err
	}
	return nil
}
