package project

import (
	"fmt"
	"io"
	"strings"
	"time"
)

type LogLevel int

const (
	LogLevelTerse LogLevel = iota
	LogLevelNormal
	LogLevelRequestTime
	LogLevelVerbose
)

type Logger struct {
	outputs []io.Writer
	level   LogLevel
	inGroup bool
	seq     int
}

func NewLogger(outputs []io.Writer, level LogLevel) *Logger {
	return &Logger{outputs: outputs, level: level}
}

func (l *Logger) PerfTrace(s string) {
	l.msg(s, "Perf")
}

func (l *Logger) Info(s string) {
	l.msg(s, "Info")
}

func (l *Logger) Error(s string) {
	l.msg(s, "Err")
}

func (l *Logger) StartGroup() {
	l.inGroup = true
}

func (l *Logger) EndGroup() {
	l.inGroup = false
}

func (l *Logger) LoggingEnabled() bool {
	return len(l.outputs) > 0
}

func (l *Logger) HasLevel(level LogLevel) bool {
	return l.LoggingEnabled() && l.level >= level
}

func (l *Logger) msg(s string, messageType string) {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("%s %d", messageType, l.seq))
	builder.WriteString(strings.Repeat(" ", max(0, 10-builder.Len())))
	builder.WriteRune('[')
	builder.WriteString(time.Now().Format("15:04:05.000"))
	builder.WriteString("] ")
	builder.WriteString(s)
	builder.WriteRune('\n')
	if !l.inGroup {
		l.seq++
	}
	l.write([]byte(builder.String()))
}

func (l *Logger) write(s []byte) {
	for _, output := range l.outputs {
		output.Write(s) //nolint: errcheck
	}
}
