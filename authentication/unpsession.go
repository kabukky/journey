package authentication

import (
	"github.com/gorilla/securecookie"
	"net/http"
)

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

type UsernamePasswordSession struct {
}

func (s *UsernamePasswordSession) SetSession(userName string, response http.ResponseWriter) {
	value := map[string]string{
		"name": userName,
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/admin/",
		}
		http.SetCookie(response, cookie)
	}
}

func (s *UsernamePasswordSession) GetUserName(request *http.Request) (userName string) {
	if cookie, err := request.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			userName = cookieValue["name"]
		}
	}
	return userName
}

func (s *UsernamePasswordSession) ClearSession(response http.ResponseWriter, _ *http.Request) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/admin/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
}

func (s *UsernamePasswordSession) RequireSession(callback func(http.ResponseWriter, *http.Request, map[string]string)) func(http.ResponseWriter, *http.Request, map[string]string) {
	return func(w http.ResponseWriter, r *http.Request, m map[string]string) {
		if s.GetUserName(r) != "" {
			callback(w, r, m)
		} else {
			if r.URL.Path == "/admin" || r.URL.Path == "/admin/" {
				http.Redirect(w, r, "/admin/login/", 302)
			} else {
				http.Error(w, "Not logged in", http.StatusForbidden)
			}
		}
	}
}

var _ SessionHandler = (*UsernamePasswordSession)(nil)
