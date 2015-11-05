package cmd_test

import (
	"testing"

	kc "github.com/itpkg/ksana/cmd"
)

func TestCmd(t *testing.T) {
	if err := kc.Run(); err != nil {
		t.Errorf("error on run: %v", err)
	}
}
