package lsp

import (
	"fmt"
	"sync"

	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/project/logging"
)

var _ logging.Logger = (*logger)(nil)

// logVerbosity mirrors the VS Code LogLevel enum values.
// Higher values mean less verbose output.
const (
	logVerbosityOff     = 0
	logVerbosityTrace   = 1
	logVerbosityDebug   = 2
	logVerbosityInfo    = 3
	logVerbosityWarning = 4
	logVerbosityError   = 5
)

type logger struct {
	server    *Server
	mu        sync.Mutex
	verbosity int
}

func newLogger(server *Server) *logger {
	return &logger{
		server:    server,
		verbosity: logVerbosityInfo,
	}
}

// maxVerbosityForMessageType returns the least-verbose log level at which
// messages of the given LSP MessageType should still be sent.
func maxVerbosityForMessageType(msgType lsproto.MessageType) int {
	switch msgType {
	case lsproto.MessageTypeError:
		return logVerbosityError
	case lsproto.MessageTypeWarning:
		return logVerbosityWarning
	case lsproto.MessageTypeInfo:
		return logVerbosityInfo
	default:
		return logVerbosityInfo
	}
}

func (l *logger) sendLogMessage(msgType lsproto.MessageType, message string) {
	if l == nil {
		return
	}

	if !l.server.initStarted.Load() {
		fmt.Fprintln(l.server.stderr, message)
		return
	}

	// Don't send messages that the client will filter out anyway.
	l.mu.Lock()
	verbosity := l.verbosity
	l.mu.Unlock()
	if verbosity > maxVerbosityForMessageType(msgType) {
		return
	}

	notification := lsproto.WindowLogMessageInfo.NewNotificationMessage(&lsproto.LogMessageParams{
		Type:    msgType,
		Message: message,
	})

	select {
	case l.server.outgoingQueue <- notification.Message():
		// sent
	case <-l.server.backgroundCtx.Done():
		fmt.Fprintln(l.server.stderr, message)
	}
}

func (l *logger) Log(msg ...any) {
	if l == nil {
		return
	}
	l.sendLogMessage(lsproto.MessageTypeInfo, fmt.Sprint(msg...))
}

func (l *logger) Logf(format string, args ...any) {
	if l == nil {
		return
	}
	l.sendLogMessage(lsproto.MessageTypeInfo, fmt.Sprintf(format, args...))
}

func (l *logger) Verbose() logging.Logger {
	if l == nil {
		return nil
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.verbosity > logVerbosityDebug {
		return nil
	}
	return l
}

func (l *logger) IsVerbose() bool {
	if l == nil {
		return false
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.verbosity <= logVerbosityDebug
}

func (l *logger) SetVerbose(verbose bool) {
	if l == nil {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if verbose {
		l.verbosity = logVerbosityDebug
	} else {
		l.verbosity = logVerbosityInfo
	}
}

func (l *logger) IsTracing() bool {
	if l == nil {
		return false
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.verbosity <= logVerbosityTrace
}

func (l *logger) SetVerbosity(verbosity int) {
	if l == nil {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.verbosity = verbosity
}

func (l *logger) Error(msg ...any) {
	if l == nil {
		return
	}
	l.sendLogMessage(lsproto.MessageTypeError, fmt.Sprint(msg...))
}

func (l *logger) Errorf(format string, args ...any) {
	if l == nil {
		return
	}
	l.sendLogMessage(lsproto.MessageTypeError, fmt.Sprintf(format, args...))
}

func (l *logger) Warn(msg ...any) {
	if l == nil {
		return
	}
	l.sendLogMessage(lsproto.MessageTypeWarning, fmt.Sprint(msg...))
}

func (l *logger) Warnf(format string, args ...any) {
	if l == nil {
		return
	}
	l.sendLogMessage(lsproto.MessageTypeWarning, fmt.Sprintf(format, args...))
}

func (l *logger) Info(msg ...any) {
	if l == nil {
		return
	}
	l.sendLogMessage(lsproto.MessageTypeInfo, fmt.Sprint(msg...))
}

func (l *logger) Infof(format string, args ...any) {
	if l == nil {
		return
	}
	l.sendLogMessage(lsproto.MessageTypeInfo, fmt.Sprintf(format, args...))
}
