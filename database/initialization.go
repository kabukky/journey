package database

import (
	"database/sql"
	"github.com/kabukky/journey/filenames"
	_ "github.com/mattn/go-sqlite3"
	"github.com/twinj/uuid"
	"time"
)

var readDB *sql.DB // Handler for read access

var stmtInitialization = `CREATE TABLE IF NOT EXISTS
	posts (
		id					integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		uuid				varchar(36) NOT NULL,
		title				varchar(150) NOT NULL,
		slug				varchar(150) NOT NULL,
		markdown			text,
		html				text,
		image				text,
		featured			tinyint NOT NULL DEFAULT '0',
		page				tinyint NOT NULL DEFAULT '0',
		status				varchar(150) NOT NULL DEFAULT 'draft',
		language			varchar(6) NOT NULL DEFAULT 'en_US',
		meta_title			varchar(150),
		meta_description	varchar(200),
		author_id			integer NOT NULL,
		created_at			datetime NOT NULL,
		created_by			integer NOT NULL,
		updated_at			datetime,
		updated_by			integer,
		published_at		datetime,
		published_by		integer
	);
	CREATE TABLE IF NOT EXISTS
	users (
		id					integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		uuid				varchar(36) NOT NULL,
		name				varchar(150) NOT NULL,
		slug				varchar(150) NOT NULL,
		password			varchar(60) NOT NULL,
		email				varchar(254) NOT NULL,
		image				text,
		cover				text,
		bio					varchar(200),
		website				text,
		location			text,
		accessibility		text,
		status				varchar(150) NOT NULL DEFAULT 'active',
		language			varchar(6) NOT NULL DEFAULT 'en_US',
		meta_title			varchar(150),
		meta_description	varchar(200),
		last_login			datetime,
		created_at			datetime NOT NULL,
		created_by			integer NOT NULL,
		updated_at			datetime,
		updated_by			integer
	);
	CREATE TABLE IF NOT EXISTS
	tags (
		id					integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		uuid				varchar(36) NOT NULL,
		name				varchar(150) NOT NULL,
		slug				varchar(150) NOT NULL,
		description			varchar(200),
		parent_id			integer,
		meta_title			varchar(150),
		meta_description	varchar(200),
		created_at			datetime NOT NULL,
		created_by			integer NOT NULL,
		updated_at			datetime,
		updated_by			integer
	);
	CREATE TABLE IF NOT EXISTS
	posts_tags (
		id		integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		post_id	integer NOT NULL,
		tag_id	integer NOT NULL
	);
	CREATE TABLE IF NOT EXISTS
	settings (
		id			integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		uuid		varchar(36) NOT NULL,
		key			varchar(150) NOT NULL,
		value		text,
		type		varchar(150) NOT NULL DEFAULT 'core',
		created_at	datetime NOT NULL,
		created_by	integer NOT NULL,
		updated_at	datetime,
		updated_by	integer
	);
	INSERT OR IGNORE INTO settings (id, uuid, key, value, type, created_at, created_by, updated_at, updated_by) VALUES (1, ?, 'title', 'My Blog', 'blog', ?, 1, ?, 1);
	INSERT OR IGNORE INTO settings (id, uuid, key, value, type, created_at, created_by, updated_at, updated_by) VALUES (2, ?, 'description', 'Just another Blog', 'blog', ?, 1, ?, 1);
	INSERT OR IGNORE INTO settings (id, uuid, key, value, type, created_at, created_by, updated_at, updated_by) VALUES (3, ?, 'email', '', 'blog', ?, 1, ?, 1);
	INSERT OR IGNORE INTO settings (id, uuid, key, value, type, created_at, created_by, updated_at, updated_by) VALUES (4, ?, 'logo', '/public/images/blog-logo.jpg', 'blog', ?, 1, ?, 1);
	INSERT OR IGNORE INTO settings (id, uuid, key, value, type, created_at, created_by, updated_at, updated_by) VALUES (5, ?, 'cover', '/public/images/blog-cover.jpg', 'blog', ?, 1, ?, 1);
	INSERT OR IGNORE INTO settings (id, uuid, key, value, type, created_at, created_by, updated_at, updated_by) VALUES (6, ?, 'postsPerPage', 5, 'blog', ?, 1, ?, 1);
	INSERT OR IGNORE INTO settings (id, uuid, key, value, type, created_at, created_by, updated_at, updated_by) VALUES (7, ?, 'activeTheme', 'promenade', 'theme', ?, 1, ?, 1);
	`

func Initialize() error {
	// TODO: If there is no journey.db, convert ghost database if available (time format needs to change to be compatible with journey)
	// Open or create database file
	var err error
	readDB, err = sql.Open("sqlite3", filenames.DatabaseFilename)
	if err != nil {
		return err
	}
	readDB.SetMaxIdleConns(256) // TODO: is this enough?
	err = readDB.Ping()
	if err != nil {
		return err
	}
	currentTime := time.Now()
	_, err = readDB.Exec(stmtInitialization, uuid.Formatter(uuid.NewV4(), uuid.CleanHyphen), currentTime, currentTime, uuid.Formatter(uuid.NewV4(), uuid.CleanHyphen), currentTime, currentTime, uuid.Formatter(uuid.NewV4(), uuid.CleanHyphen), currentTime, currentTime, uuid.Formatter(uuid.NewV4(), uuid.CleanHyphen), currentTime, currentTime, uuid.Formatter(uuid.NewV4(), uuid.CleanHyphen), currentTime, currentTime, uuid.Formatter(uuid.NewV4(), uuid.CleanHyphen), currentTime, currentTime, uuid.Formatter(uuid.NewV4(), uuid.CleanHyphen), currentTime, currentTime)
	// TODO: Is Commit()/Rollback() needed for DB.Exec()?
	if err != nil {
		return err
	}
	return nil
}
