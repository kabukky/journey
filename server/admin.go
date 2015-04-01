package server

import (
	"encoding/json"
	"github.com/kabukky/journey/authentication"
	"github.com/kabukky/journey/configuration"
	"github.com/kabukky/journey/conversion"
	"github.com/kabukky/journey/database"
	"github.com/kabukky/journey/filenames"
	"github.com/kabukky/journey/slug"
	"github.com/kabukky/journey/structure"
	"github.com/kabukky/journey/structure/methods"
	"github.com/kabukky/journey/templates"
	"github.com/twinj/uuid"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type JsonPost struct {
	Id          int64
	Title       string
	Slug        string
	Markdown    string
	Html        string
	IsFeatured  bool
	IsPage      bool
	IsPublished bool
	Image       string
	Date        time.Time
	Tags        string
}

type JsonBlog struct {
	Title        string
	Description  string
	Logo         string
	Cover        string
	Themes       []string
	ActiveTheme  string
	PostsPerPage int64
}

type JsonUser struct {
	Id               int64
	Name             string
	Email            string
	Image            string
	Cover            string
	Bio              string
	Website          string
	Location         string
	Password         string
	PasswordRepeated string
}

type JsonUserId struct {
	Id int64
}

var validAdminPath = regexp.MustCompile("^/admin/([\\.a-zA-Z0-9#-]*)$")
var validAdminApiPath = regexp.MustCompile("^/admin/api/(posts|post|upload|images|blog|userid|user)/?(\\d+)?/?$")

func loginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		http.ServeFile(w, r, filepath.Join(filenames.AdminFilepath, "login.html"))
		return
	case "POST":
		name := r.FormValue("name")
		password := r.FormValue("password")
		if name != "" && password != "" {
			if authentication.LoginIsCorrect(name, password) {
				authentication.SetSession(name, w)
			} else {
				log.Println("Failed login attempt for user " + name)
			}
		}
		http.Redirect(w, r, "/admin/", 302)
		return
	}
}

func registrationHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if database.RetrieveUsersCount() == 0 {
			http.ServeFile(w, r, filepath.Join(filenames.AdminFilepath, "registration.html"))
			return
		} else {
			http.Redirect(w, r, "/admin/", 302)
			return
		}
	case "POST":
		if database.RetrieveUsersCount() == 0 { // TODO: Or check if authenticated user is admin when adding users from inside the admin area
			name := r.FormValue("name")
			email := r.FormValue("email")
			password := r.FormValue("password")
			if name != "" && password != "" {
				hashedPassword, err := authentication.EncryptPassword(password)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				user := structure.User{Name: []byte(name), Slug: slug.Generate(name, "users"), Email: []byte(email), Image: []byte(filenames.DefaultUserImageFilename), Cover: []byte(filenames.DefaultUserCoverFilename), Role: 4}
				err = methods.SaveUser(&user, hashedPassword, 1)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				http.Redirect(w, r, "/admin/", 302)
				return
			}
			http.Redirect(w, r, "/admin/", 302)
			return
		} else {
			// TODO: Handle creation of other users (not just the first one)
			http.Error(w, "Not implemented yet.", http.StatusInternalServerError)
			return
		}
	}
}

// Not used at the moment.
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	authentication.ClearSession(w)
	http.Redirect(w, r, "/admin/", 302)
	return
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	if database.RetrieveUsersCount() == 0 {
		http.Redirect(w, r, "/admin/register/", 302)
		return
	} else {
		userName := authentication.GetUserName(r)
		if userName != "" {
			// Get arguments (files)
			m := validAdminPath.FindStringSubmatch(r.URL.Path)
			if m == nil {
				http.Redirect(w, r, "/admin/", http.StatusFound)
				return
			} else if m[1] == "" {
				http.ServeFile(w, r, filepath.Join(filenames.AdminFilepath, "admin.html"))
				return
			} else {
				http.ServeFile(w, r, filepath.Join(filenames.AdminFilepath, m[1]))
				return
			}
		} else {
			http.Redirect(w, r, "/admin/login/", 302)
			return
		}
	}
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	if database.RetrieveUsersCount() == 0 { // If there are no users in db (e.g. a brand new installation), don't allow use of api
		http.Error(w, "No users in database!", http.StatusInternalServerError)
		return
	} else {
		userName := authentication.GetUserName(r)
		if userName != "" {
			author, err := database.RetrieveUserByName([]byte(userName))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			userId := author.Id
			// Get arguments (api call)
			m := validAdminApiPath.FindStringSubmatch(r.URL.Path)
			if m == nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			switch m[1] {
			case "":
				http.Error(w, "Not a valid api function!", http.StatusInternalServerError)
				return
			// API function to get all posts by pages
			case "posts":
				page, err := strconv.Atoi(m[2])
				if err != nil || page < 1 {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				postsPerPage := int64(15)
				posts, err := database.RetrievePostsForApi(postsPerPage, ((int64(page) - 1) * postsPerPage))
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				json, err := json.Marshal(postsToJson(posts))
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write(json)
			// API function to get a specific post using its id
			case "post":
				switch r.Method {
				// Get post
				case "GET":
					id, err := strconv.ParseInt(m[2], 10, 64)
					if err != nil || id < 1 {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					post, err := database.RetrievePostById(id)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					json, err := json.Marshal(postToJson(post))
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					w.Header().Set("Content-Type", "application/json")
					w.Write(json)
				// Create post
				case "POST":
					decoder := json.NewDecoder(r.Body)
					var json JsonPost
					err := decoder.Decode(&json)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					var postSlug string
					if json.Slug != "" { // Ceck if user has submitted a custom slug
						postSlug = slug.Generate(json.Slug, "posts")
					} else {
						postSlug = slug.Generate(json.Title, "posts")
					}
					post := structure.Post{Title: []byte(json.Title), Slug: postSlug, Markdown: []byte(json.Markdown), Html: conversion.GenerateHtmlFromMarkdown([]byte(json.Markdown)), IsFeatured: json.IsFeatured, IsPage: json.IsPage, IsPublished: json.IsPublished, Image: []byte(json.Image), Date: time.Now(), Tags: methods.GenerateTagsFromCommaString(json.Tags), Author: &structure.User{Id: userId}}
					err = methods.SavePost(&post)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("Post created!"))
					return
				// Update post
				case "PATCH":
					decoder := json.NewDecoder(r.Body)
					var json JsonPost
					err := decoder.Decode(&json)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					var postSlug string
					// Get current slug of post
					post, err := database.RetrievePostById(json.Id)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					if json.Slug != post.Slug { // Check if user has submitted a custom slug
						postSlug = slug.Generate(json.Slug, "posts")
					} else {
						postSlug = post.Slug
					}
					*post = structure.Post{Id: json.Id, Title: []byte(json.Title), Slug: postSlug, Markdown: []byte(json.Markdown), Html: conversion.GenerateHtmlFromMarkdown([]byte(json.Markdown)), IsFeatured: json.IsFeatured, IsPage: json.IsPage, IsPublished: json.IsPublished, Image: []byte(json.Image), Date: time.Now(), Tags: methods.GenerateTagsFromCommaString(json.Tags), Author: &structure.User{Id: userId}}
					err = methods.UpdatePost(post)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("Post updated!"))
					return
				// Delete post
				case "DELETE":
					id, err := strconv.ParseInt(m[2], 10, 64)
					if err != nil || id < 1 {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					err = database.DeletePostById(id)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("Post deleted!"))
					return
				}
			// API function to upload images
			case "upload":
				// Create multipart reader
				reader, err := r.MultipartReader()
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				// Slice to hold all paths to the files
				allFilePaths := make([]string, 0)
				// Copy each part to destination.
				for {
					part, err := reader.NextPart()
					if err == io.EOF {
						break
					}
					// If part.FileName() is empty, skip this iteration.
					if part.FileName() == "" {
						continue
					}
					// Folder structure: year/month/randomname
					filePath := filepath.Join(filenames.ImagesFilepath, time.Now().Format("2006"), time.Now().Format("01"))
					if os.MkdirAll(filePath, 0777) != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					dst, err := os.Create(filepath.Join(filePath, strconv.FormatInt(time.Now().Unix(), 10)+"_"+uuid.Formatter(uuid.NewV4(), uuid.Clean)+filepath.Ext(part.FileName())))
					defer dst.Close()
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					if _, err := io.Copy(dst, part); err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					allFilePaths = append(allFilePaths, strings.Replace(dst.Name(), filenames.ContentFilepath, "", 1))
				}
				json, err := json.Marshal(allFilePaths)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write(json)
				return
			// API function to get all images by pages
			case "images":
				page, err := strconv.Atoi(m[2])
				if err != nil || page < 1 {
					http.Error(w, "Not a valid api function!", http.StatusInternalServerError)
					return
				}
				images := make([]string, 0)
				// Walk all files in images folder
				err = filepath.Walk(filenames.ImagesFilepath, func(filePath string, info os.FileInfo, err error) error {
					if !info.IsDir() && (strings.EqualFold(filepath.Ext(filePath), ".jpg") || strings.EqualFold(filepath.Ext(filePath), ".jpeg") || strings.EqualFold(filepath.Ext(filePath), ".gif") || strings.EqualFold(filepath.Ext(filePath), ".png") || strings.EqualFold(filepath.Ext(filePath), ".svg")) {
						// Rewrite to file path on server
						filePath = strings.Replace(filePath, filenames.ImagesFilepath, "/images", 1)
						// Prepend file to slice (thus reversing the order)
						images = append([]string{filePath}, images...)
					}
					return nil
				})
				if len(images) == 0 {
					// Write empty json array
					w.Header().Set("Content-Type", "application/json")
					w.Write([]byte("[]"))
					return
				}
				imagesPerPage := 15
				start := (page * imagesPerPage) - imagesPerPage
				end := page * imagesPerPage
				if start > (len(images) - 1) {
					// Write empty json array
					w.Header().Set("Content-Type", "application/json")
					w.Write([]byte("[]"))
					return
				}
				if end > len(images) {
					end = len(images)
				}
				json, err := json.Marshal(images[start:end])
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write(json)
				return
			// API function to get and set blog settings
			case "blog":
				switch r.Method {
				case "GET":
					blog, err := database.RetrieveBlog()
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					blogJson := JsonBlog{Title: string(blog.Title), Description: string(blog.Description), Logo: string(blog.Logo), Cover: string(blog.Cover), PostsPerPage: blog.PostsPerPage, Themes: templates.GetAllThemes(), ActiveTheme: blog.ActiveTheme}
					json, err := json.Marshal(blogJson)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					w.Header().Set("Content-Type", "application/json")
					w.Write(json)
				case "PATCH":
					decoder := json.NewDecoder(r.Body)
					var json JsonBlog
					err := decoder.Decode(&json)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					// Make sure postPerPage is over 0
					if json.PostsPerPage < 1 {
						json.PostsPerPage = 1
					}
					// Retrieve old post settings for comparison
					blog, err := database.RetrieveBlog()
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					tempBlog := structure.Blog{Url: []byte(configuration.Config.Url), Title: []byte(json.Title), Description: []byte(json.Description), Logo: []byte(json.Logo), Cover: []byte(json.Cover), AssetPath: []byte("/assets/"), PostCount: blog.PostCount, PostsPerPage: json.PostsPerPage, ActiveTheme: json.ActiveTheme}
					err = methods.UpdateBlog(&tempBlog)
					// Check if active theme setting has been changed, if so, generate templates from new theme
					if tempBlog.ActiveTheme != blog.ActiveTheme {
						err = templates.Generate()
						if err != nil {
							// If there's an error while generating the new templates, the whole program must be stopped.
							log.Fatal("Fatal error: Template data couldn't be generated from theme files: " + err.Error())
							http.Error(w, err.Error(), http.StatusInternalServerError)
							return
						}
					}
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("Blog settings updated!"))
					return
				}
			// API function to get and set user (author) settings
			case "user":
				switch r.Method {
				case "GET":
					id, err := strconv.ParseInt(m[2], 10, 64)
					if err != nil || id < 1 {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					} else if id != userId { // Make sure the authenticating user is only accessing his own data. TODO: Make sure the user is admin or multiple users have been introduced
						http.Error(w, "You don't have permission to access this data.", http.StatusForbidden)
						return
					}
					author, err := database.RetrieveUser(id)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					authorJson := JsonUser{Id: author.Id, Name: string(author.Name), Email: string(author.Email), Image: string(author.Image), Cover: string(author.Cover), Bio: string(author.Bio), Website: string(author.Website), Location: string(author.Location)}
					json, err := json.Marshal(authorJson)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					w.Header().Set("Content-Type", "application/json")
					w.Write(json)
				case "PATCH":
					decoder := json.NewDecoder(r.Body)
					var json JsonUser
					err := decoder.Decode(&json)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					// Make sure user id is over 0 and E-Mail is included.
					if json.Id < 1 {
						http.Error(w, "Wrong user id.", http.StatusInternalServerError)
						return
					} else if json.Email == "" {
						http.Error(w, "Email needs to be included.", http.StatusInternalServerError)
						return
					}
					author := structure.User{Id: json.Id, Email: []byte(json.Email), Image: []byte(json.Image), Cover: []byte(json.Cover), Bio: []byte(json.Bio), Website: []byte(json.Website), Location: []byte(json.Location)}
					err = methods.UpdateUser(&author, userId)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					if json.Password != "" && (json.Password == json.PasswordRepeated) { // Update password if a new one was submitted
						encryptedPassword, err := authentication.EncryptPassword(json.Password)
						if err != nil {
							http.Error(w, err.Error(), http.StatusInternalServerError)
							return
						}
						err = database.UpdateUserPassword(author.Id, encryptedPassword, time.Now(), json.Id)
						if err != nil {
							http.Error(w, err.Error(), http.StatusInternalServerError)
							return
						}
					}
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("User settings updated!"))
					return
				}
			// API function to get the id of the authenticated user
			case "userid":
				jsonUserId := JsonUserId{Id: userId}
				json, err := json.Marshal(jsonUserId)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write(json)
			}
		} else {
			http.Error(w, "Not logged in!", http.StatusInternalServerError)
			return
		}
	}
}

func postsToJson(posts []structure.Post) *[]JsonPost {
	jsonPosts := make([]JsonPost, len(posts))
	for index, _ := range posts {
		jsonPosts[index] = *postToJson(&posts[index])
	}
	return &jsonPosts
}

func postToJson(post *structure.Post) *JsonPost {
	var jsonPost JsonPost
	jsonPost.Id = post.Id
	jsonPost.Title = string(post.Title)
	jsonPost.Slug = post.Slug
	jsonPost.Markdown = string(post.Markdown)
	jsonPost.Html = string(post.Html)
	jsonPost.IsFeatured = post.IsFeatured
	jsonPost.IsPage = post.IsPage
	jsonPost.IsPublished = post.IsPublished
	jsonPost.Image = string(post.Image)
	jsonPost.Date = post.Date
	tags := make([]string, len(post.Tags))
	for index, _ := range post.Tags {
		tags[index] = string(post.Tags[index].Name)
	}
	jsonPost.Tags = strings.Join(tags, ",")
	return &jsonPost
}

func InitializeAdmin(mux *http.ServeMux) {
	mux.Handle("/admin/", http.HandlerFunc(adminHandler))
	mux.Handle("/admin/api/", http.HandlerFunc(apiHandler))
	mux.Handle("/admin/login/", http.HandlerFunc(loginHandler))
	mux.Handle("/admin/logout/", http.HandlerFunc(logoutHandler))
	mux.Handle("/admin/register/", http.HandlerFunc(registrationHandler))
}
