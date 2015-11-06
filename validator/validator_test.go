package validator_test

import (
	"testing"

	kv "github.com/itpkg/ksana/validator"
)

type S struct {
	Username string `valid:"username"`
	Email    string `valid:"email"`
	Password string `valid:"password"`
	Ingnore  string `valid:"-"`
}

func TestValidator(t *testing.T) {
	s := S{}

	if e := kv.To(map[string]string{
		"username": "sdewrewr",
		"password": "234wfsewr",
	}, &s); e != nil {
		t.Errorf("bad in to: %v", e)
	}

	if e := kv.To(map[string]string{
		"username": "wer:?dfr",
		"password": " 2fsewre",
	}, &s); e == nil {
		t.Errorf("bad in to")
	}
}
