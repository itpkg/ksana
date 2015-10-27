package ksana_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/BurntSushi/toml"
	ks "github.com/itpkg/ksana"
)

func TestReflect(t *testing.T) {
	en := ks.BaseEngine{}
	et1 := reflect.TypeOf((*ks.BaseEngine)(nil)).Elem()
	et2 := reflect.TypeOf(&en).Elem()
	t.Logf("%s <==> %s", et1.PkgPath(), et2.PkgPath())

}

func TestLocale(t *testing.T) {
	lcs := make(map[string]string, 0)
	lcs["aaa"] = "AAA"
	lcs["bbb"] = "BBB"

	end := toml.NewEncoder(os.Stdout)
	if err := end.Encode(lcs); err != nil {
		t.Errorf("bad on toml encode: %v", err)
	}
}
