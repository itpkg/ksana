package ksana

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

type Jwt struct {
	Key []byte `inject:"jwt.key"`
}

func (p *Jwt) key() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		kid, err := FromHex(token.Header["kid"].(string))
		if err != nil {
			return nil, err
		}
		return append(p.Key, kid...), nil
	}
}

func (p *Jwt) Parse(obj interface{}) (map[string]interface{}, error) {
	switch ty := obj.(type) {
	case string:
		if tk, er := jwt.Parse(obj.(string), p.key()); er == nil {
			return tk.Claims, nil
		} else {
			return nil, er
		}
	case *http.Request:
		if tk, er := jwt.ParseFromRequest(obj.(*http.Request), p.key()); er == nil {
			return tk.Claims, nil
		} else {
			return nil, er
		}
	default:
		return nil, errors.New(fmt.Sprintf("unknown support of type: %v", ty))
	}
}

func (p *Jwt) Create(claims map[string]interface{}) (string, error) {
	token := jwt.New(jwt.SigningMethodHS512)
	token.Claims = claims

	kid, err := RandomBytes(6)
	if err != nil {
		return "", err
	}

	token.Header["kid"] = ToHex(kid)
	return token.SignedString(append(p.Key, kid...))
}
