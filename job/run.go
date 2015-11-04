package job

import (
	"errors"

	"github.com/itpkg/ksana/logging"
	"github.com/itpkg/ksana/utils"
)

func Run(file string, store Store, logger logging.Logger) error {
	cfg := Configuration{}
	if err := utils.FromToml(file, &cfg); err != nil {
		return err
	}
	timeout := cfg.Timeout()
	queues := cfg.Queues()

	fn := func() {
		que, msg, err := store.Pop(timeout, queues...)
		if err == nil {
			logger.Info("get job %s@%s", msg.Id, que)
			wrk := workers[que]
			if wrk == nil {
				err = errors.New("unknown worker")
			} else {
				err = wrk.Do(msg)
			}
			if err == nil {

				logger.Info("job done %s@%s", msg.Id, que)
			} else {
				logger.Error("error on job %s@%s", msg.Id, que)
			}

		}
		store.Done(que, msg, err)
	}

	for i := 0; i < cfg.Workers; i++ {
		go fn()
	}
	return nil
}
