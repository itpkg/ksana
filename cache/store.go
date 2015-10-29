package cache

import (
	"time"
)

type Store interface {
	Get(key string, val interface{}) error
	Set(key string, val interface{}, exp time.Duration) error
	Delete(key string) error
	Flush() error
}
