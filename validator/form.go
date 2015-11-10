package validator

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

func Parse(req *http.Request, fm interface{}) error {
	req.ParseForm()
	return ParseValues(req.Form, fm)
}

func hasTag(t string) bool {
	return t != "" && t != "-"
}

func ParseValues(val url.Values, fm interface{}) error {
	mt := reflect.TypeOf(fm)
	mv := reflect.ValueOf(fm)
	if mt.Kind() != reflect.Ptr || mt.Elem().Kind() != reflect.Struct {
		return errors.New("bad type")
	}
	for i := 0; i < mv.Elem().NumField(); i++ {
		ft := mt.Elem().Field(i)
		fv := mv.Elem().Field(i)
		tag := ft.Tag.Get("valid")
		if !fv.CanSet() {
			return errors.New(fmt.Sprintf("%s@%s can not be set", ft.Name, mt))
		}

		pk := strings.ToLower(ft.Name)

		fk := ft.Type.Kind()
		switch fk {
		case reflect.String:
			pv := val.Get(pk)
			if hasTag(tag) {
				ph := handlers[tag]
				if ph == nil {
					return errors.New(fmt.Sprintf("can not find valiator for name %s", tag))
				}
				if err := ph(pv); err != nil {
					return err
				}
			}
			fv.Set(reflect.ValueOf(pv))
		default:
			return errors.New(fmt.Sprintf("unsupport kind of %s", fk))
		}

	}

	return nil
}
