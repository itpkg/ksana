package validator

import (
	"errors"
	"fmt"
	"regexp"
)

type Handler func(string) error

var handlers = make(map[string]Handler)

func RegisterR(name string, pattern string) {
	pat := regexp.MustCompile(pattern)
	handlers[name] = func(v string) error {
		if pat.MatchString(v) {
			return nil
		} else {
			return errors.New(fmt.Sprintf("not match with %v", pat))
		}
	}

}

func Register(name string, h Handler) {
	handlers[name] = h
}

//==============================================================================

func init() {
	RegisterR("email", `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`)
	RegisterR("url", `[a-zA-z]+://[^\s]*`)
	RegisterR("ip", `^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$`)
	RegisterR("password", `^\w{6,20}$`)
}
