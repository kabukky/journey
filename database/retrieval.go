package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/kabukky/journey/structure"
)

const stmtRetrievePostsCount = "SELECT count(*) FROM posts WHERE page = 0 AND status = 'published'"
const stmtRetrievePostsCountByUser = "SELECT count(*) FROM posts WHERE page = 0 AND status = 'published' AND author_id = ?"
const stmtRetrievePostsCountByTag = "SELECT count(*) FROM posts, posts_tags WHERE posts_tags.post_id = posts.id AND posts_tags.tag_id = ? AND page = 0 AND status = 'published'"
const stmtRetrievePostsForIndex = "SELECT id, uuid, title, slug, markdown, html, featured, page, status, meta_description, image, author_id, published_at FROM posts WHERE page = 0 AND status = 'published' ORDER BY published_at DESC LIMIT ? OFFSET ?"
const stmtRetrievePostsForApi = "SELECT id, uuid, title, slug, markdown, html, featured, page, status, meta_description, image, author_id, published_at FROM posts ORDER BY id DESC LIMIT ? OFFSET ?"
const stmtRetrievePostsByUser = "SELECT id, uuid, title, slug, markdown, html, featured, page, status, meta_description, image, author_id, published_at FROM posts WHERE page = 0 AND status = 'published' AND author_id = ? ORDER BY published_at DESC LIMIT ? OFFSET ?"
const stmtRetrievePostsByTag = "SELECT posts.id, posts.uuid, posts.title, posts.slug, posts.markdown, posts.html, posts.featured, posts.page, posts.status, posts.meta_description, posts.image, posts.author_id, posts.published_at FROM posts, posts_tags WHERE posts_tags.post_id = posts.id AND posts_tags.tag_id = ? AND page = 0 AND status = 'published' ORDER BY posts.published_at DESC LIMIT ? OFFSET ?"
const stmtRetrievePostById = "SELECT id, uuid, title, slug, markdown, html, featured, page, status, meta_description, image, author_id, published_at FROM posts WHERE id = ?"
const stmtRetrievePostBySlug = "SELECT id, uuid, title, slug, markdown, html, featured, page, status, meta_description, image, author_id, published_at FROM posts WHERE slug = ? COLLATE NOCASE"
const stmtRetrieveUserById = "SELECT id, name, slug, email, image, cover, bio, website, location FROM users WHERE id = ?"
const stmtRetrieveUserBySlug = "SELECT id, name, slug, email, image, cover, bio, website, location FROM users WHERE slug = ?"
const stmtRetrieveUserByName = "SELECT id, name, slug, email, image, cover, bio, website, location FROM users WHERE name = ?"
const stmtRetrieveTags = "SELECT tag_id FROM posts_tags WHERE post_id = ?"
const stmtRetrieveTagById = "SELECT id, name, slug FROM tags WHERE id = ?"
const stmtRetrieveTagBySlug = "SELECT id, name, slug FROM tags WHERE slug = ?"
const stmtRetrieveTagIdBySlug = "SELECT id FROM tags WHERE slug = ?"
const stmtRetrieveHashedPasswordByName = "SELECT password FROM users WHERE name = ?"
const stmtRetrieveUsersCount = "SELECT count(*) FROM users"
const stmtRetrieveBlog = "SELECT value FROM settings WHERE key = ?"
const stmtRetrievePostCreationDateById = "SELECT created_at FROM posts WHERE id = ?"
const stmtRetrieveSitemap = "SELECT slug, updated_at FROM posts WHERE status = 'published' AND slug != 404 ORDER BY updated_at DESC"

func RetrievePostById(id int64) (*structure.Post, error) {
	// Retrieve post
	row := readDB.QueryRow(stmtRetrievePostById, id)
	return extractPost(row)
}

func RetrievePostBySlug(slug string) (*structure.Post, error) {
	// Retrieve post
	row := readDB.QueryRow(stmtRetrievePostBySlug, slug)
	return extractPost(row)
}

func RetrievePostsByUser(user_id int64, limit int64, offset int64) ([]structure.Post, error) {
	// Retrieve posts
	rows, err := readDB.Query(stmtRetrievePostsByUser, user_id, limit, offset)
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

func RetrievePostsByTag(tag_id int64, limit int64, offset int64) ([]structure.Post, error) {
	// Retrieve posts
	rows, err := readDB.Query(stmtRetrievePostsByTag, tag_id, limit, offset)
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

func RetrievePostsForApi(limit int64, offset int64) ([]structure.Post, error) {
	// Retrieve posts
	rows, err := readDB.Query(stmtRetrievePostsForApi, limit, offset)
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
		var userId int64
		var status string
		err := rows.Scan(&post.Id, &post.Uuid, &post.Title, &post.Slug, &post.Markdown, &post.Html, &post.IsFeatured, &post.IsPage, &status, &post.MetaDescription, &post.Image, &userId, &post.Date)
		if err != nil {
			return nil, err
		}
		// If there was no publication date attached to the post, make its creation date the date of the post
		if post.Date == nil {
			post.Date, err = retrievePostCreationDateById(post.Id)
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
		post.Author, err = RetrieveUser(userId)
		if err != nil {
			return nil, err
		}
		// Retrieve tags
		post.Tags, err = RetrieveTags(post.Id)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return &posts, nil
}

func extractPost(row *sql.Row) (*structure.Post, error) {
	post := structure.Post{}
	var userId int64
	var status string
	err := row.Scan(&post.Id, &post.Uuid, &post.Title, &post.Slug, &post.Markdown, &post.Html, &post.IsFeatured, &post.IsPage, &status, &post.MetaDescription, &post.Image, &userId, &post.Date)
	if err != nil {
		return nil, err
	}
	// If there was no publication date attached to the post, make its creation date the date of the post
	if post.Date == nil {
		post.Date, err = retrievePostCreationDateById(post.Id)
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
	post.Author, err = RetrieveUser(userId)
	if err != nil {
		return nil, err
	}
	// Retrieve tags
	post.Tags, err = RetrieveTags(post.Id)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

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

func RetrieveNumberOfPostsByUser(user_id int64) (int64, error) {
	var count int64
	// Retrieve number of posts
	row := readDB.QueryRow(stmtRetrievePostsCountByUser, user_id)
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func RetrieveNumberOfPostsByTag(tag_id int64) (int64, error) {
	var count int64
	// Retrieve number of posts
	row := readDB.QueryRow(stmtRetrievePostsCountByTag, tag_id)
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func retrievePostCreationDateById(post_id int64) (*time.Time, error) {
	var creationDate time.Time
	// Retrieve number of posts
	row := readDB.QueryRow(stmtRetrievePostCreationDateById, post_id)
	err := row.Scan(&creationDate)
	if err != nil {
		return &creationDate, err
	}
	return &creationDate, nil
}

func RetrieveUser(id int64) (*structure.User, error) {
	user := structure.User{}
	// Retrieve user
	row := readDB.QueryRow(stmtRetrieveUserById, id)
	err := row.Scan(&user.Id, &user.Name, &user.Slug, &user.Email, &user.Image, &user.Cover, &user.Bio, &user.Website, &user.Location)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func RetrieveUserBySlug(slug string) (*structure.User, error) {
	user := structure.User{}
	// Retrieve user
	row := readDB.QueryRow(stmtRetrieveUserBySlug, slug)
	err := row.Scan(&user.Id, &user.Name, &user.Slug, &user.Email, &user.Image, &user.Cover, &user.Bio, &user.Website, &user.Location)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func RetrieveUserByName(name []byte) (*structure.User, error) {
	user := structure.User{}
	// Retrieve user
	row := readDB.QueryRow(stmtRetrieveUserByName, name)
	err := row.Scan(&user.Id, &user.Name, &user.Slug, &user.Email, &user.Image, &user.Cover, &user.Bio, &user.Website, &user.Location)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func RetrieveTags(postId int64) ([]structure.Tag, error) {
	tags := make([]structure.Tag, 0)
	// Retrieve tags
	rows, err := readDB.Query(stmtRetrieveTags, postId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var tagId int64
		err := rows.Scan(&tagId)
		if err != nil {
			return nil, err
		}
		tag, err := RetrieveTag(tagId)
		// TODO: Error while receiving individual tag is ignored right now. Keep it this way?
		if err == nil {
			tags = append(tags, *tag)
		}
	}
	return tags, nil
}

func RetrieveTag(tagId int64) (*structure.Tag, error) {
	tag := structure.Tag{}
	// Retrieve tag
	row := readDB.QueryRow(stmtRetrieveTagById, tagId)
	err := row.Scan(&tag.Id, &tag.Name, &tag.Slug)
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func RetrieveTagBySlug(slug string) (*structure.Tag, error) {
	tag := structure.Tag{}
	// Retrieve tag
	row := readDB.QueryRow(stmtRetrieveTagBySlug, slug)
	err := row.Scan(&tag.Id, &tag.Name, &tag.Slug)
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func RetrieveTagIdBySlug(slug string) (int64, error) {
	var id int64
	row := readDB.QueryRow(stmtRetrieveTagIdBySlug, slug)
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func RetrieveHashedPasswordForUser(name []byte) ([]byte, error) {
	var hashedPassword []byte
	row := readDB.QueryRow(stmtRetrieveHashedPasswordByName, name)
	err := row.Scan(&hashedPassword)
	if err != nil {
		return []byte{}, err
	}
	return hashedPassword, nil
}

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
	// GhostHead
	row = readDB.QueryRow(stmtRetrieveBlog, "ghost_head")
	err = row.Scan(&tempBlog.GhostHead)
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

func RetrieveActiveTheme() (*string, error) {
	var activeTheme string
	row := readDB.QueryRow(stmtRetrieveBlog, "activeTheme")
	err := row.Scan(&activeTheme)
	if err != nil {
		return &activeTheme, err
	}
	return &activeTheme, nil
}

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

func RetrieveSitemap() ([]structure.SmURL, error) {
	SmURL := make([]structure.SmURL, 0)
	// Retrieve sitemap URLs
	rows, err := readDB.Query(stmtRetrieveSitemap)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var url structure.SmURL
		var lastMod *time.Time

		err := rows.Scan(&url.Loc, &lastMod)
		if err != nil {
			return nil, err
		}
		if err == nil {
			url.LastMod = fmt.Sprint(lastMod.UTC().Format("2006-01-02T15:04:05-07:00"))
			SmURL = append(SmURL, url)
		}
	}
	return SmURL, nil
}
