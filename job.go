package ksana

import (
	"github.com/jrallison/go-workers"
)

func Back(queue, class string, args interface{}) (string, error) {
	return workers.Enqueue(queue, class, args)
}
