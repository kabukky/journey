package database

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const stmtInsertPost = "INSERT INTO posts (id, uuid, title, slug, markdown, html, featured, page, status, meta_description, image, author_id, created_at, created_by, updated_at, updated_by, published_at, published_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
const stmtInsertUser = "INSERT INTO users (id, uuid, name, slug, password, email, image, cover, created_at, created_by, updated_at, updated_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
const stmtInsertRoleUser = "INSERT INTO roles_users (id, role_id, user_id) VALUES (?, ?, ?)"
const stmtInsertTag = "INSERT INTO tags (id, uuid, name, slug, created_at, created_by, updated_at, updated_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
const stmtInsertPostTag = "INSERT INTO posts_tags (id, post_id, tag_id) VALUES (?, ?, ?)"
const stmtInsertSetting = "INSERT INTO settings (id, uuid, key, value, type, created_at, created_by, updated_at, updated_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"

func InsertPost(title []byte, slug string, markdown []byte, html []byte, featured bool, isPage bool, published bool, meta_description []byte, image []byte, created_at time.Time, created_by int64) (int64, error) {

	status := "draft"
	if published {
		status = "published"
	}
	writeDB, err := readDB.Begin()
	if err != nil {
		_ = writeDB.Rollback()
		return 0, err
	}
	var result sql.Result
	if published {
		result, err = writeDB.Exec(stmtInsertPost, nil, uuid.New().String(), title, slug, markdown, html, featured, isPage, status, meta_description, image, created_by, created_at, created_by, created_at, created_by, created_at, created_by)
	} else {
		result, err = writeDB.Exec(stmtInsertPost, nil, uuid.New().String(), title, slug, markdown, html, featured, isPage, status, meta_description, image, created_by, created_at, created_by, created_at, created_by, nil, nil)
	}
	if err != nil {
		_ = writeDB.Rollback()
		return 0, err
	}
	postId, err := result.LastInsertId()
	if err != nil {
		_ = writeDB.Rollback()
		return 0, err
	}
	return postId, writeDB.Commit()
}

func InsertUser(name []byte, slug string, password string, email []byte, image []byte, cover []byte, created_at time.Time, created_by int64) (int64, error) {
	writeDB, err := readDB.Begin()
	if err != nil {
		_ = writeDB.Rollback()
		return 0, err
	}
	result, err := writeDB.Exec(stmtInsertUser, nil, uuid.New().String(), name, slug, password, email, image, cover, created_at, created_by, created_at, created_by)
	if err != nil {
		_ = writeDB.Rollback()
		return 0, err
	}
	userId, err := result.LastInsertId()
	if err != nil {
		_ = writeDB.Rollback()
		return 0, err
	}
	return userId, writeDB.Commit()
}

func InsertRoleUser(role_id int, user_id int64) error {
	writeDB, err := readDB.Begin()
	if err != nil {
		_ = writeDB.Rollback()
		return err
	}
	_, err = writeDB.Exec(stmtInsertRoleUser, nil, role_id, user_id)
	if err != nil {
		_ = writeDB.Rollback()
		return err
	}
	return writeDB.Commit()
}

func InsertTag(name []byte, slug string, created_at time.Time, created_by int64) (int64, error) {
	writeDB, err := readDB.Begin()
	if err != nil {
		_ = writeDB.Rollback()
		return 0, err
	}
	result, err := writeDB.Exec(stmtInsertTag, nil, uuid.New().String(), name, slug, created_at, created_by, created_at, created_by)
	if err != nil {
		_ = writeDB.Rollback()
		return 0, err
	}
	tagId, err := result.LastInsertId()
	if err != nil {
		_ = writeDB.Rollback()
		return 0, err
	}
	return tagId, writeDB.Commit()
}

func InsertPostTag(post_id int64, tag_id int64) error {
	writeDB, err := readDB.Begin()
	if err != nil {
		_ = writeDB.Rollback()
		return err
	}
	_, err = writeDB.Exec(stmtInsertPostTag, nil, post_id, tag_id)
	if err != nil {
		_ = writeDB.Rollback()
		return err
	}
	return writeDB.Commit()
}

func insertSettingString(key string, value string, setting_type string, created_at time.Time, created_by int64) error {
	writeDB, err := readDB.Begin()
	if err != nil {
		_ = writeDB.Rollback()
		return err
	}
	_, err = writeDB.Exec(stmtInsertSetting, nil, uuid.New().String(), key, value, setting_type, created_at, created_by, created_at, created_by)
	if err != nil {
		_ = writeDB.Rollback()
		return err
	}
	return writeDB.Commit()
}

func insertSettingInt64(key string, value int64, setting_type string, created_at time.Time, created_by int64) error {
	writeDB, err := readDB.Begin()
	if err != nil {
		_ = writeDB.Rollback()
		return err
	}
	_, err = writeDB.Exec(stmtInsertSetting, nil, uuid.New().String(), key, value, setting_type, created_at, created_by, created_at, created_by)
	if err != nil {
		_ = writeDB.Rollback()
		return err
	}
	return writeDB.Commit()
}
