package job_test

import (
	"testing"

	kj "github.com/itpkg/ksana/job"
)

func TestSort(t *testing.T) {
	cfg := kj.Configuration{
		QueuesM: map[string]int{"aaa": 10, "bbb": 100, "ccc": 1, "ddd": 1000},
	}
	names := cfg.Queues()
	t.Logf("Before sort: %v", cfg)
	t.Logf("After sort: %v", names)
	if names[0] != "ddd" {
		t.Errorf("bad in sort")
	}
}
