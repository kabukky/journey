package methods

import (
	"github.com/rkuris/journey/database"
	"github.com/rkuris/journey/date"
	"github.com/rkuris/journey/structure"
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

func UpdateUser(u *structure.User, updatedByID int64) error {
	err := database.UpdateUser(u.ID, u.Name, u.Slug, u.Email, u.Image, u.Cover, u.Bio, u.Website, u.Location, date.GetCurrentTime(), updatedByID)
	if err != nil {
		return err
	}
	return nil
}
