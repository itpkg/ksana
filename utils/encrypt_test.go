package utils_test

import (
	"crypto/aes"
	"crypto/sha512"
	"testing"

	ks "github.com/itpkg/ksana/utils"
)

const salt_len = 32

func TestAes(t *testing.T) {
	key, _ := ks.RandomBytes(32)
	c, _ := aes.NewCipher(key)
	a := ks.Aes{Cip: c}

	dest1, _ := a.Encrypt([]byte(hello))
	dest2, _ := a.Encrypt([]byte(hello))
	t.Logf("AES1: %s", dest1)
	t.Logf("AES2: %s", dest2)

	src, _ := a.Decrypt(dest1)
	if string(src) != hello {
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
	t.Logf("md5: %s", ks.Md5([]byte(hello)))
}

func TestSha(t *testing.T) {

	t.Logf("sha512: %s", ks.Sha512([]byte(hello)))

	if s, e := ks.Ssha512([]byte(hello), salt_len); e == nil {
		t.Logf("ssha512: %s", s)
		if b, e := ks.Csha512(s, []byte(hello)); e == nil {
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

func TestDovecot(t *testing.T) {
	if s, e := ks.Ssha512([]byte(hello), 6); e == nil {
		t.Logf("Please run: doveadm pw -t {SSHA512}%s -p \"%s\"", s, hello)
	} else {
		t.Errorf("bad in test dovecot: %v", e)
	}
}
