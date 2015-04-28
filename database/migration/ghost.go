package migration

import (
	"database/sql"
	"errors"
	"github.com/kabukky/journey/filenames"
	"github.com/kabukky/journey/helpers"
	"log"
	"os"
	"path/filepath"
	"time"
)

const stmtRetrieveGhostPosts = "SELECT id, (created_at/1000), (updated_at/1000), (published_at/1000) FROM posts"
const stmtRetrieveGhostTags = "SELECT id, (created_at/1000), (updated_at/1000) FROM tags"
const stmtRetrieveGhostUsers = "SELECT id, name, email, (last_login/1000), (created_at/1000), (updated_at/1000) FROM users"
const stmtRetrieveGhostRoles = "SELECT id, (created_at/1000), (updated_at/1000) FROM roles"
const stmtRetrieveGhostSettings = "SELECT id, (created_at/1000), (updated_at/1000) FROM settings"
const stmtRetrieveGhostPermissions = "SELECT id, (created_at/1000), (updated_at/1000) FROM permissions"
const stmtRetrieveGhostClients = "SELECT id, (created_at/1000), (updated_at/1000) FROM clients"

const stmtUpdateGhostPost = "UPDATE posts SET created_at = ?, updated_at = ?, published_at = ? WHERE id = ?"
const stmtUpdateGhostTags = "UPDATE tags SET created_at = ?, updated_at = ? WHERE id = ?"
const stmtUpdateGhostUsers = "UPDATE users SET name = ?, email= ?, last_login = ?, created_at = ?, updated_at = ? WHERE id = ?"
const stmtUpdateGhostRoles = "UPDATE roles SET created_at = ?, updated_at = ? WHERE id = ?"
const stmtUpdateGhostSettings = "UPDATE settings SET created_at = ?, updated_at = ? WHERE id = ?"
const stmtUpdateGhostPermissions = "UPDATE permissions SET created_at = ?, updated_at = ? WHERE id = ?"
const stmtUpdateGhostClients = "UPDATE clients SET created_at = ?, updated_at = ? WHERE id = ?"
const stmtUpdateGhostTheme = "UPDATE settings SET value = ?, updated_at = ?, updated_by = ? WHERE key = 'activeTheme'"

type dateHolder struct {
	id          int64
	name        []byte
	email       []byte
	createdAt   *time.Time
	updatedAt   *time.Time
	publishedAt *time.Time
	lastLogin   *time.Time
}

// Function to convert a Ghost database to use with Journey
func Ghost() {
	// Check every file in data directory
	err := filepath.Walk(filenames.DatabaseFilepath, inspectDatabaseFile)
	if err != nil {
		log.Println("Error while looking for a Ghost database to convert:", err)
		return
	}
}

func inspectDatabaseFile(filePath string, info os.FileInfo, err error) error {
	if !info.IsDir() && filepath.Ext(filePath) == ".db" {
		err := convertGhostDatabase(filePath)
		if err != nil {
			return err
		}
	}
	return nil
}

// This function converts all fields in the Ghost db that are not compatible with Journey (only date fields for now. Ghost uses a javascript-specific unix timestamp).
func convertGhostDatabase(fileName string) error {
	// If journey.db exists already, don't convert this file
	if helpers.FileExists(filenames.DatabaseFilename) {
		return errors.New(filenames.DatabaseFilename + " already exists.")
	}
	log.Println("Trying to convert " + fileName + "...")
	readDB, err := sql.Open("sqlite3", fileName)
	if err != nil {
		log.Println("Error:", err)
		return err
	}
	err = readDB.Ping()
	if err != nil {
		log.Println("Error:", err)
		return err
	}
	// Convert posts
	err = convertPosts(readDB)
	if err != nil {
		log.Println("Error:", err)
		return err
	}
	// Convert users
	err = convertUsers(readDB)
	if err != nil {
		log.Println("Error:", err)
		return err
	}
	// Convert tags
	err = convertDates(readDB, stmtRetrieveGhostTags, stmtUpdateGhostTags)
	if err != nil {
		log.Println("Error:", err)
		return err
	}
	// Convert roles
	err = convertDates(readDB, stmtRetrieveGhostRoles, stmtUpdateGhostRoles)
	if err != nil {
		log.Println("Error:", err)
		return err
	}
	// Convert settings
	err = convertDates(readDB, stmtRetrieveGhostSettings, stmtUpdateGhostSettings)
	if err != nil {
		log.Println("Error:", err)
		return err
	}
	// Set default theme
	err = setDefaultTheme(readDB)
	if err != nil {
		log.Println("Error:", err)
		return err
	}
	// Convert permissions (not used by Journey at the moment)
	err = convertDates(readDB, stmtRetrieveGhostPermissions, stmtUpdateGhostPermissions)
	if err != nil {
		log.Println("Error:", err)
		return err
	}
	// Convert clients (not used by Journey at the moment)
	err = convertDates(readDB, stmtRetrieveGhostClients, stmtUpdateGhostClients)
	if err != nil {
		log.Println("Error:", err)
		return err
	}
	// All went well. Close database connection.
	err = readDB.Close()
	if err != nil {
		log.Println("Error:", err)
		return err
	}
	// Rename file to Journey format
	err = os.Rename(fileName, filenames.DatabaseFilename)
	if err != nil {
		log.Println("Error:", err)
		return err
	}
	log.Println("Success!")
	return nil
}

func convertPosts(readDB *sql.DB) error {
	allRows := make([]dateHolder, 0)
	// Retrieve posts
	rows, err := readDB.Query(stmtRetrieveGhostPosts)
	if err != nil {
		return err
	}
	// Read all rows to structs
	for rows.Next() {
		row := dateHolder{}
		var createdAt sql.NullInt64
		var updatedAt sql.NullInt64
		var publishedAt sql.NullInt64
		err := rows.Scan(&row.id, &createdAt, &updatedAt, &publishedAt)
		if err != nil {
			return err
		}
		// Convert dates
		if createdAt.Valid {
			date := time.Unix(createdAt.Int64, 0)
			row.createdAt = &date
		}
		if updatedAt.Valid {
			date := time.Unix(updatedAt.Int64, 0)
			row.updatedAt = &date
		}
		if publishedAt.Valid {
			date := time.Unix(publishedAt.Int64, 0)
			row.publishedAt = &date
		}
		allRows = append(allRows, row)
	}
	rows.Close()
	// Write all new dates
	for _, row := range allRows {
		writeDB, err := readDB.Begin()
		if err != nil {
			writeDB.Rollback()
			return err
		}
		// Update the database with new date formats
		_, err = writeDB.Exec(stmtUpdateGhostPost, row.createdAt, row.updatedAt, row.publishedAt, row.id)
		if err != nil {
			writeDB.Rollback()
			return err
		}
		err = writeDB.Commit()
		if err != nil {
			writeDB.Rollback()
			return err
		}
	}
	return nil
}

func convertUsers(readDB *sql.DB) error {
	allRows := make([]dateHolder, 0)
	// Retrieve posts
	rows, err := readDB.Query(stmtRetrieveGhostUsers)
	if err != nil {
		return err
	}
	// Read all rows to structs
	for rows.Next() {
		row := dateHolder{}
		var name sql.NullString
		var email sql.NullString
		var lastLogin sql.NullInt64
		var createdAt sql.NullInt64
		var updatedAt sql.NullInt64
		err := rows.Scan(&row.id, &name, &email, &lastLogin, &createdAt, &updatedAt)
		if err != nil {
			return err
		}
		// Convert strings to byte array since that is how Journey saves the user name and email (Login won't work wihout this).
		if name.Valid {
			row.name = []byte(name.String)
		}
		if email.Valid {
			row.email = []byte(email.String)
		}
		// Convert dates
		if lastLogin.Valid {
			date := time.Unix(lastLogin.Int64, 0)
			row.lastLogin = &date
		}
		if createdAt.Valid {
			date := time.Unix(createdAt.Int64, 0)
			row.createdAt = &date
		}
		if updatedAt.Valid {
			date := time.Unix(updatedAt.Int64, 0)
			row.updatedAt = &date
		}
		allRows = append(allRows, row)
	}
	rows.Close()
	// Write all new dates
	for _, row := range allRows {
		writeDB, err := readDB.Begin()
		if err != nil {
			writeDB.Rollback()
			return err
		}
		// Update the database with new date formats
		_, err = writeDB.Exec(stmtUpdateGhostUsers, row.name, row.email, row.lastLogin, row.createdAt, row.updatedAt, row.id)
		if err != nil {
			writeDB.Rollback()
			return err
		}
		err = writeDB.Commit()
		if err != nil {
			writeDB.Rollback()
			return err
		}
	}
	return nil
}

func convertDates(readDB *sql.DB, stmtRetrieve string, stmtUpdate string) error {
	allRows := make([]dateHolder, 0)
	// Retrieve posts
	rows, err := readDB.Query(stmtRetrieve)
	if err != nil {
		return err
	}
	// Read all rows to structs
	for rows.Next() {
		row := dateHolder{}
		var createdAt sql.NullInt64
		var updatedAt sql.NullInt64
		err := rows.Scan(&row.id, &createdAt, &updatedAt)
		if err != nil {
			return err
		}
		// Convert dates
		if createdAt.Valid {
			date := time.Unix(createdAt.Int64, 0)
			row.createdAt = &date
		}
		if updatedAt.Valid {
			date := time.Unix(updatedAt.Int64, 0)
			row.updatedAt = &date
		}
		allRows = append(allRows, row)
	}
	rows.Close()
	// Write all new dates
	for _, row := range allRows {
		writeDB, err := readDB.Begin()
		if err != nil {
			writeDB.Rollback()
			return err
		}
		// Update the database with new date formats
		_, err = writeDB.Exec(stmtUpdate, row.createdAt, row.updatedAt, row.id)
		if err != nil {
			writeDB.Rollback()
			return err
		}
		err = writeDB.Commit()
		if err != nil {
			writeDB.Rollback()
			return err
		}
	}
	return nil
}

func setDefaultTheme(readDB *sql.DB) error {
	writeDB, err := readDB.Begin()
	if err != nil {
		writeDB.Rollback()
		return err
	}
	// Update the database with the default theme (promenade)
	date := time.Now()
	_, err = writeDB.Exec(stmtUpdateGhostTheme, "promenade", date, 1)
	if err != nil {
		writeDB.Rollback()
		return err
	}
	return writeDB.Commit()
}
