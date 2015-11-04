package job

type Store interface {
	Push(queue string, args interface{}) error
	Do(queue string, args interface{}) error
}
