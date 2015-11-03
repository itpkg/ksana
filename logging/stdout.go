package logging

import (
	"os"
)

type StdoutLogger struct {
	level int
}

func (p *StdoutLogger) Debug(format string, args ...interface{}) {
	write(os.Stdout, DEBUG, format, args...)
}

func (p *StdoutLogger) Warn(format string, args ...interface{}) {
	write(os.Stdout, WARN, format, args...)
}

func (p *StdoutLogger) Info(format string, args ...interface{}) {
	write(os.Stdout, INFO, format, args...)
}

func (p *StdoutLogger) Error(format string, args ...interface{}) {
	write(os.Stdout, ERROR, format, args...)
}

func (p *StdoutLogger) Close() {}

//==============================================================================

func NewStdoutLogger(level int) Logger {
	return &StdoutLogger{level: level}
}
