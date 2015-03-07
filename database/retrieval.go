package database

import (
	"database/sql"
	"github.com/kabukky/journey/structure"
)

const stmtRetrievePostsCount = "SELECT count(*) FROM posts WHERE page = 0 AND status = 'published'"
const stmtRetrievePostsCountByAuthor = "SELECT count(*) FROM posts WHERE page = 0 AND status = 'published' AND author_id = ?"
const stmtRetrievePostsCountByTag = "SELECT count(*) FROM posts, posts_tags WHERE posts_tags.post_id = posts.id AND posts_tags.tag_id = ? AND page = 0 AND status = 'published'"
const stmtRetrievePostsForIndex = "SELECT id, uuid, title, slug, markdown, html, featured, page, status, image, author_id, created_at FROM posts WHERE page = 0 AND status = 'published' ORDER BY id DESC LIMIT ? OFFSET ?"
const stmtRetrievePostsForApi = "SELECT id, uuid, title, slug, markdown, html, featured, page, status, image, author_id, created_at FROM posts ORDER BY id DESC LIMIT ? OFFSET ?"
const stmtRetrievePostsByAuthor = "SELECT id, uuid, title, slug, markdown, html, featured, page, status, image, author_id, created_at FROM posts WHERE page = 0 AND status = 'published' AND author_id = ? ORDER BY id DESC LIMIT ? OFFSET ?"
const stmtRetrievePostsByTag = "SELECT posts.id, posts.uuid, posts.title, posts.slug, posts.markdown, posts.html, posts.featured, posts.page, posts.status, posts.image, posts.author_id, posts.created_at FROM posts, posts_tags WHERE posts_tags.post_id = posts.id AND posts_tags.tag_id = ? AND page = 0 AND status = 'published' ORDER BY posts.id DESC LIMIT ? OFFSET ?"
const stmtRetrievePostById = "SELECT id, uuid, title, slug, markdown, html, featured, page, status, image, author_id, created_at FROM posts WHERE id = ?"
const stmtRetrievePostBySlug = "SELECT id, uuid, title, slug, markdown, html, featured, page, status, image, author_id, created_at FROM posts WHERE slug = ?"
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

func RetrievePostsByAuthor(author_id int64, limit int64, offset int64) ([]structure.Post, error) {
	// Retrieve posts
	rows, err := readDB.Query(stmtRetrievePostsByAuthor, author_id, limit, offset)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	posts, err := extractPosts(rows)
	if err != nil {
		return nil, err
	}
	return *posts, nil
}

func RetrievePostsByTag(tag_id int64, limit int64, offset int64) ([]structure.Post, error) {
	// Retrieve posts
	rows, err := readDB.Query(stmtRetrievePostsByTag, tag_id, limit, offset)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	posts, err := extractPosts(rows)
	if err != nil {
		return nil, err
	}
	return *posts, nil
}

func RetrievePostsForIndex(limit int64, offset int64) ([]structure.Post, error) {
	// Retrieve posts
	rows, err := readDB.Query(stmtRetrievePostsForIndex, limit, offset)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	posts, err := extractPosts(rows)
	if err != nil {
		return nil, err
	}
	return *posts, nil
}

func RetrievePostsForApi(limit int64, offset int64) ([]structure.Post, error) {
	// Retrieve posts
	rows, err := readDB.Query(stmtRetrievePostsForApi, limit, offset)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
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
		var authorId int64
		var status string
		err := rows.Scan(&post.Id, &post.Uuid, &post.Title, &post.Slug, &post.Markdown, &post.Html, &post.IsFeatured, &post.IsPage, &status, &post.Image, &authorId, &post.Date)
		if err != nil {
			return nil, err
		}
		// Evaluate status
		if status == "published" {
			post.IsPublished = true
		} else {
			post.IsPublished = false
		}
		// Retrieve author
		post.Author, err = RetrieveAuthor(authorId)
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
	var authorId int64
	var status string
	err := row.Scan(&post.Id, &post.Uuid, &post.Title, &post.Slug, &post.Markdown, &post.Html, &post.IsFeatured, &post.IsPage, &status, &post.Image, &authorId, &post.Date)
	if err != nil {
		return nil, err
	}
	// Evaluate status
	if status == "published" {
		post.IsPublished = true
	} else {
		post.IsPublished = false
	}
	// Retrieve author
	post.Author, err = RetrieveAuthor(authorId)
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

func RetrieveNumberOfPostsByAuthor(author_id int64) (int64, error) {
	var count int64
	// Retrieve number of posts
	row := readDB.QueryRow(stmtRetrievePostsCountByAuthor, author_id)
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

func RetrieveAuthor(id int64) (*structure.Author, error) {
	author := structure.Author{}
	// Retrieve author
	row := readDB.QueryRow(stmtRetrieveUserById, id)
	err := row.Scan(&author.Id, &author.Name, &author.Slug, &author.Email, &author.Image, &author.Cover, &author.Bio, &author.Website, &author.Location)
	if err != nil {
		return nil, err
	}
	return &author, nil
}

func RetrieveAuthorBySlug(slug string) (*structure.Author, error) {
	author := structure.Author{}
	// Retrieve author
	row := readDB.QueryRow(stmtRetrieveUserBySlug, slug)
	err := row.Scan(&author.Id, &author.Name, &author.Slug, &author.Email, &author.Image, &author.Cover, &author.Bio, &author.Website, &author.Location)
	if err != nil {
		return nil, err
	}
	return &author, nil
}

func RetrieveAuthorByName(name []byte) (*structure.Author, error) {
	author := structure.Author{}
	// Retrieve author
	row := readDB.QueryRow(stmtRetrieveUserByName, name)
	err := row.Scan(&author.Id, &author.Name, &author.Slug, &author.Email, &author.Image, &author.Cover, &author.Bio, &author.Website, &author.Location)
	if err != nil {
		return nil, err
	}
	return &author, nil
}

func RetrieveTags(postId int64) ([]structure.Tag, error) {
	tags := make([]structure.Tag, 0)
	// Retrieve tags
	rows, err := readDB.Query(stmtRetrieveTags, postId)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
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
	postCount, err := RetrieveNumberOfPosts()
	if err != nil {
		return &tempBlog, err
	}
	tempBlog.PostCount = postCount
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
