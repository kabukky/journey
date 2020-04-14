package authentication

import "net/http"

type SessionHandler interface {
	SetSession(userName string, response http.ResponseWriter)
	GetUserName(request *http.Request) string
	ClearSession(response http.ResponseWriter, request *http.Request)
	RequireSession(func(http.ResponseWriter, *http.Request, map[string]string)) func(http.ResponseWriter, *http.Request, map[string]string)
}
