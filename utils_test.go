package ksana_test

import (
	"testing"
	"time"

	"github.com/itpkg/ksana"
)

const hello = "Hello, IT-PACKAGE!!!"

var now = time.Now()
var obj = map[string]interface{}{"message": hello, "ok": true, "value": 1.1, "time": now}

func TestJson(t *testing.T) {
	if jsn, err := ksana.ToJson(obj); err == nil {
		t.Logf("TO JSON: %v", jsn)
		var obj1 = make(map[string]interface{}, 0)
		if err := ksana.FromJson(jsn, &obj1); err == nil && obj["value"] == obj1["value"] {
			t.Logf("FROM JSON: %v", obj1)
		} else {
			t.Errorf("From json error! %v VS %v", obj, obj1)
		}
	} else {
		t.Errorf("To json error! %v", err)
	}
}

func TestBits(t *testing.T) {
	if buf, err := ksana.ToBits(obj); err == nil {
		t.Logf("TO BITES: %v", buf)
		var obj1 = make(map[string]interface{}, 0)
		if err := ksana.FromBits(buf, &obj1); err == nil && obj["value"] == obj1["value"] {
			t.Logf("FROM BITS: %v", obj1)
		} else {
			t.Errorf("From bites error! %v VS %v", obj, obj1)
		}
	} else {
		t.Errorf("To bites error! %v", err)
	}
}

func TestOthers(t *testing.T) {
	t.Logf("UUID: %s", ksana.Uuid())

	if buf, err := ksana.RandomBytes(16); err == nil {
		t.Logf("Random bytes %v", buf)
		bs := ksana.ToBase64(buf)
		t.Logf("Base encode %s", bs)
		if buf1, err := ksana.FromBase64(bs); err == nil && ksana.Equal(buf, buf1) {
			t.Logf("Base64 Decode: %s", buf1)
		} else {
			t.Errorf("Decode base64 error! %v", err)
		}
	} else {
		t.Errorf("Random bytes error! %v", err)
	}

	if err := ksana.Shell("/usr/bin/uname", "-a"); err != nil {
		t.Errorf("Run script error! %v", err)
	}
}
