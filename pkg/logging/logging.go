package logging

import (
	"fmt"
	"io"
	"log"
)

type LogLevel int8

const (
	Debug LogLevel = iota
	Info
	Warning
	Error
	Fatal
)

func (l *LogLevel) String() string {
	switch *l {
	case Debug:
		return "DEBUG"
	case Info:
		return "INFO"
	case Warning:
		return "WARNING"
	case Error:
		return "ERROR"
	case Fatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

type ILogger interface {
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Warningf(format string, args ...any)
	Errorf(format string, args ...any)
	Fatalf(format string, args ...any)
	Println(v ...any)
	Errorln(v ...any)
}

type Logger struct {
	ILogger
	level  LogLevel
	prefix string
	logger *log.Logger
	depth  int
}

var defaultInstance = NewLogger()

func NewLogger() *Logger {
	return NewLoggerWithLevel(Info)
}

func NewLoggerWithLevel(level LogLevel) *Logger {
	return NewLoggerWithPrefixAndLevel("", level)
}

func NewLoggerWithPrefix(prefix string) *Logger {
	return NewLoggerWithPrefixAndLevel(prefix, Info)
}

func NewLoggerWithPrefixAndLevel(prefix string, level LogLevel) *Logger {
	if level < Debug || level > Fatal {
		panic("Invalid log level")
	}
	log.Println()
	return &Logger{
		level:  level,
		prefix: prefix,
		logger: log.New(log.Writer(), "", log.Ldate|log.Ltime|log.Llongfile|log.Lmsgprefix),
		depth:  2,
	}

}

func (l *Logger) Debugf(format string, args ...any) {
	if l.level <= Debug {
		l.logger.Output(l.depth, l.fmtMessage(format, args...))
	}
}

func (l *Logger) Infof(format string, args ...any) {
	if l.level <= Info {
		l.logger.Output(l.depth, l.fmtMessage(format, args...))
	}
}

func (l *Logger) Warningf(format string, args ...any) {
	if l.level <= Warning {
		l.logger.Output(l.depth, l.fmtMessage(format, args...))
	}
}

func (l *Logger) Errorf(format string, args ...any) {
	if l.level <= Error {
		l.logger.Output(l.depth, l.fmtMessage(format, args...))
	}
}

func (l *Logger) Fatalf(format string, args ...any) {
	if l.level <= Fatal {
		l.logger.Output(l.depth, l.fmtMessage(format, args...))
	}
}

func (l *Logger) fmtMessage(format string, args ...any) string {
	if len(args) > 0 {
		format = fmt.Sprintf(format, args...)
	}
	if len(l.prefix) > 0 {
		format = fmt.Sprintf("%s %s", l.prefix, format)
	}
	return fmt.Sprintf("[%s] %s", l.level.String(), format)
}

func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

func (l *Logger) SetPrefix(prefix string) {
	l.prefix = prefix
}

func (l *Logger) GetLevel() LogLevel {
	return l.level
}

func (l *Logger) GetPrefix() string {
	return l.prefix
}

func (l *Logger) SetOutput(w io.Writer) {
	l.logger.SetOutput(w)
}

func (l *Logger) Println(v ...any) {
	t := l.level
	l.level = Info
	l.logger.Output(l.depth, l.fmtMessage(fmt.Sprint(v...)))
	l.level = t
}

func (l *Logger) Printf(format string, args ...any) {
	t := l.level
	l.level = Info
	l.logger.Output(l.depth, l.fmtMessage(format, args...))
	l.level = t
}

func (l *Logger) Errorln(v ...any) {
	t := l.level
	l.level = Error
	l.logger.Output(l.depth, l.fmtMessage(fmt.Sprint(v...)))
	l.level = t
}

func (l *Logger) SetDepth(d int) {
	l.depth = d
}

func Println(v ...any) {
	defaultInstance.Println(v...)
}

func Printf(format string, args ...any) {
	defaultInstance.Printf(format, args...)
}

func Errorln(v ...any) {
	defaultInstance.Errorln(v...)
}

func Errorf(format string, args ...any) {
	defaultInstance.Errorf(format, args...)
}
