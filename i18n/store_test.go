package i18n_test

import (
	"testing"

	ki "github.com/itpkg/ksana/i18n"
	kl "github.com/itpkg/ksana/logging"
	ko "github.com/itpkg/ksana/orm"
	_ "github.com/lib/pq"
)

func TestDatabase(t *testing.T) {
	db, err := ko.Open("test")
	if err != nil {
		t.Errorf("error on open: %v", err)
	}

	db.Logger = kl.NewStdoutLogger(kl.DEBUG)

	ds := ki.DatabaseStore{Db: db}

	if err = db.Migrate(); err != nil {
		t.Errorf("error on migrate: %v", err)
	}

	test_store(t, &ds)
}

func test_store(t *testing.T, s ki.Store) {
	code := "hello"
	lang := "zh-CN"
	msg := "你好, KSANA!"

	if err := s.Set(lang, code, msg); err != nil {
		t.Errorf("set error: %v", err)
	}
	if msg1, err := s.Get(lang, code); err != nil {
		t.Errorf("get error: %v", err)
	} else if msg1 != msg {
		t.Errorf("%s != %s", msg, msg1)
	}

	if err := s.Loop(lang, func(code, msg string) error {
		t.Logf("%s: %s", code, msg)
		return nil
	}); err != nil {
		t.Errorf("loop error: %v", err)
	}

	/*
		c1:="aaa"
		m1:="AAA"
		c2:="bbb"
		m2:="BBB"
		if err:= s.SetM(lang, map[string]string{c1:m1,c2:m2}); err!=nil{
			log.Errorf("setM error: %v", err)
		}
		if err:= s.GetM(lang, map[string]string{c1:"", c2:""});err!=nil{
			log.Errorf("getM error: %v", err)
		}
	*/
}
