package base_test

import (
	"crypto/aes"
	"crypto/sha512"
	"testing"

	ks "github.com/itpkg/ksana"
	kb "github.com/itpkg/ksana/base"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

func TestMigrate(t *testing.T) {
	en := kb.BaseEngine{Db: get_db(t)}

	if err := en.Migrate(); err != nil {
		t.Errorf("error on migrate: %v", err)
	}
}

func TestDao(t *testing.T) {
	db := get_db(t)

	key, _ := ks.RandomBytes(32)
	cip, _ := aes.NewCipher(key)

	dao := kb.Dao{
		Aes: &ks.Aes{Cip: cip},
		Hmac: &ks.Hmac{
			Fn:  sha512.New,
			Key: key,
		},
	}

	k1 := "aaa"
	k2 := "bbb"
	v11 := 1111
	v21 := 2222
	if err := dao.Set(db, k1, v11, true); err != nil {
		t.Errorf("error on set1: %v", err)
	}

	if err := dao.Set(db, k2, v21, false); err != nil {
		t.Errorf("error on set1: %v", err)
	}
	var v12 int
	var v22 int
	if err := dao.Get(db, k1, &v12); err != nil || v11 != v12 {
		t.Errorf("error on get1: %v", err)
	}
	if err := dao.Get(db, k2, &v22); err != nil || v21 != v22 {
		t.Errorf("error on get2: %v", err)
	}

}

func get_db(t *testing.T) *gorm.DB {
	if db, err := gorm.Open("sqlite3", "test.db"); err == nil {
		db.LogMode(true)
		db.DB().Ping()
		return &db
	} else {
		t.Errorf("error in open database: %v", err)
	}
	return nil
}
