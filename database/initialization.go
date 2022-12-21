package database

import (
	"database/sql"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"

	"github.com/kabukky/journey/database/migration"
	"github.com/kabukky/journey/date"
	"github.com/kabukky/journey/filenames"
	"github.com/kabukky/journey/helpers"
	"github.com/kabukky/journey/structure"
)

// Handler for read access
var readDB *sql.DB

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
	INSERT OR IGNORE INTO settings (id, uuid, key, value, type, created_at, created_by, updated_at, updated_by) VALUES (8, ?, 'navigation', '[{"label":"Home", "url":"/"}]', 'blog', ?, 1, ?, 1);
	INSERT OR IGNORE INTO settings (id, uuid, key, value, type, created_at, created_by, updated_at, updated_by) VALUES (9, ?, 'ghost_head', '<!-- Custom header -->', 'blog', ?, 1, ?, 1);
	CREATE TABLE IF NOT EXISTS
	roles (
		id			integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		uuid		varchar(36) NOT NULL,
		name		varchar(150) NOT NULL,
		description	varchar(200),
		created_at	datetime NOT NULL,
		created_by	integer NOT NULL,
		updated_at	datetime,
		updated_by	integer
	);
	INSERT OR IGNORE INTO roles (id, uuid, name, description, created_at, created_by, updated_at, updated_by) VALUES (1, ?, 'Administrator', 'Administrators', ?, 1, ?, 1);
	INSERT OR IGNORE INTO roles (id, uuid, name, description, created_at, created_by, updated_at, updated_by) VALUES (2, ?, 'Editor', 'Editors', ?, 1, ?, 1);
	INSERT OR IGNORE INTO roles (id, uuid, name, description, created_at, created_by, updated_at, updated_by) VALUES (3, ?, 'Author', 'Authors', ?, 1, ?, 1);
	INSERT OR IGNORE INTO roles (id, uuid, name, description, created_at, created_by, updated_at, updated_by) VALUES (4, ?, 'Owner', 'Blog Owner', ?, 1, ?, 1);
	CREATE TABLE IF NOT EXISTS
	roles_users (
		id		integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		role_id	integer NOT NULL,
		user_id	integer NOT NULL
	);
	`

func Initialize() error {
	// If journey.db does not exist, look for a Ghost database to convert
	if !helpers.FileExists(filenames.DatabaseFilename) {
		// Convert Ghost database if available (time format needs to change to be compatible with journey)
		migration.Ghost()
	}
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
	currentTime := date.GetCurrentTime()
	_, err = readDB.Exec(stmtInitialization, uuid.New().String(), currentTime, currentTime, uuid.New().String(), currentTime, currentTime, uuid.New().String(), currentTime, currentTime, uuid.New().String(), currentTime, currentTime, uuid.New().String(), currentTime, currentTime, uuid.New().String(), currentTime, currentTime, uuid.New().String(), currentTime, currentTime, uuid.New().String(), currentTime, currentTime, uuid.New().String(), currentTime, currentTime, uuid.New().String(), currentTime, currentTime, uuid.New().String(), currentTime, currentTime, uuid.New().String(), currentTime, currentTime, uuid.New().String(), currentTime, currentTime)
	// TODO: Is Commit()/Rollback() needed for DB.Exec()?
	if err != nil {
		return err
	}
	err = checkBlogSettings()
	if err != nil {
		return err
	}
	return nil
}

// Function to check and insert any missing blog settings into the database (settings could be missing if migrating from Ghost).
func checkBlogSettings() error {
	tempBlog := structure.Blog{}
	// Check for title
	row := readDB.QueryRow(stmtRetrieveBlog, "title")
	err := row.Scan(&tempBlog.Title)
	if err != nil {
		// Insert title
		err = insertSettingString("title", "My Blog", "blog", date.GetCurrentTime(), 1)
		if err != nil {
			return err
		}
	}
	// Check for description
	row = readDB.QueryRow(stmtRetrieveBlog, "description")
	err = row.Scan(&tempBlog.Description)
	if err != nil {
		// Insert description
		err = insertSettingString("description", "Just another Blog", "blog", date.GetCurrentTime(), 1)
		if err != nil {
			return err
		}
	}
	// Check for email
	var email []byte
	row = readDB.QueryRow(stmtRetrieveBlog, "email")
	err = row.Scan(&email)
	if err != nil {
		// Insert email
		err = insertSettingString("email", "", "blog", date.GetCurrentTime(), 1)
		if err != nil {
			return err
		}
	}
	// Check for logo
	row = readDB.QueryRow(stmtRetrieveBlog, "logo")
	err = row.Scan(&tempBlog.Logo)
	if err != nil {
		// Insert logo
		err = insertSettingString("logo", "/public/images/blog-logo.jpg", "blog", date.GetCurrentTime(), 1)
		if err != nil {
			return err
		}
	}
	// Check for cover
	row = readDB.QueryRow(stmtRetrieveBlog, "cover")
	err = row.Scan(&tempBlog.Cover)
	if err != nil {
		// Insert cover
		err = insertSettingString("cover", "/public/images/blog-cover.jpg", "blog", date.GetCurrentTime(), 1)
		if err != nil {
			return err
		}
	}
	// Check for postsPerPage
	row = readDB.QueryRow(stmtRetrieveBlog, "postsPerPage")
	err = row.Scan(&tempBlog.PostsPerPage)
	if err != nil {
		// Insert postsPerPage
		err = insertSettingInt64("postsPerPage", 5, "blog", date.GetCurrentTime(), 1)
		if err != nil {
			return err
		}
	}
	// Check for activeTheme
	row = readDB.QueryRow(stmtRetrieveBlog, "activeTheme")
	err = row.Scan(&tempBlog.ActiveTheme)
	if err != nil {
		// Insert activeTheme
		err = insertSettingString("activeTheme", "promenade", "theme", date.GetCurrentTime(), 1)
		if err != nil {
			return err
		}
	}
	// Check for navigation
	var navigation []byte
	row = readDB.QueryRow(stmtRetrieveBlog, "navigation")
	err = row.Scan(&navigation)
	if err != nil {
		// Insert navigation
		err = insertSettingString("navigation", "[{\"label\":\"Home\", \"url\":\"/\"}]", "blog", date.GetCurrentTime(), 1)
		if err != nil {
			return err
		}
	}
	return nil
}
