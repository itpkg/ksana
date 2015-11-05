package settings_test

import (
	"crypto/aes"
	"testing"

	kl "github.com/itpkg/ksana/logging"
	ko "github.com/itpkg/ksana/orm"
	ks "github.com/itpkg/ksana/settings"
	ku "github.com/itpkg/ksana/utils"
	_ "github.com/lib/pq"
)

func TestDatabase(t *testing.T) {
	db, err := ko.Open("test")
	if err != nil {
		t.Errorf("error on open: %v", err)
	}
	db.Logger = kl.NewStdoutLogger(kl.DEBUG)

	key, _ := ku.RandomBytes(32)
	cip, _ := aes.NewCipher(key)
	aes := ku.Aes{Cip: cip}

	ds := ks.DatabaseStore{Db: db, Aes: &aes}

	if err = db.Migrate(); err != nil {
		t.Errorf("error on migrate: %v", err)
	}

	test_store(t, &ds)
}

func test_store(t *testing.T, s ks.Store) {
	key := "hello"
	val := "你好, KSANA!"

	if err := s.Set(key, val, false); err != nil {
		t.Errorf("set error: %v", err)
	}
	var msg1 string
	if err := s.Get(key, &msg1); err != nil {
		t.Errorf("get error: %v", err)
	} else if msg1 != val {
		t.Errorf("%s != %s", val, msg1)
	}

	if err := s.Set(key, val, true); err != nil {
		t.Errorf("set error: %v", err)
	}
	var msg2 string
	if err := s.Get(key, &msg2); err != nil {
		t.Errorf("get error: %v", err)
	} else if msg2 != val {
		t.Errorf("%s != %s", val, msg2)
	}

}
