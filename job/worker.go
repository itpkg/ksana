package job

type Worker interface{
	Do(msg *Message) error
}

workers := map[string]Worker

//==============================================================================
func Register(queue string, worker Worker){
	workers[queue] :=  worker
}

func Run(file cfg) {
	fn := func() error{
		store.Get()
	}
	for i:=0; i<threads; i++{
		
	}
}

