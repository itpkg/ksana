package ioc

import (
	"github.com/facebookgo/inject"
)

var beans inject.Graph

func Map(objects map[string]interface{}) error {
	items := make([]*inject.Object, 0)
	for k, v := range objects {
		items = append(items, &inject.Object{Value: v, Name: k})
	}
	return beans.Provide(items...)
}

func Use(objects ...interface{}) error {
	items := make([]*inject.Object, 0)
	for _, v := range objects {
		items = append(items, &inject.Object{Value: v})
	}
	return beans.Provide(items...)
}

func Loop(fn func(_ interface{}) error) error {
	for _, obj := range beans.Objects() {
		if err := fn(obj); err != nil {
			return err
		}
	}
	return nil
}

func Build() error {
	return beans.Populate()
}
