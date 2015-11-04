package job

type Store interface {
	Push(queue string, msg *Message) error
	Pop(queues ...string) (string, *Message, error)
}
