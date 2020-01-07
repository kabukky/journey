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

func UpdatePost(id int64, title []byte, slug string, markdown []byte, html []byte, featured bool, isPage bool, published bool, meta_description []byte, image []byte, updated_at time.Time, updated_by int64) error {
	currentPost, err := RetrievePostById(id)
	if err != nil {
		return err
	}
	status := "draft"
	if published {
		status = "published"
	}
	writeDB, err := readDB.Begin()
	if err != nil {
		_ = writeDB.Rollback()
		return err
	}
	// If the updated post is published for the first time, add publication date and user
	if published && !currentPost.IsPublished {
		_, err = writeDB.Exec(stmtUpdatePostPublished, title, slug, markdown, html, featured, isPage, status, meta_description, image, updated_at, updated_by, updated_at, updated_by, id)
	} else {
		_, err = writeDB.Exec(stmtUpdatePost, title, slug, markdown, html, featured, isPage, status, meta_description, image, updated_at, updated_by, id)
	}
	if err != nil {
		_ = writeDB.Rollback()
		return err
	}
	return writeDB.Commit()
}

func UpdateSettings(title []byte, description []byte, logo []byte, cover []byte, postsPerPage int64, activeTheme string, navigation []byte, updated_at time.Time, updated_by int64) error {
	writeDB, err := readDB.Begin()
	if err != nil {
		_ = writeDB.Rollback()
		return err
	}
	// Title
	_, err = writeDB.Exec(stmtUpdateSettings, title, updated_at, updated_by, "title")
	if err != nil {
		_ = writeDB.Rollback()
		return err
	}
	// Description
	_, err = writeDB.Exec(stmtUpdateSettings, description, updated_at, updated_by, "description")
	if err != nil {
		_ = writeDB.Rollback()
		return err
	}
	// Logo
	_, err = writeDB.Exec(stmtUpdateSettings, logo, updated_at, updated_by, "logo")
	if err != nil {
		_ = writeDB.Rollback()
		return err
	}
	// Cover
	_, err = writeDB.Exec(stmtUpdateSettings, cover, updated_at, updated_by, "cover")
	if err != nil {
		_ = writeDB.Rollback()
		return err
	}
	// PostsPerPage
	_, err = writeDB.Exec(stmtUpdateSettings, postsPerPage, updated_at, updated_by, "postsPerPage")
	if err != nil {
		_ = writeDB.Rollback()
		return err
	}
	// ActiveTheme
	_, err = writeDB.Exec(stmtUpdateSettings, activeTheme, updated_at, updated_by, "activeTheme")
	if err != nil {
		_ = writeDB.Rollback()
		return err
	}
	// Navigation
	_, err = writeDB.Exec(stmtUpdateSettings, navigation, updated_at, updated_by, "navigation")
	if err != nil {
		_ = writeDB.Rollback()
		return err
	}
	return writeDB.Commit()
}

func UpdateActiveTheme(activeTheme string, updated_at time.Time, updated_by int64) error {
	writeDB, err := readDB.Begin()
	if err != nil {
		_ = writeDB.Rollback()
		return err
	}
	_, err = writeDB.Exec(stmtUpdateSettings, activeTheme, updated_at, updated_by, "activeTheme")
	if err != nil {
		_ = writeDB.Rollback()
		return err
	}
	return writeDB.Commit()
}

func UpdateUser(id int64, name []byte, slug string, email []byte, image []byte, cover []byte, bio []byte, website []byte, location []byte, updated_at time.Time, updated_by int64) error {
	writeDB, err := readDB.Begin()
	if err != nil {
		_ = writeDB.Rollback()
		return err
	}
	_, err = writeDB.Exec(stmtUpdateUser, name, slug, email, image, cover, bio, website, location, updated_at, updated_by, id)
	if err != nil {
		_ = writeDB.Rollback()
		return err
	}
	return writeDB.Commit()
}

func UpdateLastLogin(logInDate time.Time, userId int64) error {
	writeDB, err := readDB.Begin()
	if err != nil {
		_ = writeDB.Rollback()
		return err
	}
	_, err = writeDB.Exec(stmtUpdateLastLogin, logInDate, userId)
	if err != nil {
		_ = writeDB.Rollback()
		return err
	}
	return writeDB.Commit()
}

func UpdateUserPassword(id int64, password string, updated_at time.Time, updated_by int64) error {
	writeDB, err := readDB.Begin()
	if err != nil {
		_ = writeDB.Rollback()
		return err
	}
	_, err = writeDB.Exec(stmtUpdateUserPassword, password, updated_at, updated_by, id)
	if err != nil {
		_ = writeDB.Rollback()
		return err
	}
	return writeDB.Commit()
}
