package job_test

import (
	"testing"
	"time"

	"github.com/garyburd/redigo/redis"
	kj "github.com/itpkg/ksana/job"
)

func TestStore(t *testing.T) {
	test_store(t, kj.NewRedisStore(&pool))
}

var pool = redis.Pool{
	MaxIdle:     5,
	IdleTimeout: 240 * time.Second,
	Dial: func() (redis.Conn, error) {
		return redis.Dial("tcp", "localhost:6379")

	},
	TestOnBorrow: func(c redis.Conn, t time.Time) error {
		_, err := c.Do("PING")
		return err
	},
}

func test_store(t *testing.T, s kj.Store) {
	queue := "test"
	args := []interface{}{"aaa", 111, time.Now()}

	if err := s.Push(queue, args...); err != nil {
		t.Errorf("error on push: %v", err)
	}

	var hello string
	var version int
	var now time.Time

	if err := s.Pop(queue, &hello, &version, &now); err == nil {
		args1 := []interface{}{hello, version, now}
		t.Logf("Get %v", args1)
		if args1[1] != args[1] || args1[0] != args[0] {
			t.Errorf("Buf want %v", args)
		}
	} else {
		t.Errorf("error on pop: %v", err)
	}

}
