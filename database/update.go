package database

import (
	"time"
)

const stmtUpdatePost = "UPDATE posts SET title = ?, slug = ?, markdown = ?, html = ?, featured = ?, page = ?, status = ?, image = ?, updated_at = ?, updated_by = ? WHERE id = ?"
const stmtUpdateSettings = "UPDATE settings SET value = ? WHERE key = ?"
const stmtUpdateUser = "UPDATE users SET email = ?, image = ?, cover = ?, bio = ?, website = ?, location = ?, updated_at = ?, updated_by = ? WHERE id = ?"
const stmtUpdateUserPassword = "UPDATE users SET password = ?, updated_at = ?, updated_by = ? WHERE id = ?"

func UpdatePost(id int64, title []byte, slug string, markdown []byte, html []byte, featured bool, isPage bool, published bool, image []byte, updated_at time.Time, updated_by int64) error {
	status := "draft"
	if published {
		status = "published"
	}
	writeDB, err := readDB.Begin()
	if err != nil {
		writeDB.Rollback()
		return err
	}
	_, err = writeDB.Exec(stmtUpdatePost, title, slug, markdown, html, featured, isPage, status, image, updated_at, updated_by, id)
	if err != nil {
		writeDB.Rollback()
		return err
	}
	return writeDB.Commit()
}

func UpdateSettings(title []byte, description []byte, logo []byte, cover []byte, postsPerPage int64, activeTheme string) error {
	writeDB, err := readDB.Begin()
	if err != nil {
		writeDB.Rollback()
		return err
	}
	// Title
	_, err = writeDB.Exec(stmtUpdateSettings, title, "title")
	if err != nil {
		writeDB.Rollback()
		return err
	}
	// Description
	_, err = writeDB.Exec(stmtUpdateSettings, description, "description")
	if err != nil {
		writeDB.Rollback()
		return err
	}
	// Logo
	_, err = writeDB.Exec(stmtUpdateSettings, logo, "logo")
	if err != nil {
		writeDB.Rollback()
		return err
	}
	// Cover
	_, err = writeDB.Exec(stmtUpdateSettings, cover, "cover")
	if err != nil {
		writeDB.Rollback()
		return err
	}
	// PostsPerPage
	_, err = writeDB.Exec(stmtUpdateSettings, postsPerPage, "postsPerPage")
	if err != nil {
		writeDB.Rollback()
		return err
	}
	// ActiveTheme
	_, err = writeDB.Exec(stmtUpdateSettings, activeTheme, "activeTheme")
	if err != nil {
		writeDB.Rollback()
		return err
	}
	return writeDB.Commit()
}

func UpdateUser(id int64, email []byte, image []byte, cover []byte, bio []byte, website []byte, location []byte, updated_at time.Time, updated_by int64) error {
	writeDB, err := readDB.Begin()
	if err != nil {
		writeDB.Rollback()
		return err
	}
	_, err = writeDB.Exec(stmtUpdateUser, email, image, cover, bio, website, location, updated_at, updated_by, id)
	if err != nil {
		writeDB.Rollback()
		return err
	}
	return writeDB.Commit()
}

func UpdateUserPassword(id int64, password string, updated_at time.Time, updated_by int64) error {
	writeDB, err := readDB.Begin()
	if err != nil {
		writeDB.Rollback()
		return err
	}
	_, err = writeDB.Exec(stmtUpdateUserPassword, password, updated_at, updated_by, id)
	if err != nil {
		writeDB.Rollback()
		return err
	}
	return writeDB.Commit()
}
