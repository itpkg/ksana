package job

type Store interface {
	Push(queue string, args ...interface{}) error
	Pop(queue string, args ...interface{}) error
}
