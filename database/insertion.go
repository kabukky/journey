package database

import (
	"github.com/twinj/uuid"
	"time"
)

const stmtInsertPost = "INSERT INTO posts (id, uuid, title, slug, markdown, html, featured, page, status, image, author_id, created_at, created_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
const stmtInsertUser = "INSERT INTO users (id, uuid, name, slug, password, email, image, cover, created_at, created_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
const stmtInsertTag = "INSERT INTO tags (id, uuid, name, slug, created_at, created_by) VALUES (?, ?, ?, ?, ?, ?)"
const stmtInsertPostTag = "INSERT INTO posts_tags (id, post_id, tag_id) VALUES (?, ?, ?)"

func InsertPost(title []byte, slug string, markdown []byte, html []byte, featured bool, isPage bool, published bool, image []byte, created_at time.Time, created_by int64) (int64, error) {
	status := "draft"
	if published {
		status = "published"
	}
	writeDB, err := readDB.Begin()
	if err != nil {
		writeDB.Rollback()
		return 0, err
	}
	result, err := writeDB.Exec(stmtInsertPost, nil, uuid.Formatter(uuid.NewV4(), uuid.CleanHyphen), title, slug, markdown, html, featured, isPage, status, image, created_by, created_at, created_by)
	if err != nil {
		writeDB.Rollback()
		return 0, err
	}
	postId, err := result.LastInsertId()
	if err != nil {
		writeDB.Rollback()
		return 0, err
	}
	return postId, writeDB.Commit()
}

func InsertUser(name []byte, slug string, password string, email []byte, image []byte, cover []byte, created_at time.Time, created_by int64) error {
	writeDB, err := readDB.Begin()
	if err != nil {
		writeDB.Rollback()
		return err
	}
	_, err = writeDB.Exec(stmtInsertUser, nil, uuid.Formatter(uuid.NewV4(), uuid.CleanHyphen), name, slug, password, email, image, cover, created_at, created_by)
	if err != nil {
		writeDB.Rollback()
		return err
	}
	return writeDB.Commit()
}

func InsertTag(name []byte, slug string, created_at time.Time, created_by int64) (int64, error) {
	writeDB, err := readDB.Begin()
	if err != nil {
		writeDB.Rollback()
		return 0, err
	}
	result, err := writeDB.Exec(stmtInsertTag, nil, uuid.Formatter(uuid.NewV4(), uuid.CleanHyphen), name, slug, created_at, created_by)
	if err != nil {
		writeDB.Rollback()
		return 0, err
	}
	tagId, err := result.LastInsertId()
	if err != nil {
		writeDB.Rollback()
		return 0, err
	}
	return tagId, writeDB.Commit()
}

func InsertPostTag(post_id int64, tag_id int64) error {
	writeDB, err := readDB.Begin()
	if err != nil {
		writeDB.Rollback()
		return err
	}
	_, err = writeDB.Exec(stmtInsertPostTag, nil, post_id, tag_id)
	if err != nil {
		writeDB.Rollback()
		return err
	}
	return writeDB.Commit()
}
