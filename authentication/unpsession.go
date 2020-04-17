package authentication

import (
	"net/http"

	"github.com/gorilla/securecookie"
)

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

// UsernamePasswordSession is an implementation of sessions
// using cookies and database users
type UsernamePasswordSession struct {
}

// SetSession sets the session information.
func (s *UsernamePasswordSession) SetSession(userName string, response http.ResponseWriter) {
	value := map[string]string{
		"name": userName,
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/admin",
		}
		http.SetCookie(response, cookie)
	}
}

// GetUserName gets the name of the current user from the session,
// or returns "" if there is no session
func (s *UsernamePasswordSession) GetUserName(request *http.Request) (userName string) {
	if cookie, err := request.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			userName = cookieValue["name"]
		}
	}
	return userName
}

// ClearSession wipes the session cookie
func (s *UsernamePasswordSession) ClearSession(response http.ResponseWriter, request *http.Request) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/admin",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
	http.Redirect(response, request, "/admin/login", http.StatusFound)
}

// RequireSession is a function wrapper that requires a session. It will
// redirect to the login page if there isn't one.
func (s *UsernamePasswordSession) RequireSession(callback func(http.ResponseWriter, *http.Request, map[string]string)) func(http.ResponseWriter, *http.Request, map[string]string) {
	return func(w http.ResponseWriter, r *http.Request, m map[string]string) {
		if s.GetUserName(r) != "" {
			callback(w, r, m)
		} else {
			if r.URL.Path == "/admin" || r.URL.Path == "/admin/" {
				http.Redirect(w, r, "/admin/login", http.StatusFound)
			} else {
				http.Error(w, "Not logged in", http.StatusForbidden)
			}
		}
	}
}

var _ SessionHandler = (*UsernamePasswordSession)(nil)
