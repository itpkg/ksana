package cache_test

import (
	"testing"
	"time"

	kc "github.com/itpkg/ksana/cache"
	ku "github.com/itpkg/ksana/utils"
)

func TestRedisStore(t *testing.T) {
	test_store(t, &kc.RedisStore{Pool: ku.OpenRedisPool("localhost", 6379, 0)})
}

func test_store(t *testing.T, s kc.Store) {
	key := "aaa"
	val := 1234

	if err := s.Set(key, val, 3*time.Hour); err != nil {
		t.Errorf("error on set: %v", err)
	}
	var val1 int
	if err := s.Get(key, &val1); err == nil {
		if val1 != val {
			t.Errorf("error on get: %s != %s", val, val1)
		}
	} else {
		t.Errorf("error on get: %v", err)
	}

	if err := s.Delete(key); err != nil {
		t.Errorf("error on delete: %v", err)
	}
	s.Set(key+"aaa", val+2, 3*time.Hour)
	if err := s.Flush(); err != nil {
		t.Errorf("error on flush: %v", err)
	}
}
