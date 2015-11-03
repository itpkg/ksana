package logging_test

import (
	"testing"

	kl "github.com/itpkg/ksana/logging"
)

func TestLogger(t *testing.T) {
	if log, err := kl.NewFileLogger("test.log", kl.DEBUG); err == nil {
		test_logger(t, log)
	} else {
		t.Errorf("error on open file logger: %v", err)
	}

	if log, err := kl.NewSyslogLogger("test", kl.DEBUG); err == nil {
		test_logger(t, log)
	} else {
		t.Errorf("error on open syslog logger: %v", err)
	}

	log := kl.NewStdoutLogger(kl.DEBUG)
	test_logger(t, log)
}

func test_logger(t *testing.T, log kl.Logger) {
	log.Debug("aaa")
	log.Info("bbb")
}
