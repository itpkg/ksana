package ioc

import (
	"errors"
	"fmt"
	"reflect"
)

var beans = make(map[string]interface{})

func Use(k string, v interface{}) {
	beans[k] = v
}

func Loop(fn func(string, interface{}) error) error {
	for k, v := range beans {
		if e := fn(k, v); e != nil {
			return e
		}
	}
	return nil
}

func Ping() error {
	for _, v := range beans {
		rt := reflect.TypeOf(v)
		rv := reflect.ValueOf(v)
		if rt.Kind() == reflect.Ptr && rt.Elem().Kind() == reflect.Struct {
			for i := 0; i < rv.Elem().NumField(); i++ {
				ft := rt.Elem().Field(i)
				fv := rv.Elem().Field(i)

				tag := ft.Tag.Get("inject")
				switch {
				case tag == "" || tag == "-":
					break
				default:
					if !fv.CanSet() {
						return errors.New(fmt.Sprintf("can not set field %s:%s@%v", ft.Name, tag, rt))
					}
					o := beans[tag]
					if o == nil {
						return errors.New(fmt.Sprintf("can not find bean name by %s", tag))
					}
					fv.Set(reflect.ValueOf(o))
				}
			}
		}
	}

	return nil
}
