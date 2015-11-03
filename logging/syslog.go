package logging

import (
	"fmt"
	"log/syslog"
)

type SyslogLogger struct {
	out *syslog.Writer
}

func (p *SyslogLogger) Debug(format string, args ...interface{}) {
	p.out.Debug(fmt.Sprintf(format, args...))
}

func (p *SyslogLogger) Warn(format string, args ...interface{}) {
	p.out.Warning(fmt.Sprintf(format, args...))
}

func (p *SyslogLogger) Info(format string, args ...interface{}) {
	p.out.Info(fmt.Sprintf(format, args...))
}

func (p *SyslogLogger) Error(format string, args ...interface{}) {
	p.out.Err(fmt.Sprintf(format, args...))
}

func (p *SyslogLogger) Close() {
	p.out.Close()
}

//==============================================================================

func NewSyslogLogger(tag string, level int) (Logger, error) {
	var lvl syslog.Priority
	switch level {
	case DEBUG:
		lvl = syslog.LOG_DEBUG
	case INFO:
		lvl = syslog.LOG_INFO
	case WARN:
		lvl = syslog.LOG_WARNING
	case ERROR:
		lvl = syslog.LOG_ERR
	default:
		lvl = syslog.LOG_DEBUG
	}

	wrt, err := syslog.New(lvl, tag)
	if err != nil {
		return nil, err
	}
	return &SyslogLogger{out: wrt}, nil
}
