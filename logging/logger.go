package logging

import (
	"fmt"
	"io"
	"time"
)

const (
	_ = iota
	DEBUG
	WARN
	INFO
	ERROR
)

func write(writer io.Writer, level int, format string, args ...interface{}) {
	var lvl string
	switch level {
	case DEBUG:
		lvl = "debug"
	case INFO:
		lvl = "info"
	case WARN:
		lvl = "warn"
	case ERROR:
		lvl = "error"
	default:
		lvl = "unknown"
	}
	fmt.Fprintf(writer, "%v[%s]: %s\n",
		time.Now(),
		lvl,
		fmt.Sprintf(format, args...),
	)
}

type Logger interface {
	Debug(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Info(format string, args ...interface{})
	Error(format string, args ...interface{})
	Close()
}
