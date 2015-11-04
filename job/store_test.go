package job_test

import (
	"testing"

	kj "github.com/itpkg/ksana/job"
	ku "github.com/itpkg/ksana/utils"
)

var redis_p = ku.OpenRedisPool("localhost", 6379, 0)
var hello = "Hello, Ksana-Job!"

func TestStore(t *testing.T) {
	test_store(t, kj.NewRedisStore(redis_p))
}

func test_store(t *testing.T, s kj.Store) {
	msg, err := kj.NewMessage(hello)
	if err != nil {
		t.Errorf("bad in new message: %v", err)
	}

	if err := s.Push("test", msg); err != nil {
		t.Errorf("bad in push: %v", err)
	}
	if name, msg1, err := s.Pop("test", "fuck"); err == nil {
		t.Logf("get message from [%s]", name)
		var hello1 string
		if err = msg1.Parse(&hello1); err != nil {
			t.Errorf("bad in parse message: %v", err)
		}
		if hello1 != hello {
			t.Errorf("Want %s, But get %v", hello, hello1)
		}

	} else {
		t.Errorf("bad in pop: %v", err)
	}
}
