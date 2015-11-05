package logging

import (
	"fmt"
	"os"

	"github.com/itpkg/ksana/cmd"
)

func Open(env string) Logger {
	if cmd.IsProduction(env) {
		if log, err := NewSyslogLogger("ksana", INFO); err == nil {
			return log
		} else {
			fmt.Fprintf(os.Stderr, "error on create syslog logger: %v", err)
		}

		if log, err := NewFileLogger(fmt.Sprintf("%s.log", env), INFO); err == nil {
			return log
		} else {
			fmt.Fprintf(os.Stderr, "error on create file logger: %v", err)
		}
		return NewStdoutLogger(INFO)
	}

	return NewStdoutLogger(DEBUG)
}
