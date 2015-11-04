package job

import (
	"time"
)

type Store interface {
	Push(queue string, msg *Message) error
	Pop(timeout time.Duration, queues ...string) (string, *Message, error)
	Done(queue string, msg *Message, err error)
}
