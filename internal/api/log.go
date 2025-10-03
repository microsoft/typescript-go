package api

import "github.com/microsoft/typescript-go/internal/project/logging"

type NoLogger struct{}

// SetVerbose implements logging.Logger.
func (n NoLogger) SetVerbose(verbose bool) {
	panic("unimplemented")
}

var _ logging.Logger = (*NoLogger)(nil)

func (n NoLogger) Log(msg ...any)                  {}
func (n NoLogger) Logf(format string, args ...any) {}
func (n NoLogger) Write(msg string)                {}
func (n NoLogger) Verbose() logging.Logger {
	return n
}

func (n NoLogger) IsVerbose() bool {
	return false
}
