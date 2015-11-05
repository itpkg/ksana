package logging

import (
	"fmt"
	"os"

	"github.com/itpkg/ksana/cmd"
)

func Open(env string) Logger {
	log_d := "log"

	if cmd.IsProduction(env) {
		if log, err := NewSyslogLogger("ksana", INFO); err == nil {
			return log
		} else {
			fmt.Fprintf(os.Stderr, "error on create syslog logger: %v", err)
		}

		if _, err := os.Stat(log_d); os.IsNotExist(err) {
			os.Mkdir(log_d, 0700)
		}

		if log, err := NewFileLogger(fmt.Sprintf("log/%s.log", env), INFO); err == nil {
			return log
		} else {
			fmt.Fprintf(os.Stderr, "error on create file logger: %v", err)
		}
		return NewStdoutLogger(INFO)
	}

	return NewStdoutLogger(DEBUG)
}
