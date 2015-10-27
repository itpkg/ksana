package ksana_test

import (
	"strings"
	"testing"
	"time"

	ks "github.com/itpkg/ksana"
)

const hello = "Hello, KSANA."

var now = time.Now()
var obj = map[string]interface{}{"message": hello, "ok": true, "value": 1.1, "time": now}

func TestPkgRoot(t *testing.T) {
	pr := ks.PkgRoot(&ks.BaseEngine{})
	if strings.HasSuffix(pr, "github.com/itpkg/ksana") {
		t.Logf("pkg root: %s", pr)
	} else {
		t.Errorf("bad pkg root")
	}
}

func TestJson(t *testing.T) {
	if jsn, err := ks.ToJson(obj); err == nil {
		t.Logf("TO JSON: %v", jsn)
		var obj1 = make(map[string]interface{}, 0)
		if err := ks.FromJson(jsn, &obj1); err == nil && obj["value"] == obj1["value"] {
			t.Logf("FROM JSON: %v", obj1)
		} else {
			t.Errorf("From json error! %v VS %v", obj, obj1)
		}
	} else {
		t.Errorf("To json error! %v", err)
	}
}

func TestBits(t *testing.T) {
	if buf, err := ks.ToBits(obj); err == nil {
		t.Logf("TO BITES: %v", buf)
		var obj1 = make(map[string]interface{}, 0)
		if err := ks.FromBits(buf, &obj1); err == nil && obj["value"] == obj1["value"] {
			t.Logf("FROM BITS: %v", obj1)
		} else {
			t.Errorf("From bites error! %v VS %v", obj, obj1)
		}
	} else {
		t.Errorf("To bites error! %v", err)
	}
}

func TestOthers(t *testing.T) {
	t.Logf("UUID: %s", ks.Uuid())

	if buf, err := ks.RandomBytes(16); err == nil {
		t.Logf("Random bytes %v", buf)
		bs := ks.ToBase64(buf)
		t.Logf("Base encode %s", bs)
		if buf1, err := ks.FromBase64(bs); err == nil && ks.Equal(buf, buf1) {
			t.Logf("Base64 Decode: %s", buf1)
		} else {
			t.Errorf("Decode base64 error! %v", err)
		}
	} else {
		t.Errorf("Random bytes error! %v", err)
	}

	if err := ks.Shell("/usr/bin/uname", "-a"); err != nil {
		t.Errorf("Run script error! %v", err)
	}
}
