package database

import (
	"time"
)

const stmtUpdatePost = "UPDATE posts SET title = ?, slug = ?, markdown = ?, html = ?, featured = ?, page = ?, status = ?, meta_description = ?, image = ?, updated_at = ?, updated_by = ? WHERE id = ?"
const stmtUpdatePostPublished = "UPDATE posts SET title = ?, slug = ?, markdown = ?, html = ?, featured = ?, page = ?, status = ?, meta_description = ?, image = ?, updated_at = ?, updated_by = ?, published_at = ?, published_by = ? WHERE id = ?"
const stmtUpdateSettings = "UPDATE settings SET value = ?, updated_at = ?, updated_by = ? WHERE key = ?"
const stmtUpdateUser = "UPDATE users SET name = ?, slug = ?, email = ?, image = ?, cover = ?, bio = ?, website = ?, location = ?, updated_at = ?, updated_by = ? WHERE id = ?"
const stmtUpdateLastLogin = "UPDATE users SET last_login = ? WHERE id = ?"
const stmtUpdateUserPassword = "UPDATE users SET password = ?, updated_at = ?, updated_by = ? WHERE id = ?"

// UpdatePost ...
func UpdatePost(id int64, title []byte, slug string, markdown []byte, html []byte, featured bool, isPage bool, published bool, metaDescription []byte, image []byte, updatedAt time.Time, updatedBy int64) (bool, error) {
	currentPost, err := RetrievePostByID(id)
	if err != nil {
		return false, err
	}
	status := "draft"
	if published {
		status = "published"
	}
	writeDB, err := readDB.Begin()
	if err != nil {
		writeDB.Rollback()
		return false, err
	}
	// If the updated post is published for the first time, add publication date and user
	if published && !currentPost.IsPublished {
		_, err = writeDB.Exec(stmtUpdatePostPublished, title, slug, markdown, html, featured, isPage, status, metaDescription, image, updatedAt, updatedBy, updatedAt, updatedBy, id)
	} else {
		_, err = writeDB.Exec(stmtUpdatePost, title, slug, markdown, html, featured, isPage, status, metaDescription, image, updatedAt, updatedBy, id)
	}
	if err != nil {
		writeDB.Rollback()
		return false, err
	}
	return published, writeDB.Commit()
}

// UpdateSettings ...
func UpdateSettings(title []byte, description []byte, logo []byte, cover []byte, postsPerPage int64, activeTheme string, navigation []byte, updatedAt time.Time, updatedBy int64) error {
	writeDB, err := readDB.Begin()
	if err != nil {
		writeDB.Rollback()
		return err
	}
	// Title
	_, err = writeDB.Exec(stmtUpdateSettings, title, updatedAt, updatedBy, "title")
	if err != nil {
		writeDB.Rollback()
		return err
	}
	// Description
	_, err = writeDB.Exec(stmtUpdateSettings, description, updatedAt, updatedBy, "description")
	if err != nil {
		writeDB.Rollback()
		return err
	}
	// Logo
	_, err = writeDB.Exec(stmtUpdateSettings, logo, updatedAt, updatedBy, "logo")
	if err != nil {
		writeDB.Rollback()
		return err
	}
	// Cover
	_, err = writeDB.Exec(stmtUpdateSettings, cover, updatedAt, updatedBy, "cover")
	if err != nil {
		writeDB.Rollback()
		return err
	}
	// PostsPerPage
	_, err = writeDB.Exec(stmtUpdateSettings, postsPerPage, updatedAt, updatedBy, "postsPerPage")
	if err != nil {
		writeDB.Rollback()
		return err
	}
	// ActiveTheme
	_, err = writeDB.Exec(stmtUpdateSettings, activeTheme, updatedAt, updatedBy, "activeTheme")
	if err != nil {
		writeDB.Rollback()
		return err
	}
	// Navigation
	_, err = writeDB.Exec(stmtUpdateSettings, navigation, updatedAt, updatedBy, "navigation")
	if err != nil {
		writeDB.Rollback()
		return err
	}
	return writeDB.Commit()
}

// UpdateActiveTheme ...
func UpdateActiveTheme(activeTheme string, updatedAt time.Time, updatedBy int64) error {
	writeDB, err := readDB.Begin()
	if err != nil {
		writeDB.Rollback()
		return err
	}
	_, err = writeDB.Exec(stmtUpdateSettings, activeTheme, updatedAt, updatedBy, "activeTheme")
	if err != nil {
		writeDB.Rollback()
		return err
	}
	return writeDB.Commit()
}

// UpdateUser ...
func UpdateUser(id int64, name []byte, slug string, email []byte, image []byte, cover []byte, bio []byte, website []byte, location []byte, updatedAt time.Time, updatedBy int64) error {
	writeDB, err := readDB.Begin()
	if err != nil {
		writeDB.Rollback()
		return err
	}
	_, err = writeDB.Exec(stmtUpdateUser, name, slug, email, image, cover, bio, website, location, updatedAt, updatedBy, id)
	if err != nil {
		writeDB.Rollback()
		return err
	}
	return writeDB.Commit()
}

// UpdateLastLogin ...
func UpdateLastLogin(logInDate time.Time, userID int64) error {
	writeDB, err := readDB.Begin()
	if err != nil {
		writeDB.Rollback()
		return err
	}
	_, err = writeDB.Exec(stmtUpdateLastLogin, logInDate, userID)
	if err != nil {
		writeDB.Rollback()
		return err
	}
	return writeDB.Commit()
}

// UpdateUserPassword ...
func UpdateUserPassword(id int64, password string, updatedAt time.Time, updatedBy int64) error {
	writeDB, err := readDB.Begin()
	if err != nil {
		writeDB.Rollback()
		return err
	}
	_, err = writeDB.Exec(stmtUpdateUserPassword, password, updatedAt, updatedBy, id)
	if err != nil {
		writeDB.Rollback()
		return err
	}
	return writeDB.Commit()
}
