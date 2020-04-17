package database

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/kabukky/journey/structure"
)

const stmtRetrievePostsCount = "SELECT count(*) FROM posts WHERE page = 0 AND status = 'published'"
const stmtRetrievePostsCountByUser = "SELECT count(*) FROM posts WHERE page = 0 AND status = 'published' AND author_id = ?"
const stmtRetrievePostsCountByTag = "SELECT count(*) FROM posts, posts_tags WHERE posts_tags.post_id = posts.id AND posts_tags.tag_id = ? AND page = 0 AND status = 'published'"
const stmtRetrievePostsForIndex = "SELECT id, uuid, title, slug, markdown, html, featured, page, status, meta_description, image, author_id, published_at FROM posts WHERE page = 0 AND status = 'published' ORDER BY published_at DESC LIMIT ? OFFSET ?"
const stmtRetrievePostsForAPI = "SELECT id, uuid, title, slug, markdown, html, featured, page, status, meta_description, image, author_id, published_at FROM posts ORDER BY id DESC LIMIT ? OFFSET ?"
const stmtRetrievePostsByUser = "SELECT id, uuid, title, slug, markdown, html, featured, page, status, meta_description, image, author_id, published_at FROM posts WHERE page = 0 AND status = 'published' AND author_id = ? ORDER BY published_at DESC LIMIT ? OFFSET ?"
const stmtRetrievePostsByTag = "SELECT posts.id, posts.uuid, posts.title, posts.slug, posts.markdown, posts.html, posts.featured, posts.page, posts.status, posts.meta_description, posts.image, posts.author_id, posts.published_at FROM posts, posts_tags WHERE posts_tags.post_id = posts.id AND posts_tags.tag_id = ? AND page = 0 AND status = 'published' ORDER BY posts.published_at DESC LIMIT ? OFFSET ?"
const stmtRetrievePostByID = "SELECT id, uuid, title, slug, markdown, html, featured, page, status, meta_description, image, author_id, published_at FROM posts WHERE id = ?"
const stmtRetrievePostBySlug = "SELECT id, uuid, title, slug, markdown, html, featured, page, status, meta_description, image, author_id, published_at FROM posts WHERE slug = ?"
const stmtRetrieveUserByID = "SELECT id, name, slug, email, image, cover, bio, website, location FROM users WHERE id = ?"
const stmtRetrieveUserBySlug = "SELECT id, name, slug, email, image, cover, bio, website, location FROM users WHERE slug = ?"
const stmtRetrieveUserByName = "SELECT id, name, slug, email, image, cover, bio, website, location FROM users WHERE name = ?"
const stmtRetrieveTags = "SELECT tag_id FROM posts_tags WHERE post_id = ?"
const stmtRetrieveTagByID = "SELECT id, name, slug FROM tags WHERE id = ?"
const stmtRetrieveTagBySlug = "SELECT id, name, slug FROM tags WHERE slug = ?"
const stmtRetrieveTagIDBySlug = "SELECT id FROM tags WHERE slug = ?"
const stmtRetrieveHashedPasswordByName = "SELECT password FROM users WHERE name = ?"
const stmtRetrieveUsersCount = "SELECT count(*) FROM users"
const stmtRetrieveBlog = "SELECT value FROM settings WHERE key = ?"
const stmtRetrievePostCreationDateByID = "SELECT created_at FROM posts WHERE id = ?"

// RetrievePostByID ...
func RetrievePostByID(id int64) (*structure.Post, error) {
	// Retrieve post
	row := readDB.QueryRow(stmtRetrievePostByID, id)
	return extractPost(row)
}

// RetrievePostBySlug ...
func RetrievePostBySlug(slug string) (*structure.Post, error) {
	// Retrieve post
	row := readDB.QueryRow(stmtRetrievePostBySlug, slug)
	return extractPost(row)
}

// RetrievePostsByUser ...
func RetrievePostsByUser(userID int64, limit int64, offset int64) ([]structure.Post, error) {
	// Retrieve posts
	rows, err := readDB.Query(stmtRetrievePostsByUser, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	posts, err := extractPosts(rows)
	if err != nil {
		return nil, err
	}
	return *posts, nil
}

// RetrievePostsByTag ...
func RetrievePostsByTag(tagID int64, limit int64, offset int64) ([]structure.Post, error) {
	// Retrieve posts
	rows, err := readDB.Query(stmtRetrievePostsByTag, tagID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	posts, err := extractPosts(rows)
	if err != nil {
		return nil, err
	}
	return *posts, nil
}

// RetrievePostsForIndex ...
func RetrievePostsForIndex(limit int64, offset int64) ([]structure.Post, error) {
	// Retrieve posts
	rows, err := readDB.Query(stmtRetrievePostsForIndex, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	posts, err := extractPosts(rows)
	if err != nil {
		return nil, err
	}
	return *posts, nil
}

// RetrievePostsForAPI ...
func RetrievePostsForAPI(limit int64, offset int64) ([]structure.Post, error) {
	// Retrieve posts
	rows, err := readDB.Query(stmtRetrievePostsForAPI, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	posts, err := extractPosts(rows)
	if err != nil {
		return nil, err
	}
	return *posts, nil
}

func extractPosts(rows *sql.Rows) (*[]structure.Post, error) {
	posts := make([]structure.Post, 0)
	for rows.Next() {
		post := structure.Post{}
		var userID int64
		var status string
		err := rows.Scan(&post.ID, &post.UUID, &post.Title, &post.Slug, &post.Markdown, &post.HTML, &post.IsFeatured, &post.IsPage, &status, &post.MetaDescription, &post.Image, &userID, &post.Date)
		if err != nil {
			return nil, err
		}
		// If there was no publication date attached to the post, make its creation date the date of the post
		if post.Date == nil {
			post.Date, err = retrievePostCreationDateByID(post.ID)
			if err != nil {
				return nil, err
			}
		}
		// Evaluate status
		if status == "published" {
			post.IsPublished = true
		} else {
			post.IsPublished = false
		}
		// Retrieve user
		post.Author, err = RetrieveUser(userID)
		if err != nil {
			return nil, err
		}
		// Retrieve tags
		post.Tags, err = RetrieveTags(post.ID)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return &posts, nil
}

func extractPost(row *sql.Row) (*structure.Post, error) {
	post := structure.Post{}
	var userID int64
	var status string
	err := row.Scan(&post.ID, &post.UUID, &post.Title, &post.Slug, &post.Markdown, &post.HTML, &post.IsFeatured, &post.IsPage, &status, &post.MetaDescription, &post.Image, &userID, &post.Date)
	if err != nil {
		return nil, err
	}
	// If there was no publication date attached to the post, make its creation date the date of the post
	if post.Date == nil {
		post.Date, err = retrievePostCreationDateByID(post.ID)
		if err != nil {
			return nil, err
		}
	}
	// Evaluate status
	if status == "published" {
		post.IsPublished = true
	} else {
		post.IsPublished = false
	}
	// Retrieve user
	post.Author, err = RetrieveUser(userID)
	if err != nil {
		return nil, err
	}
	// Retrieve tags
	post.Tags, err = RetrieveTags(post.ID)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// RetrieveNumberOfPosts ...
func RetrieveNumberOfPosts() (int64, error) {
	var count int64
	// Retrieve number of posts
	row := readDB.QueryRow(stmtRetrievePostsCount)
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// RetrieveNumberOfPostsByUser ...
func RetrieveNumberOfPostsByUser(userID int64) (int64, error) {
	var count int64
	// Retrieve number of posts
	row := readDB.QueryRow(stmtRetrievePostsCountByUser, userID)
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// RetrieveNumberOfPostsByTag ...
func RetrieveNumberOfPostsByTag(tagID int64) (int64, error) {
	var count int64
	// Retrieve number of posts
	row := readDB.QueryRow(stmtRetrievePostsCountByTag, tagID)
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func retrievePostCreationDateByID(postID int64) (*time.Time, error) {
	var creationDate time.Time
	// Retrieve number of posts
	row := readDB.QueryRow(stmtRetrievePostCreationDateByID, postID)
	err := row.Scan(&creationDate)
	if err != nil {
		return &creationDate, err
	}
	return &creationDate, nil
}

// RetrieveUser ...
func RetrieveUser(id int64) (*structure.User, error) {
	user := structure.User{}
	// Retrieve user
	row := readDB.QueryRow(stmtRetrieveUserByID, id)
	err := row.Scan(&user.ID, &user.Name, &user.Slug, &user.Email, &user.Image, &user.Cover, &user.Bio, &user.Website, &user.Location)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// RetrieveUserBySlug ...
func RetrieveUserBySlug(slug string) (*structure.User, error) {
	user := structure.User{}
	// Retrieve user
	row := readDB.QueryRow(stmtRetrieveUserBySlug, slug)
	err := row.Scan(&user.ID, &user.Name, &user.Slug, &user.Email, &user.Image, &user.Cover, &user.Bio, &user.Website, &user.Location)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// RetrieveUserByName ...
func RetrieveUserByName(name []byte) (*structure.User, error) {
	user := structure.User{}
	// Retrieve user
	row := readDB.QueryRow(stmtRetrieveUserByName, name)
	err := row.Scan(&user.ID, &user.Name, &user.Slug, &user.Email, &user.Image, &user.Cover, &user.Bio, &user.Website, &user.Location)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// RetrieveTags ...
func RetrieveTags(postID int64) ([]structure.Tag, error) {
	tags := make([]structure.Tag, 0)
	// Retrieve tags
	rows, err := readDB.Query(stmtRetrieveTags, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var tagID int64
		err := rows.Scan(&tagID)
		if err != nil {
			return nil, err
		}
		tag, err := RetrieveTag(tagID)
		// TODO: Error while receiving individual tag is ignored right now. Keep it this way?
		if err == nil {
			tags = append(tags, *tag)
		}
	}
	return tags, nil
}

// RetrieveTag ...
func RetrieveTag(tagID int64) (*structure.Tag, error) {
	tag := structure.Tag{}
	// Retrieve tag
	row := readDB.QueryRow(stmtRetrieveTagByID, tagID)
	err := row.Scan(&tag.ID, &tag.Name, &tag.Slug)
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// RetrieveTagBySlug ...
func RetrieveTagBySlug(slug string) (*structure.Tag, error) {
	tag := structure.Tag{}
	// Retrieve tag
	row := readDB.QueryRow(stmtRetrieveTagBySlug, slug)
	err := row.Scan(&tag.ID, &tag.Name, &tag.Slug)
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// RetrieveTagIDBySlug ...
func RetrieveTagIDBySlug(slug string) (int64, error) {
	var id int64
	row := readDB.QueryRow(stmtRetrieveTagIDBySlug, slug)
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// RetrieveHashedPasswordForUser ...
func RetrieveHashedPasswordForUser(name []byte) ([]byte, error) {
	var hashedPassword []byte
	row := readDB.QueryRow(stmtRetrieveHashedPasswordByName, name)
	err := row.Scan(&hashedPassword)
	if err != nil {
		return []byte{}, err
	}
	return hashedPassword, nil
}

// RetrieveBlog retrieves a blog entry from the db
func RetrieveBlog() (*structure.Blog, error) {
	tempBlog := structure.Blog{}
	// Title
	row := readDB.QueryRow(stmtRetrieveBlog, "title")
	err := row.Scan(&tempBlog.Title)
	if err != nil {
		return &tempBlog, err
	}
	// Description
	row = readDB.QueryRow(stmtRetrieveBlog, "description")
	err = row.Scan(&tempBlog.Description)
	if err != nil {
		return &tempBlog, err
	}
	// Logo
	row = readDB.QueryRow(stmtRetrieveBlog, "logo")
	err = row.Scan(&tempBlog.Logo)
	if err != nil {
		return &tempBlog, err
	}
	// Cover
	row = readDB.QueryRow(stmtRetrieveBlog, "cover")
	err = row.Scan(&tempBlog.Cover)
	if err != nil {
		return &tempBlog, err
	}
	// PostsPerPage
	row = readDB.QueryRow(stmtRetrieveBlog, "postsPerPage")
	err = row.Scan(&tempBlog.PostsPerPage)
	if err != nil {
		return &tempBlog, err
	}
	// ActiveTheme
	row = readDB.QueryRow(stmtRetrieveBlog, "activeTheme")
	err = row.Scan(&tempBlog.ActiveTheme)
	if err != nil {
		return &tempBlog, err
	}
	// Post count
	postCount, err := RetrieveNumberOfPosts()
	if err != nil {
		return &tempBlog, err
	}
	tempBlog.PostCount = postCount
	// Navigation
	var navigation []byte
	row = readDB.QueryRow(stmtRetrieveBlog, "navigation")
	err = row.Scan(&navigation)
	if err != nil {
		return &tempBlog, err
	}
	tempBlog.NavigationItems, err = makeNavigation(navigation)
	if err != nil {
		return &tempBlog, err
	}
	return &tempBlog, err
}

// RetrieveActiveTheme retrieves the active theme from the db
func RetrieveActiveTheme() (*string, error) {
	var activeTheme string
	row := readDB.QueryRow(stmtRetrieveBlog, "activeTheme")
	err := row.Scan(&activeTheme)
	if err != nil {
		return &activeTheme, err
	}
	return &activeTheme, nil
}

// RetrieveUsersCount reads the number of users in the db
func RetrieveUsersCount() int {
	userCount := -1
	row := readDB.QueryRow(stmtRetrieveUsersCount)
	err := row.Scan(&userCount)
	if err != nil {
		return -1
	}
	return userCount
}

func makeNavigation(navigation []byte) ([]structure.Navigation, error) {
	navigationItems := make([]structure.Navigation, 0)
	err := json.Unmarshal(navigation, &navigationItems)
	if err != nil {
		return navigationItems, err
	}
	return navigationItems, nil
}
