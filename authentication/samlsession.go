package authentication

import (
	"net/http"

	"github.com/crewjam/saml/samlsp"
)

// SAMLSession is an implemention of SessionHandler for SAML sessions
type SAMLSession struct {
	*samlsp.Middleware
}

var _ SessionHandler = (*SAMLSession)(nil)

// SetSession doesn't do anything for a SAML session
func (*SAMLSession) SetSession(userName string, response http.ResponseWriter) {
}

// GetUserName gets the email address from the session context
func (*SAMLSession) GetUserName(request *http.Request) string {
	session := samlsp.SessionFromContext(request.Context())
	if session == nil {
		return ""
	}
	return session.(samlsp.JWTSessionClaims).StandardClaims.Subject
}

// ClearSession removes the session cookie. Note that this probably
// doesn't really log you out of your IDP, but it does clear your
// cookie and won't get you a new one until you hit /admin again.
func (s *SAMLSession) ClearSession(response http.ResponseWriter, request *http.Request) {
	s.Session.DeleteSession(response, request)
	http.Redirect(response, request, "/", http.StatusFound)
}

// RequireSession is a function call wrapper to make sure you have
// a session. It will redirect to the auth flow if there is no session.
func (s *SAMLSession) RequireSession(callback func(http.ResponseWriter, *http.Request, map[string]string)) func(http.ResponseWriter, *http.Request, map[string]string) {
	return func(w http.ResponseWriter, r *http.Request, m map[string]string) {
		var err error
		if s.Session == nil {
			s.HandleStartAuthFlow(w, r)
			return
		}
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

		s.OnError(w, r, err)
		return
	}
}
