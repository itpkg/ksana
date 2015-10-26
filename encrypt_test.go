package ksana_test

import (
	"crypto/aes"
	"crypto/sha512"
	"testing"

	ks "github.com/itpkg/ksana"
)

var hello = []byte("Hello, KSANA!")

const salt_len = 32

func TestAes(t *testing.T) {
	key, _ := ks.RandomBytes(32)
	c, _ := aes.NewCipher(key)
	a := ks.Aes{Cip: c}

	dest1, iv1, _ := a.Encrypt(hello)
	dest2, iv2, _ := a.Encrypt(hello)
	t.Logf("AES1(%d, iv=%x): %s => %x", len(dest1), iv1, hello, dest1)
	t.Logf("AES2(%d, iv=%x): %s => %x", len(dest2), iv2, hello, dest2)

	src := a.Decrypt(dest1, iv1)
	if string(src) != string(hello) {
		t.Errorf("val == %x, want %x", src, hello)
	}

}

func TestHmac(t *testing.T) {
	key, _ := ks.RandomBytes(32)
	hm := ks.Hmac{
		Fn:  sha512.New,
		Key: key,
	}

	dest1 := hm.Sum([]byte(hello))
	dest2 := hm.Sum([]byte(hello))

	t.Logf("HMAC1(%d): %x", len(dest1), dest1)
	t.Logf("HMAC2(%d): %x", len(dest2), dest2)
	if !hm.Equal(dest1, dest2) {
		t.Errorf("HMAC FAILED!")
	}

}

func TestMd5(t *testing.T) {
	t.Logf("md5: %s", ks.Md5(hello))
}
func TestSha(t *testing.T) {

	t.Logf("sha512: %s", ks.Sha512(hello))

	if s, e := ks.Ssha512(hello, salt_len); e == nil {
		t.Logf("ssha512: %s", s)
		if b, e := ks.Csha512(s, hello); e == nil {
			if !b {
				t.Errorf("csha512 error")
			}
		} else {
			t.Errorf("csha512 error: %v", e)
		}
	} else {
		t.Errorf("ssha512 error: %v", e)
	}
}
