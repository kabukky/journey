package authentication

import (
	"github.com/crewjam/saml/samlsp"
	"net/http"
)

type SAMLSession struct {
	*samlsp.Middleware
}

var _ SessionHandler = (*SAMLSession)(nil)

func (_ *SAMLSession) SetSession(userName string, response http.ResponseWriter) {
}
func (_ *SAMLSession) GetUserName(request *http.Request) string {
	session := samlsp.SessionFromContext(request.Context())
	if session == nil {
		return ""
	}
	return session.(samlsp.JWTSessionClaims).StandardClaims.Subject
}
func (s *SAMLSession) ClearSession(response http.ResponseWriter, request *http.Request) {
	s.Session.DeleteSession(response, request)
	http.Redirect(response, request, "/admin", 302)
}
func (s *SAMLSession) RequireSession(callback func(http.ResponseWriter, *http.Request, map[string]string)) func(http.ResponseWriter, *http.Request, map[string]string) {
	return func(w http.ResponseWriter, r *http.Request, m map[string]string) {
		var err error
		if s.Session == nil {
			s.HandleStartAuthFlow(w, r)
			return
		} else {
			session, err := s.Session.GetSession(r)
			if session != nil {
				r = r.WithContext(samlsp.ContextWithSession(r.Context(), session))
				callback(w, r, m)
				return
			}
			if err == samlsp.ErrNoSession {
				s.HandleStartAuthFlow(w, r)
				return
			}
		}

		s.OnError(w, r, err)
		return
	}
}
