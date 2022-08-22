package authentication

import "net/http"

// SessionHandler is the interface to session functions. There are
// two, the Username/Password and SAML versions
type SessionHandler interface {
	SetSession(userName string, response http.ResponseWriter)
	GetUserName(request *http.Request) string
	ClearSession(response http.ResponseWriter, request *http.Request)
	RequireSession(func(http.ResponseWriter, *http.Request, map[string]string)) func(http.ResponseWriter, *http.Request, map[string]string)
}
