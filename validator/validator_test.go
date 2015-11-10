package validator_test

import (
	"net/url"
	"reflect"
	"testing"

	kv "github.com/itpkg/ksana/validator"
)

type S struct {
	Username string
	Email    string `valid:"email"`
	Password string `valid:"password"`
	Ingnore1 string `valid:"-"`
	Ingnore2 string
}

func values() url.Values {

	v := url.Values{}
	v.Set("username", "Ava")
	v.Add("password", "123456")
	v.Add("email", "test@test.com")
	v.Add("emails", "aaa@aaa.com")
	v.Add("emails", "bbb@bbb.com")
	return v
}

func TestValidator(t *testing.T) {
	s := S{}
	v := values()
	if e := kv.ParseValues(v, &s); e == nil {
		t.Logf("S = %v", s)
	} else {
		t.Errorf("bad in to: %v", e)
	}

}

func TestStruct(t *testing.T) {
	s := S{}

	st1 := reflect.ValueOf(&s)
	t.Logf("Kind: %v", st1.Kind())
	t2 := st1.Elem().Interface()

	st := reflect.TypeOf(t2)
	t.Logf("Kind: %v", st.Kind())

	for i := 0; i < st.NumField(); i++ {
		sf := st.Field(i)
		t.Logf(
			"Field.%d: name=%s, type=%v, tag=%v",
			i,
			sf.Name,
			sf.Type.Kind(),
			sf.Tag.Get("valid"),
		)
	}
}
