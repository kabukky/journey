package database

import (
	"database/sql"
	"time"

	"github.com/gofrs/uuid"
)

const stmtInsertPost = "INSERT INTO posts (id, uuid, title, slug, markdown, html, featured, page, status, meta_description, image, author_id, created_at, created_by, updated_at, updated_by, published_at, published_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
const stmtInsertUser = "INSERT INTO users (id, uuid, name, slug, password, email, image, cover, created_at, created_by, updated_at, updated_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
const stmtInsertRoleUser = "INSERT INTO roles_users (id, role_id, user_id) VALUES (?, ?, ?)"
const stmtInsertTag = "INSERT INTO tags (id, uuid, name, slug, created_at, created_by, updated_at, updated_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
const stmtInsertPostTag = "INSERT INTO posts_tags (id, post_id, tag_id) VALUES (?, ?, ?)"
const stmtInsertSetting = "INSERT INTO settings (id, uuid, key, value, type, created_at, created_by, updated_at, updated_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"

// InsertPost inserts a post
func InsertPost(title []byte, slug string, markdown []byte, html []byte, featured bool, isPage bool, published bool, metaDescription []byte, image []byte, createdAt time.Time, createdBy int64) (int64, error) {

	status := "draft"
	if published {
		status = "published"
	}
	writeDB, err := readDB.Begin()
	if err != nil {
		writeDB.Rollback()
		return 0, err
	}
	var result sql.Result
	if published {
		result, err = writeDB.Exec(stmtInsertPost, nil, uuid.Must(uuid.NewV4()).String(), title, slug, markdown, html, featured, isPage, status, metaDescription, image, createdBy, createdAt, createdBy, createdAt, createdBy, createdAt, createdBy)
	} else {
		result, err = writeDB.Exec(stmtInsertPost, nil, uuid.Must(uuid.NewV4()).String(), title, slug, markdown, html, featured, isPage, status, metaDescription, image, createdBy, createdAt, createdBy, createdAt, createdBy, nil, nil)
	}
	if err != nil {
		writeDB.Rollback()
		return 0, err
	}
	postID, err := result.LastInsertId()
	if err != nil {
		writeDB.Rollback()
		return 0, err
	}
	return postID, writeDB.Commit()
}

// InsertUser inserts a user
func InsertUser(name []byte, slug string, password string, email []byte, image []byte, cover []byte, createdAt time.Time, createdBy int64) (int64, error) {
	writeDB, err := readDB.Begin()
	if err != nil {
		writeDB.Rollback()
		return 0, err
	}
	result, err := writeDB.Exec(stmtInsertUser, nil, uuid.Must(uuid.NewV4()).String(), name, slug, password, email, image, cover, createdAt, createdBy, createdAt, createdBy)
	if err != nil {
		writeDB.Rollback()
		return 0, err
	}
	userID, err := result.LastInsertId()
	if err != nil {
		writeDB.Rollback()
		return 0, err
	}
	return userID, writeDB.Commit()
}

// InsertRoleUser inserts a role for a user
func InsertRoleUser(roleID int, userID int64) error {
	writeDB, err := readDB.Begin()
	if err != nil {
		writeDB.Rollback()
		return err
	}
	_, err = writeDB.Exec(stmtInsertRoleUser, nil, roleID, userID)
	if err != nil {
		writeDB.Rollback()
		return err
	}
	return writeDB.Commit()
}

// InsertTag inserts a tag
func InsertTag(name []byte, slug string, createdAt time.Time, createdBy int64) (int64, error) {
	writeDB, err := readDB.Begin()
	if err != nil {
		writeDB.Rollback()
		return 0, err
	}
	result, err := writeDB.Exec(stmtInsertTag, nil, uuid.Must(uuid.NewV4()).String(), name, slug, createdAt, createdBy, createdAt, createdBy)
	if err != nil {
		writeDB.Rollback()
		return 0, err
	}
	tagID, err := result.LastInsertId()
	if err != nil {
		writeDB.Rollback()
		return 0, err
	}
	return tagID, writeDB.Commit()
}

// InsertPostTag inserts a post tag
func InsertPostTag(postID int64, tagID int64) error {
	writeDB, err := readDB.Begin()
	if err != nil {
		writeDB.Rollback()
		return err
	}
	_, err = writeDB.Exec(stmtInsertPostTag, nil, postID, tagID)
	if err != nil {
		writeDB.Rollback()
		return err
	}
	return writeDB.Commit()
}

func insertSettingString(key string, value string, settingType string, createdAt time.Time, createdBy int64) error {
	writeDB, err := readDB.Begin()
	if err != nil {
		writeDB.Rollback()
		return err
	}
	_, err = writeDB.Exec(stmtInsertSetting, nil, uuid.Must(uuid.NewV4()).String(), key, value, settingType, createdAt, createdBy, createdAt, createdBy)
	if err != nil {
		writeDB.Rollback()
		return err
	}
	return writeDB.Commit()
}

func insertSettingInt64(key string, value int64, settingType string, createdAt time.Time, createdBy int64) error {
	writeDB, err := readDB.Begin()
	if err != nil {
		writeDB.Rollback()
		return err
	}
	_, err = writeDB.Exec(stmtInsertSetting, nil, uuid.Must(uuid.NewV4()).String(), key, value, settingType, createdAt, createdBy, createdAt, createdBy)
	if err != nil {
		writeDB.Rollback()
		return err
	}
	return writeDB.Commit()
}
