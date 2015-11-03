package logging

import (
	"os"
)

type FileLogger struct {
	out   *os.File
	level int
}

func (p *FileLogger) Debug(format string, args ...interface{}) {
	write(p.out, DEBUG, format, args...)
}

func (p *FileLogger) Warn(format string, args ...interface{}) {
	write(p.out, WARN, format, args...)
}

func (p *FileLogger) Info(format string, args ...interface{}) {
	write(p.out, INFO, format, args...)
}

func (p *FileLogger) Error(format string, args ...interface{}) {
	write(p.out, ERROR, format, args...)
}

func (p *FileLogger) Close() {
	p.out.Close()
}

//==============================================================================

func NewFileLogger(name string, level int) (Logger, error) {
	fi, err := os.OpenFile(name, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0600)
	if err != nil {
		return nil, err
	}
	return &FileLogger{out: fi, level: level}, nil
}
