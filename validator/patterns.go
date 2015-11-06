package validator

import (
	"regexp"
)

var patterns = make(map[string]*regexp.Regexp)

func Register(name string, pattern string) {
	patterns[name] = regexp.MustCompile(pattern)
}

//==============================================================================

func init() {
	Register("username", "")
	Register("email", "")
	Register("password", "")
	Register("not_empty", "")
}
