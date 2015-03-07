package templates

import (
	"github.com/kabukky/journey/structure"
)

// Helpers are created during parsing of the theme (template files). Helpers should never be altered during template execution (Helpers are shared across all requests).
type Helper struct {
	Name       string
	Arguments  []Helper
	Unescaped  bool
	Position   int
	Block      []byte
	Children   []Helper
	Function   func(*Helper, *structure.RequestData) []byte
	BodyHelper *Helper
}
