package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/dimfeld/httptreemux"
	"github.com/satori/go.uuid"

	"github.com/kabukky/journey/authentication"
	"github.com/kabukky/journey/configuration"
	"github.com/kabukky/journey/conversion"
	"github.com/kabukky/journey/database"
	"github.com/kabukky/journey/date"
	"github.com/kabukky/journey/filenames"
	"github.com/kabukky/journey/helpers"
	"github.com/kabukky/journey/slug"
	"github.com/kabukky/journey/structure"
	"github.com/kabukky/journey/structure/methods"
	"github.com/kabukky/journey/templates"
)

type JsonPost struct {
	Id              int64
	Title           string
	Slug            string
	Markdown        string
	Html            string
	IsFeatured      bool
	IsPage          bool
	IsPublished     bool
	Image           string
	MetaDescription string
	Date            *time.Time
	Tags            string
}

type JsonBlog struct {
	Url             string
	Title           string
	Description     string
	Logo            string
	Cover           string
	Themes          []string
	ActiveTheme     string
	PostsPerPage    int64
	NavigationItems []structure.Navigation
}

type JsonUser struct {
	Id               int64
	Name             string
	Slug             string
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

type JsonImage struct {
	Filename string
}

// Function to serve the login page
func getLoginHandler(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	if database.RetrieveUsersCount() == 0 {
		http.Redirect(w, r, "/admin/register/", 302)
		return
	}
	http.ServeFile(w, r, filepath.Join(filenames.AdminFilepath, "login.html"))
}

// Function to receive a login form
func postLoginHandler(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	name := r.FormValue("name")
	password := r.FormValue("password")
	if name != "" && password != "" {
		if authentication.LoginIsCorrect(name, password) {
			logInUser(name, w)
		} else {
			log.Println("Failed login attempt for user " + name)
		}
	}
	http.Redirect(w, r, "/admin/", 302)
}

// Function to serve the registration form
func getRegistrationHandler(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	if database.RetrieveUsersCount() == 0 {
		http.ServeFile(w, r, filepath.Join(filenames.AdminFilepath, "registration.html"))
		return
	}
	http.Redirect(w, r, "/admin/", 302)
}

// Function to recieve a registration form.
func postRegistrationHandler(w http.ResponseWriter, r *http.Request, _ map[string]string) {
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
	}

	// TODO: Handle creation of other users (not just the first one)
	http.Error(w, "Not implemented yet.", http.StatusInternalServerError)
}

// Function to log out the user. Not used at the moment.
func logoutHandler(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	authentication.ClearSession(w)
	http.Redirect(w, r, "/admin/login/", 302)
}

// Function to route the /admin/ url accordingly. (Is user logged in? Is at least one user registered?)
func adminHandler(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	if database.RetrieveUsersCount() == 0 {
		http.Redirect(w, r, "/admin/register/", 302)
		return
	}
	userName := authentication.GetUserName(r)
	if userName != "" {
		http.ServeFile(w, r, filepath.Join(filenames.AdminFilepath, "admin.html"))
		return
	}
	http.Redirect(w, r, "/admin/login/", 302)
}

// Function to serve files belonging to the admin interface.
func adminFileHandler(w http.ResponseWriter, r *http.Request, params map[string]string) {
	userName := authentication.GetUserName(r)
	if userName != "" {
		path := filepath.Join(filenames.AdminFilepath, params["filepath"])
		if !helpers.FileExists(path) {
			http.Error(w, "Still lost?", http.StatusNotFound)
			return
		}
		// Get arguments (files)
		http.ServeFile(w, r, path)
		return
	}
	http.Error(w, "Still lost?", http.StatusNotFound)
}

// API function to get all posts by pages
func apiPostsHandler(w http.ResponseWriter, r *http.Request, params map[string]string) {
	userName := authentication.GetUserName(r)
	if userName != "" {
		number := params["number"]
		page, err := strconv.Atoi(number)
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
		return
	}
	http.Error(w, "Not logged in!", http.StatusInternalServerError)
}

// API function to get a post by id
func getApiPostHandler(w http.ResponseWriter, r *http.Request, params map[string]string) {
	userName := authentication.GetUserName(r)
	if userName != "" {
		id := params["id"]
		// Get post
		postId, err := strconv.ParseInt(id, 10, 64)
		if err != nil || postId < 1 {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		post, err := database.RetrievePostById(postId)
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
		return
	}
	http.Error(w, "Not logged in!", http.StatusInternalServerError)
}

// API function to create a post
func postApiPostHandler(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	userName := authentication.GetUserName(r)
	if userName != "" {
		userId, err := getUserId(userName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Create post
		decoder := json.NewDecoder(r.Body)
		var json JsonPost
		err = decoder.Decode(&json)
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
		currentTime := date.GetCurrentTime()
		post := structure.Post{Title: []byte(json.Title), Slug: postSlug, Markdown: []byte(json.Markdown), Html: conversion.GenerateHtmlFromMarkdown([]byte(json.Markdown)), IsFeatured: json.IsFeatured, IsPage: json.IsPage, IsPublished: json.IsPublished, MetaDescription: []byte(json.MetaDescription), Image: []byte(json.Image), Date: &currentTime, Tags: methods.GenerateTagsFromCommaString(json.Tags), Author: &structure.User{Id: userId}}
		err = methods.SavePost(&post)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Post created!"))
		return
	}
	http.Error(w, "Not logged in!", http.StatusInternalServerError)
}

// API function to update a post.
func patchApiPostHandler(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	userName := authentication.GetUserName(r)
	if userName != "" {
		userId, err := getUserId(userName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Update post
		decoder := json.NewDecoder(r.Body)
		var json JsonPost
		err = decoder.Decode(&json)
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
		currentTime := date.GetCurrentTime()
		*post = structure.Post{Id: json.Id, Title: []byte(json.Title), Slug: postSlug, Markdown: []byte(json.Markdown), Html: conversion.GenerateHtmlFromMarkdown([]byte(json.Markdown)), IsFeatured: json.IsFeatured, IsPage: json.IsPage, IsPublished: json.IsPublished, MetaDescription: []byte(json.MetaDescription), Image: []byte(json.Image), Date: &currentTime, Tags: methods.GenerateTagsFromCommaString(json.Tags), Author: &structure.User{Id: userId}}
		err = methods.UpdatePost(post)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Post updated!"))
		return
	}
	http.Error(w, "Not logged in!", http.StatusInternalServerError)
}

// API function to delete a post by id.
func deleteApiPostHandler(w http.ResponseWriter, r *http.Request, params map[string]string) {
	userName := authentication.GetUserName(r)
	if userName != "" {
		id := params["id"]
		// Delete post
		postId, err := strconv.ParseInt(id, 10, 64)
		if err != nil || postId < 1 {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = methods.DeletePost(postId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Post deleted!"))
		return
	}
	http.Error(w, "Not logged in!", http.StatusInternalServerError)
}

// API function to upload images
func apiUploadHandler(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	userName := authentication.GetUserName(r)
	if userName != "" {
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
			currentDate := date.GetCurrentTime()
			filePath := filepath.Join(filenames.ImagesFilepath, currentDate.Format("2006"), currentDate.Format("01"))
			if os.MkdirAll(filePath, 0777) != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			dst, err := os.Create(filepath.Join(filePath, strconv.FormatInt(currentDate.Unix(), 10)+"_"+uuid.NewV4().String()+filepath.Ext(part.FileName())))
			defer dst.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if _, err := io.Copy(dst, part); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// Rewrite to file path on server
			filePath = strings.Replace(dst.Name(), filenames.ImagesFilepath, "/images", 1)
			// Make sure to always use "/" as path separator (to make a valid url that we can use on the blog)
			filePath = filepath.ToSlash(filePath)
			allFilePaths = append(allFilePaths, filePath)
		}
		json, err := json.Marshal(allFilePaths)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
		return
	}
	http.Error(w, "Not logged in!", http.StatusInternalServerError)
}

// API function to get all images by pages
func apiImagesHandler(w http.ResponseWriter, r *http.Request, params map[string]string) {
	userName := authentication.GetUserName(r)
	if userName != "" {
		number := params["number"]
		page, err := strconv.Atoi(number)
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
				// Make sure to always use "/" as path separator (to make a valid url that we can use on the blog)
				filePath = filepath.ToSlash(filePath)
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
	}
	http.Error(w, "Not logged in!", http.StatusInternalServerError)
}

// API function to delete an image by its filename.
func deleteApiImageHandler(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	userName := authentication.GetUserName(r)
	if userName != "" { // TODO: Check if the user has permissions to delete the image
		// Get the file name from the json data
		decoder := json.NewDecoder(r.Body)
		var json JsonImage
		err := decoder.Decode(&json)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = filepath.Walk(filenames.ImagesFilepath, func(filePath string, info os.FileInfo, err error) error {
			if !info.IsDir() && filepath.Base(filePath) == filepath.Base(json.Filename) {
				err := os.Remove(filePath)
				if err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Image deleted!"))
		return
	}
	http.Error(w, "Not logged in!", http.StatusInternalServerError)
}

// API function to get blog settings
func getApiBlogHandler(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	userName := authentication.GetUserName(r)
	if userName != "" {
		// Read lock the global blog
		methods.Blog.RLock()
		defer methods.Blog.RUnlock()
		blogJson := blogToJson(methods.Blog)
		json, err := json.Marshal(blogJson)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
		return
	}
	http.Error(w, "Not logged in!", http.StatusInternalServerError)
}

// API function to update blog settings
func patchApiBlogHandler(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	userName := authentication.GetUserName(r)
	if userName != "" {
		userId, err := getUserId(userName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		decoder := json.NewDecoder(r.Body)
		var json JsonBlog
		err = decoder.Decode(&json)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Make sure postPerPage is over 0
		if json.PostsPerPage < 1 {
			json.PostsPerPage = 1
		}
		// Remove blog url in front of navigation urls
		for index := range json.NavigationItems {
			if strings.HasPrefix(json.NavigationItems[index].Url, json.Url) {
				json.NavigationItems[index].Url = strings.Replace(json.NavigationItems[index].Url, json.Url, "", 1)
				// If we removed the blog url, there should be a / in front of the url
				if !strings.HasPrefix(json.NavigationItems[index].Url, "/") {
					json.NavigationItems[index].Url = "/" + json.NavigationItems[index].Url
				}
			}
		}
		// Retrieve old blog settings for comparison
		blog, err := database.RetrieveBlog()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tempBlog := structure.Blog{Url: []byte(configuration.Config.Url), Title: []byte(json.Title), Description: []byte(json.Description), Logo: []byte(json.Logo), Cover: []byte(json.Cover), AssetPath: []byte("/assets/"), PostCount: blog.PostCount, PostsPerPage: json.PostsPerPage, ActiveTheme: json.ActiveTheme, NavigationItems: json.NavigationItems}
		err = methods.UpdateBlog(&tempBlog, userId)
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
	http.Error(w, "Not logged in!", http.StatusInternalServerError)
}

// API function to get user settings
func getApiUserHandler(w http.ResponseWriter, r *http.Request, params map[string]string) {
	userName := authentication.GetUserName(r)
	if userName != "" {
		userId, err := getUserId(userName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		id := params["id"]
		userIdToGet, err := strconv.ParseInt(id, 10, 64)
		if err != nil || userIdToGet < 1 {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else if userIdToGet != userId { // Make sure the authenticated user is only accessing his/her own data. TODO: Make sure the user is admin when multiple users have been introduced
			http.Error(w, "You don't have permission to access this data.", http.StatusForbidden)
			return
		}
		user, err := database.RetrieveUser(userIdToGet)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		userJson := userToJson(user)
		json, err := json.Marshal(userJson)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
		return
	}
	http.Error(w, "Not logged in!", http.StatusInternalServerError)
}

// API function to patch user settings
func patchApiUserHandler(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	userName := authentication.GetUserName(r)
	if userName != "" {
		userId, err := getUserId(userName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		decoder := json.NewDecoder(r.Body)
		var json JsonUser
		err = decoder.Decode(&json)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Make sure user id is over 0
		if json.Id < 1 {
			http.Error(w, "Wrong user id.", http.StatusInternalServerError)
			return
		} else if userId != json.Id { // Make sure the authenticated user is only changing his/her own data. TODO: Make sure the user is admin when multiple users have been introduced
			http.Error(w, "You don't have permission to change this data.", http.StatusInternalServerError)
			return
		}
		// Get old user data to compare
		tempUser, err := database.RetrieveUser(json.Id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Make sure user email is provided
		if json.Email == "" {
			json.Email = string(tempUser.Email)
		}
		// Make sure user name is provided
		if json.Name == "" {
			json.Name = string(tempUser.Name)
		}
		// Make sure user slug is provided
		if json.Slug == "" {
			json.Slug = tempUser.Slug
		}
		// Check if new name is already taken
		if json.Name != string(tempUser.Name) {
			_, err = database.RetrieveUserByName([]byte(json.Name))
			if err == nil {
				// The new user name is already taken. Assign the old name.
				// TODO: Return error that will be displayed in the admin interface.
				json.Name = string(tempUser.Name)
			}
		}
		// Check if new slug is already taken
		if json.Slug != tempUser.Slug {
			_, err = database.RetrieveUserBySlug(json.Slug)
			if err == nil {
				// The new user slug is already taken. Assign the old slug.
				// TODO: Return error that will be displayed in the admin interface.
				json.Slug = tempUser.Slug
			}
		}
		user := structure.User{Id: json.Id, Name: []byte(json.Name), Slug: json.Slug, Email: []byte(json.Email), Image: []byte(json.Image), Cover: []byte(json.Cover), Bio: []byte(json.Bio), Website: []byte(json.Website), Location: []byte(json.Location)}
		err = methods.UpdateUser(&user, userId)
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
			err = database.UpdateUserPassword(user.Id, encryptedPassword, date.GetCurrentTime(), json.Id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		// Check if the user name was changed. If so, update the session cookie to the new user name.
		if json.Name != string(tempUser.Name) {
			logInUser(json.Name, w)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("User settings updated!"))
		return
	}
	http.Error(w, "Not logged in!", http.StatusInternalServerError)
}

// API function to get the id of the currently authenticated user
func getApiUserIdHandler(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	userName := authentication.GetUserName(r)
	if userName != "" {
		userId, err := getUserId(userName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonUserId := JsonUserId{Id: userId}
		json, err := json.Marshal(jsonUserId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
		return
	}
	http.Error(w, "Not logged in!", http.StatusInternalServerError)
}

func getUserId(userName string) (int64, error) {
	user, err := database.RetrieveUserByName([]byte(userName))
	if err != nil {
		return 0, err
	}
	return user.Id, nil
}

func logInUser(name string, w http.ResponseWriter) {
	authentication.SetSession(name, w)
	userId, err := getUserId(name)
	if err != nil {
		log.Println("Couldn't get id of logged in user:", err)
	}
	err = database.UpdateLastLogin(date.GetCurrentTime(), userId)
	if err != nil {
		log.Println("Couldn't update last login date of a user:", err)
	}
}

func postsToJson(posts []structure.Post) *[]JsonPost {
	jsonPosts := make([]JsonPost, len(posts))
	for index := range posts {
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
	jsonPost.MetaDescription = string(post.MetaDescription)
	jsonPost.Image = string(post.Image)
	jsonPost.Date = post.Date
	tags := make([]string, len(post.Tags))
	for index := range post.Tags {
		tags[index] = string(post.Tags[index].Name)
	}
	jsonPost.Tags = strings.Join(tags, ",")
	return &jsonPost
}

func blogToJson(blog *structure.Blog) *JsonBlog {
	var jsonBlog JsonBlog
	jsonBlog.Url = string(blog.Url)
	jsonBlog.Title = string(blog.Title)
	jsonBlog.Description = string(blog.Description)
	jsonBlog.Logo = string(blog.Logo)
	jsonBlog.Cover = string(blog.Cover)
	jsonBlog.PostsPerPage = blog.PostsPerPage
	jsonBlog.Themes = templates.GetAllThemes()
	jsonBlog.ActiveTheme = blog.ActiveTheme
	jsonBlog.NavigationItems = blog.NavigationItems
	return &jsonBlog
}

func userToJson(user *structure.User) *JsonUser {
	var jsonUser JsonUser
	jsonUser.Id = user.Id
	jsonUser.Name = string(user.Name)
	jsonUser.Slug = user.Slug
	jsonUser.Email = string(user.Email)
	jsonUser.Image = string(user.Image)
	jsonUser.Cover = string(user.Cover)
	jsonUser.Bio = string(user.Bio)
	jsonUser.Website = string(user.Website)
	jsonUser.Location = string(user.Location)
	return &jsonUser
}

func InitializeAdmin(router *httptreemux.TreeMux) {
	// For admin panel
	router.GET("/admin/", adminHandler)
	router.GET("/admin/login/", getLoginHandler)
	router.POST("/admin/login/", postLoginHandler)
	router.GET("/admin/register/", getRegistrationHandler)
	router.POST("/admin/register/", postRegistrationHandler)
	router.GET("/admin/logout/", logoutHandler)
	router.GET("/admin/*filepath", adminFileHandler)

	// For admin API (no trailing slash)
	// Posts
	router.GET("/admin/api/posts/:number", apiPostsHandler)
	// Post
	router.GET("/admin/api/post/:id", getApiPostHandler)
	router.POST("/admin/api/post", postApiPostHandler)
	router.PATCH("/admin/api/post", patchApiPostHandler)
	router.DELETE("/admin/api/post/:id", deleteApiPostHandler)
	// Upload
	router.POST("/admin/api/upload", apiUploadHandler)
	// Images
	router.GET("/admin/api/images/:number", apiImagesHandler)
	router.DELETE("/admin/api/image", deleteApiImageHandler)
	// Blog
	router.GET("/admin/api/blog", getApiBlogHandler)
	router.PATCH("/admin/api/blog", patchApiBlogHandler)
	// User
	router.GET("/admin/api/user/:id", getApiUserHandler)
	router.PATCH("/admin/api/user", patchApiUserHandler)
	// User id
	router.GET("/admin/api/userid", getApiUserIdHandler)
}
