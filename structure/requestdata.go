// +build !noplugins

package structure

import (
	"github.com/yuin/gopher-lua"
)

// RequestData used for template/helper execution. Contains data specific to the incoming request.
type RequestData struct {
	PluginVMs              map[string]*lua.LState
	Posts                  []Post
	Blog                   *Blog
	CurrentTag             *Tag
	CurrentIndexPage       int
	CurrentPostIndex       int
	CurrentTagIndex        int
	CurrentNavigationIndex int
	CurrentHelperContext   int      // 0 = index, 1 = post, 2 = tag, 3 = author, 4 = navigation - used by block helpers
	CurrentTemplate        int      // 0 = index, 1 = post, 2 = tag, 3 = author - never changes during execution. Used by funcs like body_classFunc etc to output the correct class
	ContentForHelpers      []Helper // contentFor helpers that are attached to the currently rendering helper
	CurrentPath            string   // path of the the url of this request
}
