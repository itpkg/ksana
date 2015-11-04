package job

type Worker interface {
	Do(msg *Message) error
}

var workers = make(map[string]Worker, 0)

//==============================================================================
func Register(queue string, worker Worker) {
	workers[queue] = worker
}

func Run(file cfg) {
}
