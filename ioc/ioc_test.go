package ioc_test

import (
	"testing"
	"time"

	ki "github.com/itpkg/ksana/ioc"
)

type S1 struct {
	Version int        `inject:"version"`
	Now     *time.Time `inject:"now"`
	S2      *S2        `inject:"s2"`
}

type S2 struct {
	Version int
	Now     *time.Time `inject:"now"`
}

func TestIoc(t *testing.T) {
	now := time.Now()
	ki.Use("now", &now)
	ki.Use("version", 20141110)
	s1 := S1{}
	s2 := S2{}
	ki.Use("s2", &s2)
	ki.Use("s1", &s1)

	if err := ki.Ping(); err == nil {
		ki.Loop(func(k string, v interface{}) error {
			t.Logf("%s = %v", k, v)
			return nil
		})
	} else {
		t.Errorf("bad in fill %v", err)
	}
}
