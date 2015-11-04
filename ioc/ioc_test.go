package ioc_test

import (
	"testing"
	"time"

	ki "github.com/itpkg/ksana/ioc"
)

type S struct {
	Version int
	Now     time.Time
}

func TestIoc(t *testing.T) {
	version := 111
	now := time.Now()
	if err := ki.Use(&S{Version: version, Now: now}); err != nil {
		t.Errorf("error on use: %v", err)
	}
	if err := ki.Map(map[string]interface{}{"a.1": "aaa", "b.1": "bbb"}); err != nil {
		t.Errorf("error on map: %v", err)
	}
	if err := ki.Loop(func(o interface{}) error {
		t.Logf("GET %v", o)
		return nil
	}); err != nil {
		t.Errorf("error on loop: %v", err)
	}
}
