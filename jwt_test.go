package ksana_test

import (
	"testing"
	"time"

	ks "github.com/itpkg/ksana"
)

func TestToken(t *testing.T) {
	key, _ := ks.RandomBytes(32)
	jwt := ks.Jwt{Key: key}

	val := map[string]interface{}{"aaa": "hello", "bbb": time.Now(), "ccc": 123}
	if tkn, err := jwt.Create(val); err == nil {
		t.Logf("Token: %s", tkn)

		if val1, err := jwt.Parse(tkn); err == nil {
			t.Logf("%v VS %v", val, val1)

			if val["ccc"].(int) != val["ccc"].(int) {
				t.Errorf("bad value")
			}

		} else {
			t.Errorf("error on parse: %v", err)
		}

	} else {
		t.Errorf("error on create token: %v", err)
	}

}
