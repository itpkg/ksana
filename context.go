package ksana

import (
	"github.com/facebookgo/inject"
)

var beans inject.Graph

func LoopEngine(f EngineHandler) error {
	for _, obj := range beans.Objects() {
		switch obj.Value.(type) {
		case Engine:
			if err := f(obj.Value.(Engine)); err != nil {
				return err
			}
		default:
		}
	}
	return nil
}

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
